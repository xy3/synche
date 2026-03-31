// Package client implements the concurrent chunk upload client.
package client

import (
	"bytes"
	"context"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"github.com/theo/synche2/internal/cache"
	"github.com/theo/synche2/internal/chunk"
	"github.com/theo/synche2/internal/protocol"
	"golang.org/x/net/http2"
)

const (
	// maxUploadRetries is how many times to retry a failed chunk upload.
	maxUploadRetries = 3
	// retryBaseDelay is the base delay between retries (doubled each attempt).
	retryBaseDelay = 500 * time.Millisecond
)

// Config holds client configuration.
type Config struct {
	ServerURL   string // e.g. "http://localhost:8420"
	DevicePath  string // e.g. "/dev/sda" or a file path
	Concurrency int    // number of parallel upload workers
	ProbeBatch  int    // how many hashes to send per probe request
	CacheDir    string // local manifest cache dir (empty = ~/.cache/synche)
	NoCache     bool   // disable local caching
	APIKey      string // API key for server authentication (empty = no auth)
}

// Stats tracks upload progress.
type Stats struct {
	TotalChunks   int64
	UploadedNew   int64
	SkippedServer int64 // server already had (discovered via probe)
	SkippedCache  int64 // unchanged since last run (no network needed)
	Failed        int64
	BytesSent     int64
	BytesRead     int64
	BytesTotal    int64
	StartTime     time.Time
}

// Uploader handles concurrent chunk uploading with dedup probing.
type Uploader struct {
	cfg   Config
	http  *http.Client
	stats Stats
	cache *cache.Cache
}

// NewUploader creates an uploader with the given configuration.
// Uses HTTP/2 cleartext (h2c) for plain HTTP and HTTP/2 over TLS for HTTPS,
// multiplexing all requests over a single TCP connection.
func NewUploader(cfg Config) *Uploader {
	if cfg.Concurrency < 1 {
		cfg.Concurrency = 8
	}
	if cfg.ProbeBatch < 1 {
		cfg.ProbeBatch = 256
	}

	var transport http.RoundTripper

	if strings.HasPrefix(cfg.ServerURL, "https://") {
		// HTTP/2 over TLS — Go's default transport handles this via ALPN.
		transport = &http2.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true, // allow self-signed certs
			},
		}
	} else {
		// HTTP/2 cleartext (h2c) — prior knowledge, single TCP connection.
		transport = &http2.Transport{
			AllowHTTP: true,
			DialTLSContext: func(ctx context.Context, network, addr string, _ *tls.Config) (net.Conn, error) {
				// h2c: dial a plain TCP connection (no TLS).
				var d net.Dialer
				return d.DialContext(ctx, network, addr)
			},
		}
	}

	return &Uploader{
		cfg: cfg,
		http: &http.Client{
			Timeout:   120 * time.Second,
			Transport: transport,
		},
		stats: Stats{StartTime: time.Now()},
	}
}

// setAuth adds the Authorization header to a request if an API key is configured.
func (u *Uploader) setAuth(req *http.Request) {
	if u.cfg.APIKey != "" {
		req.Header.Set("Authorization", "Bearer "+u.cfg.APIKey)
	}
}

type pendingChunk struct {
	meta protocol.ChunkMeta
	data []byte
}

// uploadedHash tracks a chunk that was confirmed present on the server
// (either uploaded successfully or server already had it). Only these
// should be written to the local cache.
type uploadedHash struct {
	index uint64
	hash  string
}

// Run executes the upload. If the source is a directory, it uploads each file
// as a separate manifest. Otherwise it uploads a single file/device.
func (u *Uploader) Run() error {
	fi, err := os.Stat(u.cfg.DevicePath)
	if err != nil {
		return fmt.Errorf("stat source: %w", err)
	}
	if fi.IsDir() {
		return u.runDirectory(u.cfg.DevicePath)
	}
	_, err = u.runFile(u.cfg.DevicePath)
	if err != nil {
		return err
	}
	return nil
}

// FileResult holds the result of uploading a single file.
type FileResult struct {
	Path       string
	ManifestID string
	Stats      Stats
}

// dirFileInfo holds all chunks for a single file during directory upload.
type dirFileInfo struct {
	relPath  string
	absPath  string
	manifest protocol.Manifest
	chunks   []pendingChunk // only populated for chunks that need uploading
}

// runDirectory walks a directory and uploads each regular file as its own
// manifest. Uses a multi-phase pipeline to minimize round-trips:
//   - Phase 1: Hash all files concurrently (local I/O, parallel)
//   - Phase 2: Batch probe ALL hashes across all files (few round-trips)
//   - Phase 3: Upload only needed chunks (concurrent)
//   - Phase 4: Upload all manifests concurrently
func (u *Uploader) runDirectory(dir string) error {
	var files []string
	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() && info.Size() > 0 {
			files = append(files, path)
		}
		return nil
	})
	if err != nil {
		return fmt.Errorf("walk directory: %w", err)
	}

	if len(files) == 0 {
		return fmt.Errorf("no files found in %s", dir)
	}

	log.Printf("found %d files in %s", len(files), dir)
	totalStart := time.Now()

	// ── Phase 1: Hash all files concurrently ──────────────────────────
	log.Printf("phase 1: hashing all files...")
	phaseStart := time.Now()

	fileInfos := make([]dirFileInfo, len(files))
	// allChunkData maps hash -> data for chunks we may need to upload.
	allChunkData := make(map[string][]byte)
	var dataMu sync.Mutex

	hashWorkers := u.cfg.Concurrency
	if hashWorkers > len(files) {
		hashWorkers = len(files)
	}

	fileCh := make(chan int, hashWorkers)
	var hashWg sync.WaitGroup
	var hashErr error
	var hashErrMu sync.Mutex

	for w := 0; w < hashWorkers; w++ {
		hashWg.Add(1)
		go func() {
			defer hashWg.Done()
			for idx := range fileCh {
				path := files[idx]
				relPath, _ := filepath.Rel(dir, path)

				pipeline, err := chunk.NewPipeline(path, 2) // light concurrency per file
				if err != nil {
					hashErrMu.Lock()
					if hashErr == nil {
						hashErr = fmt.Errorf("hash %s: %w", relPath, err)
					}
					hashErrMu.Unlock()
					continue
				}

				fi := dirFileInfo{
					relPath: relPath,
					absPath: path,
					manifest: protocol.Manifest{
						SourceDevice: path,
						TotalBytes:   pipeline.DeviceSize(),
						ChunkSize:    protocol.ChunkSize,
					},
				}

				results := pipeline.Run()
				var metas []protocol.ChunkMeta
				for r := range results {
					metas = append(metas, r.Meta)
					dataMu.Lock()
					if _, exists := allChunkData[r.Meta.Hash]; !exists {
						allChunkData[r.Meta.Hash] = r.Data
					}
					dataMu.Unlock()
				}

				sort.Slice(metas, func(i, j int) bool {
					return metas[i].Index < metas[j].Index
				})
				fi.manifest.Chunks = metas
				fileInfos[idx] = fi
			}
		}()
	}

	for i := range files {
		fileCh <- i
	}
	close(fileCh)
	hashWg.Wait()

	if hashErr != nil {
		return hashErr
	}

	// Collect all unique hashes.
	allHashes := make([]string, 0, len(allChunkData))
	for h := range allChunkData {
		allHashes = append(allHashes, h)
	}

	var totalBytes int64
	for _, fi := range fileInfos {
		totalBytes += int64(fi.manifest.TotalBytes)
	}
	atomic.StoreInt64(&u.stats.BytesRead, totalBytes)

	log.Printf("phase 1 complete: %d files, %d unique chunks, %.1f MiB hashed in %s",
		len(files), len(allHashes), float64(totalBytes)/(1<<20),
		time.Since(phaseStart).Round(time.Millisecond))

	// ── Phase 2: Batch probe all hashes ───────────────────────────────
	log.Printf("phase 2: probing %d unique hashes...", len(allHashes))
	phaseStart = time.Now()

	neededSet := make(map[string]struct{})
	batchSize := u.cfg.ProbeBatch
	if batchSize < 256 {
		batchSize = 256
	}

	// Probe in parallel batches.
	type probeResult struct {
		needed []string
		err    error
	}
	probeCh := make(chan []string, u.cfg.Concurrency)
	probeResults := make(chan probeResult, u.cfg.Concurrency)

	// Launch probe workers.
	probeWorkers := 4
	if probeWorkers > u.cfg.Concurrency {
		probeWorkers = u.cfg.Concurrency
	}
	var probeWg sync.WaitGroup
	for w := 0; w < probeWorkers; w++ {
		probeWg.Add(1)
		go func() {
			defer probeWg.Done()
			for batch := range probeCh {
				needed, err := u.probe(batch)
				probeResults <- probeResult{needed: needed, err: err}
			}
		}()
	}

	// Feed batches.
	go func() {
		for i := 0; i < len(allHashes); i += batchSize {
			end := i + batchSize
			if end > len(allHashes) {
				end = len(allHashes)
			}
			probeCh <- allHashes[i:end]
		}
		close(probeCh)
		probeWg.Wait()
		close(probeResults)
	}()

	// Collect probe results.
	var probeErr error
	for pr := range probeResults {
		if pr.err != nil {
			probeErr = pr.err
			continue
		}
		for _, h := range pr.needed {
			neededSet[h] = struct{}{}
		}
	}

	if probeErr != nil {
		// On probe failure, upload everything.
		log.Printf("probe had errors, uploading all chunks: %v", probeErr)
		for _, h := range allHashes {
			neededSet[h] = struct{}{}
		}
	}

	skipped := len(allHashes) - len(neededSet)
	atomic.StoreInt64(&u.stats.SkippedServer, int64(skipped))

	log.Printf("phase 2 complete: %d needed, %d already on server, in %s",
		len(neededSet), skipped, time.Since(phaseStart).Round(time.Millisecond))

	// ── Phase 3: Upload needed chunks concurrently ────────────────────
	log.Printf("phase 3: uploading %d chunks...", len(neededSet))
	phaseStart = time.Now()

	uploadCh := make(chan pendingChunk, u.cfg.Concurrency*2)
	var uploadWg sync.WaitGroup
	var uploadFailed int64

	for w := 0; w < u.cfg.Concurrency; w++ {
		uploadWg.Add(1)
		go func() {
			defer uploadWg.Done()
			for pc := range uploadCh {
				var lastErr error
				success := false
				for attempt := 0; attempt < maxUploadRetries; attempt++ {
					if attempt > 0 {
						delay := retryBaseDelay * time.Duration(1<<(attempt-1))
						time.Sleep(delay)
					}
					if err := u.uploadChunk(pc.meta.Hash, pc.data); err != nil {
						lastErr = err
						continue
					}
					success = true
					break
				}
				if success {
					atomic.AddInt64(&u.stats.UploadedNew, 1)
					atomic.AddInt64(&u.stats.BytesSent, int64(pc.meta.Size))
				} else {
					log.Printf("upload failed chunk %s after %d retries: %v",
						pc.meta.Hash[:12], maxUploadRetries, lastErr)
					atomic.AddInt64(&uploadFailed, 1)
				}
			}
		}()
	}

	// Feed unique needed chunks.
	for hash := range neededSet {
		data := allChunkData[hash]
		uploadCh <- pendingChunk{
			meta: protocol.ChunkMeta{Hash: hash, Size: len(data)},
			data: data,
		}
	}
	close(uploadCh)
	uploadWg.Wait()

	if uploadFailed > 0 {
		return fmt.Errorf("%d chunks failed to upload", uploadFailed)
	}

	// Free chunk data — no longer needed.
	allChunkData = nil

	log.Printf("phase 3 complete: %d chunks uploaded, %.1f MiB sent in %s",
		len(neededSet),
		float64(atomic.LoadInt64(&u.stats.BytesSent))/(1<<20),
		time.Since(phaseStart).Round(time.Millisecond))

	// ── Phase 4: Upload all manifests concurrently ────────────────────
	log.Printf("phase 4: uploading %d manifests...", len(fileInfos))
	phaseStart = time.Now()

	type manifestResult struct {
		idx        int
		manifestID string
		err        error
	}

	manifestCh := make(chan int, len(fileInfos))
	manifestResults := make(chan manifestResult, len(fileInfos))

	manifestWorkers := u.cfg.Concurrency
	if manifestWorkers > len(fileInfos) {
		manifestWorkers = len(fileInfos)
	}

	var mWg sync.WaitGroup
	for w := 0; w < manifestWorkers; w++ {
		mWg.Add(1)
		go func() {
			defer mWg.Done()
			for idx := range manifestCh {
				mid, err := u.uploadManifest(&fileInfos[idx].manifest)
				manifestResults <- manifestResult{idx: idx, manifestID: mid, err: err}
			}
		}()
	}

	for i := range fileInfos {
		manifestCh <- i
	}
	close(manifestCh)
	mWg.Wait()
	close(manifestResults)

	var results []FileResult
	var failed int
	for mr := range manifestResults {
		if mr.err != nil {
			log.Printf("FAILED manifest for %s: %v", fileInfos[mr.idx].relPath, mr.err)
			failed++
			continue
		}
		results = append(results, FileResult{
			Path:       fileInfos[mr.idx].relPath,
			ManifestID: mr.manifestID,
		})
	}

	log.Printf("phase 4 complete: %d manifests uploaded in %s",
		len(results), time.Since(phaseStart).Round(time.Millisecond))

	// Sort results by path for deterministic output.
	sort.Slice(results, func(i, j int) bool {
		return results[i].Path < results[j].Path
	})

	elapsed := time.Since(totalStart)

	// Print directory summary.
	sent := atomic.LoadInt64(&u.stats.BytesSent)
	log.Printf("--- directory upload complete ---")
	log.Printf("  files:         %d uploaded, %d failed", len(results), failed)
	log.Printf("  total read:    %.2f MiB", float64(totalBytes)/(1<<20))
	log.Printf("  total sent:    %.2f MiB", float64(sent)/(1<<20))
	log.Printf("  total time:    %s", elapsed.Round(time.Millisecond))
	if elapsed.Seconds() > 0 {
		log.Printf("  throughput:    %.2f MiB/s (read), %.2f MiB/s (sent)",
			float64(totalBytes)/(1<<20)/elapsed.Seconds(),
			float64(sent)/(1<<20)/elapsed.Seconds())
	}
	for _, r := range results {
		log.Printf("  manifest ID:   %s  (%s)", r.ManifestID, r.Path)
	}

	if failed > 0 {
		return fmt.Errorf("%d files failed to upload", failed)
	}
	return nil
}

// runFile uploads a single file/device and returns its manifest ID.
// This is the core upload logic extracted from the original Run().
func (u *Uploader) runFile(filePath string) (string, error) {
	// Initialize cache.
	var cached *cache.CachedManifest
	if !u.cfg.NoCache {
		c, err := cache.New(u.cfg.CacheDir)
		if err != nil {
			log.Printf("warning: cache init failed, continuing without cache: %v", err)
		} else {
			u.cache = c
			cached, err = c.Load(filePath)
			if err != nil {
				log.Printf("warning: cache load failed: %v", err)
			}
			if cached != nil {
				log.Printf("loaded cached manifest: %d chunks from previous run", len(cached.Hashes))
			} else {
				log.Printf("no cached manifest found (first run for this source)")
			}
		}
	}

	log.Printf("opening source: %s", filePath)

	pipeline, err := chunk.NewPipeline(filePath, u.cfg.Concurrency)
	if err != nil {
		return "", fmt.Errorf("create pipeline: %w", err)
	}

	deviceSize := pipeline.DeviceSize()
	u.stats.BytesTotal = int64(deviceSize)
	log.Printf("source size: %d bytes (%.2f GiB)", deviceSize, float64(deviceSize)/(1<<30))

	results := pipeline.Run()

	var manifest protocol.Manifest
	manifest.SourceDevice = filePath
	manifest.TotalBytes = deviceSize
	manifest.ChunkSize = protocol.ChunkSize

	var allMeta []protocol.ChunkMeta
	var metaMu sync.Mutex

	// Track chunks confirmed on the server (uploaded or server-skip).
	var confirmedHashes []uploadedHash
	var confirmedMu sync.Mutex

	// Upload worker pool with retry.
	uploadCh := make(chan pendingChunk, u.cfg.Concurrency*2)
	var uploadWg sync.WaitGroup
	for i := 0; i < u.cfg.Concurrency; i++ {
		uploadWg.Add(1)
		go func() {
			defer uploadWg.Done()
			for pc := range uploadCh {
				var lastErr error
				success := false
				for attempt := 0; attempt < maxUploadRetries; attempt++ {
					if attempt > 0 {
						delay := retryBaseDelay * time.Duration(1<<(attempt-1))
						time.Sleep(delay)
					}
					if err := u.uploadChunk(pc.meta.Hash, pc.data); err != nil {
						lastErr = err
						continue
					}
					success = true
					break
				}
				if success {
					atomic.AddInt64(&u.stats.UploadedNew, 1)
					atomic.AddInt64(&u.stats.BytesSent, int64(pc.meta.Size))
					confirmedMu.Lock()
					confirmedHashes = append(confirmedHashes, uploadedHash{index: pc.meta.Index, hash: pc.meta.Hash})
					confirmedMu.Unlock()
				} else {
					log.Printf("upload failed chunk %d (%s) after %d retries: %v",
						pc.meta.Index, pc.meta.Hash[:12], maxUploadRetries, lastErr)
					atomic.AddInt64(&u.stats.Failed, 1)
				}
			}
		}()
	}

	// Concurrent probe+upload pipeline: probe batches are dispatched to a
	// pool of goroutines so probing one batch doesn't block reading or
	// uploading of other batches.
	probeCh := make(chan []pendingChunk, u.cfg.Concurrency)
	var probeWg sync.WaitGroup
	probeWorkers := u.cfg.Concurrency
	if probeWorkers > 4 {
		probeWorkers = 4
	}
	for i := 0; i < probeWorkers; i++ {
		probeWg.Add(1)
		go func() {
			defer probeWg.Done()
			for batch := range probeCh {
				u.probeAndUpload(batch, uploadCh, &confirmedHashes, &confirmedMu)
			}
		}()
	}

	// Progress reporter.
	done := make(chan struct{})
	go func() {
		ticker := time.NewTicker(2 * time.Second)
		defer ticker.Stop()
		for {
			select {
			case <-ticker.C:
				u.printProgress()
			case <-done:
				return
			}
		}
	}()

	// Use smaller probe batches to keep the pipeline flowing.
	probeBatchSize := u.cfg.Concurrency * 2
	if probeBatchSize > u.cfg.ProbeBatch {
		probeBatchSize = u.cfg.ProbeBatch
	}

	// Track cache-skipped hashes for cache saving (but they still go through
	// the probe pipeline to verify the server has them).

	var batch []pendingChunk

	for result := range results {
		atomic.AddInt64(&u.stats.TotalChunks, 1)
		atomic.AddInt64(&u.stats.BytesRead, int64(result.Meta.Size))

		metaMu.Lock()
		allMeta = append(allMeta, result.Meta)
		metaMu.Unlock()

		// Cache check: if hash matches the cached hash for this index,
		// count it as a cache hit for stats, but still probe the server
		// in case the chunk was deleted server-side.
		if cached != nil {
			if cachedHash, ok := cached.Hashes[result.Meta.Index]; ok && cachedHash == result.Meta.Hash {
				atomic.AddInt64(&u.stats.SkippedCache, 1)
			}
		}

		batch = append(batch, pendingChunk{meta: result.Meta, data: result.Data})

		if len(batch) >= probeBatchSize {
			probeCh <- batch
			batch = nil
		}
	}

	// Flush remaining batch.
	if len(batch) > 0 {
		probeCh <- batch
	}

	// Wait for all probes to complete, then all uploads.
	close(probeCh)
	probeWg.Wait()
	close(uploadCh)
	uploadWg.Wait()
	close(done)

	// Sort chunks by index for correct ordering in the manifest.
	sort.Slice(allMeta, func(i, j int) bool {
		return allMeta[i].Index < allMeta[j].Index
	})
	manifest.Chunks = allMeta

	// Check for failures before saving.
	failed := atomic.LoadInt64(&u.stats.Failed)
	if failed > 0 {
		u.printFinalStats("")
		return "", fmt.Errorf("%d chunks failed to upload — manifest NOT saved (re-run to retry)", failed)
	}

	// Save updated cache with only confirmed hashes.
	if u.cache != nil {
		allConfirmed := make(map[uint64]string, len(confirmedHashes))
		for _, ch := range confirmedHashes {
			allConfirmed[ch.index] = ch.hash
		}
		if err := u.cache.SaveFromMap(filePath, &manifest, allConfirmed); err != nil {
			log.Printf("warning: failed to save cache: %v", err)
		} else {
			log.Printf("saved manifest cache (%d chunks)", len(allConfirmed))
		}
	}

	// Upload manifest to server.
	log.Println("uploading manifest...")
	manifestID, err := u.uploadManifest(&manifest)
	if err != nil {
		return "", fmt.Errorf("upload manifest: %w", err)
	}

	u.printFinalStats(manifestID)
	return manifestID, nil
}

// probeAndUpload probes the server for which chunks in the batch are needed,
// then sends only those to the upload channel. Chunks the server already has
// are tracked as confirmed.
func (u *Uploader) probeAndUpload(batch []pendingChunk, uploadCh chan<- pendingChunk, confirmed *[]uploadedHash, confirmedMu *sync.Mutex) {
	if len(batch) == 0 {
		return
	}

	hashes := make([]string, len(batch))
	for i, pc := range batch {
		hashes[i] = pc.meta.Hash
	}

	needed, err := u.probe(hashes)
	if err != nil {
		// Probe failed — upload everything (server will reject dupes).
		log.Printf("probe failed, uploading all %d chunks: %v", len(batch), err)
		for _, pc := range batch {
			uploadCh <- pc
		}
		return
	}

	needSet := make(map[string]struct{}, len(needed))
	for _, h := range needed {
		needSet[h] = struct{}{}
	}

	skipped := 0
	for _, pc := range batch {
		if _, need := needSet[pc.meta.Hash]; need {
			uploadCh <- pc
		} else {
			skipped++
			// Server already has this chunk — confirmed.
			confirmedMu.Lock()
			*confirmed = append(*confirmed, uploadedHash{index: pc.meta.Index, hash: pc.meta.Hash})
			confirmedMu.Unlock()
		}
	}
	atomic.AddInt64(&u.stats.SkippedServer, int64(skipped))
}

// probe sends a batch of hashes to the server and returns which ones are needed.
func (u *Uploader) probe(hashes []string) ([]string, error) {
	reqBody, err := json.Marshal(protocol.ProbeRequest{Hashes: hashes})
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest(http.MethodPost, u.cfg.ServerURL+"/api/probe", bytes.NewReader(reqBody))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	u.setAuth(req)

	resp, err := u.http.Do(req)
	if err != nil {
		return nil, fmt.Errorf("probe request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("probe returned %d: %s", resp.StatusCode, body)
	}

	var probeResp protocol.ProbeResponse
	if err := json.NewDecoder(resp.Body).Decode(&probeResp); err != nil {
		return nil, fmt.Errorf("decode probe response: %w", err)
	}
	return probeResp.Needed, nil
}

// uploadChunk sends a single chunk to the server.
func (u *Uploader) uploadChunk(hash string, data []byte) error {
	url := fmt.Sprintf("%s/api/chunk/%s", u.cfg.ServerURL, hash)
	req, err := http.NewRequest(http.MethodPut, url, bytes.NewReader(data))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/octet-stream")
	u.setAuth(req)

	resp, err := u.http.Do(req)
	if err != nil {
		return fmt.Errorf("upload request: %w", err)
	}
	defer resp.Body.Close()

	// Accept both 200 OK and 204 No Content as success.
	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusNoContent {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("upload returned %d: %s", resp.StatusCode, body)
	}
	// Drain body to allow connection reuse.
	io.Copy(io.Discard, resp.Body)

	return nil
}

// uploadManifest sends the completed manifest to the server.
func (u *Uploader) uploadManifest(m *protocol.Manifest) (string, error) {
	body, err := json.Marshal(m)
	if err != nil {
		return "", err
	}

	req, err := http.NewRequest(http.MethodPost, u.cfg.ServerURL+"/api/manifest", bytes.NewReader(body))
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", "application/json")
	u.setAuth(req)

	resp, err := u.http.Do(req)
	if err != nil {
		return "", fmt.Errorf("manifest request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		respBody, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("manifest returned %d: %s", resp.StatusCode, respBody)
	}

	var mResp protocol.ManifestUploadResponse
	if err := json.NewDecoder(resp.Body).Decode(&mResp); err != nil {
		return "", err
	}
	if mResp.Error != "" {
		return "", fmt.Errorf("server error: %s", mResp.Error)
	}
	return mResp.ID, nil
}

// printProgress logs current upload progress.
func (u *Uploader) printProgress() {
	total := atomic.LoadInt64(&u.stats.TotalChunks)
	uploaded := atomic.LoadInt64(&u.stats.UploadedNew)
	skippedCache := atomic.LoadInt64(&u.stats.SkippedCache)
	skippedServer := atomic.LoadInt64(&u.stats.SkippedServer)
	failed := atomic.LoadInt64(&u.stats.Failed)
	sent := atomic.LoadInt64(&u.stats.BytesSent)
	read := atomic.LoadInt64(&u.stats.BytesRead)
	elapsed := time.Since(u.stats.StartTime)

	processed := uploaded + skippedCache + skippedServer + failed
	pct := float64(0)
	if u.stats.BytesTotal > 0 {
		pct = float64(read) / float64(u.stats.BytesTotal) * 100
	}

	rate := float64(0)
	if elapsed.Seconds() > 0 {
		rate = float64(read) / (1 << 20) / elapsed.Seconds()
	}

	log.Printf("progress: %d/%d chunks (%.1f%%) | read %.1f MiB/s | sent %.1f MiB | cache-skip %d | server-skip %d | uploaded %d | failed %d",
		processed, total, pct, rate, float64(sent)/(1<<20), skippedCache, skippedServer, uploaded, failed)
}

// printFinalStats logs the final summary.
func (u *Uploader) printFinalStats(manifestID string) {
	elapsed := time.Since(u.stats.StartTime)
	total := atomic.LoadInt64(&u.stats.TotalChunks)
	uploaded := atomic.LoadInt64(&u.stats.UploadedNew)
	skippedCache := atomic.LoadInt64(&u.stats.SkippedCache)
	skippedServer := atomic.LoadInt64(&u.stats.SkippedServer)
	failed := atomic.LoadInt64(&u.stats.Failed)
	sent := atomic.LoadInt64(&u.stats.BytesSent)
	read := atomic.LoadInt64(&u.stats.BytesRead)

	log.Printf("--- complete ---")
	log.Printf("  chunks:        %d total", total)
	log.Printf("  cache hit:     %d (hash unchanged since last run)", skippedCache)
	log.Printf("  server skip:   %d (server already had)", skippedServer)
	log.Printf("  uploaded:      %d new chunks", uploaded)
	log.Printf("  failed:        %d", failed)
	log.Printf("  read:          %.2f MiB (%.2f MiB/s)", float64(read)/(1<<20), float64(read)/(1<<20)/elapsed.Seconds())
	log.Printf("  transferred:   %.2f MiB (%.2f MiB/s)", float64(sent)/(1<<20), float64(sent)/(1<<20)/elapsed.Seconds())
	log.Printf("  time:          %s", elapsed.Round(time.Millisecond))
	log.Printf("  manifest ID:   %s", manifestID)

	if total > 0 {
		skipped := total - uploaded - failed
		savingsPct := float64(skipped) / float64(total) * 100
		log.Printf("  dedup savings: %.1f%% of chunks skipped", savingsPct)
	}
}
