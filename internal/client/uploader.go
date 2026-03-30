// Package client implements the concurrent chunk upload client.
package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"sync"
	"sync/atomic"
	"time"

	"github.com/theo/synche2/internal/cache"
	"github.com/theo/synche2/internal/chunk"
	"github.com/theo/synche2/internal/protocol"
)

// Config holds client configuration.
type Config struct {
	ServerURL   string // e.g. "http://localhost:8420"
	DevicePath  string // e.g. "/dev/sda" or a file path
	Concurrency int    // number of parallel upload workers
	ProbeBatch  int    // how many hashes to send per probe request
	CacheDir    string // local manifest cache dir (empty = ~/.cache/synche)
	NoCache     bool   // disable local caching
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
func NewUploader(cfg Config) *Uploader {
	if cfg.Concurrency < 1 {
		cfg.Concurrency = 8
	}
	if cfg.ProbeBatch < 1 {
		cfg.ProbeBatch = 256
	}

	return &Uploader{
		cfg: cfg,
		http: &http.Client{
			Timeout: 120 * time.Second,
			Transport: &http.Transport{
				MaxIdleConns:        cfg.Concurrency * 2,
				MaxIdleConnsPerHost: cfg.Concurrency * 2,
				MaxConnsPerHost:     cfg.Concurrency * 2,
				IdleConnTimeout:     90 * time.Second,
			},
		},
		stats: Stats{StartTime: time.Now()},
	}
}

type pendingChunk struct {
	meta protocol.ChunkMeta
	data []byte
}

// Run executes the full upload pipeline:
// 1. Load cached manifest (if any)
// 2. Read raw blocks from device + hash with BLAKE3
// 3. Compare each chunk hash against cache — skip if unchanged
// 4. Probe server in small concurrent batches, upload needed chunks
// 5. Save updated manifest to cache
// 6. Upload the final manifest to server
func (u *Uploader) Run() error {
	// Initialize cache.
	var cached *cache.CachedManifest
	if !u.cfg.NoCache {
		c, err := cache.New(u.cfg.CacheDir)
		if err != nil {
			log.Printf("warning: cache init failed, continuing without cache: %v", err)
		} else {
			u.cache = c
			cached, err = c.Load(u.cfg.DevicePath)
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

	log.Printf("opening source: %s", u.cfg.DevicePath)

	pipeline, err := chunk.NewPipeline(u.cfg.DevicePath, u.cfg.Concurrency)
	if err != nil {
		return fmt.Errorf("create pipeline: %w", err)
	}

	deviceSize := pipeline.DeviceSize()
	u.stats.BytesTotal = int64(deviceSize)
	log.Printf("source size: %d bytes (%.2f GiB)", deviceSize, float64(deviceSize)/(1<<30))

	results := pipeline.Run()

	var manifest protocol.Manifest
	manifest.SourceDevice = u.cfg.DevicePath
	manifest.TotalBytes = deviceSize
	manifest.ChunkSize = protocol.ChunkSize

	var allMeta []protocol.ChunkMeta
	var metaMu sync.Mutex

	// Upload worker pool.
	uploadCh := make(chan pendingChunk, u.cfg.Concurrency*2)
	var uploadWg sync.WaitGroup
	for i := 0; i < u.cfg.Concurrency; i++ {
		uploadWg.Add(1)
		go func() {
			defer uploadWg.Done()
			for pc := range uploadCh {
				if err := u.uploadChunk(pc.meta.Hash, pc.data); err != nil {
					log.Printf("upload error chunk %d (%s): %v", pc.meta.Index, pc.meta.Hash[:12], err)
					atomic.AddInt64(&u.stats.Failed, 1)
				} else {
					atomic.AddInt64(&u.stats.UploadedNew, 1)
					atomic.AddInt64(&u.stats.BytesSent, int64(pc.meta.Size))
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
				u.probeAndUpload(batch, uploadCh)
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

	var batch []pendingChunk

	for result := range results {
		atomic.AddInt64(&u.stats.TotalChunks, 1)
		atomic.AddInt64(&u.stats.BytesRead, int64(result.Meta.Size))

		metaMu.Lock()
		allMeta = append(allMeta, result.Meta)
		metaMu.Unlock()

		// Cache check: skip unchanged chunks entirely.
		if cached != nil {
			if cachedHash, ok := cached.Hashes[result.Meta.Index]; ok && cachedHash == result.Meta.Hash {
				atomic.AddInt64(&u.stats.SkippedCache, 1)
				continue
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

	manifest.Chunks = allMeta

	// Save updated cache.
	if u.cache != nil {
		if err := u.cache.Save(u.cfg.DevicePath, &manifest); err != nil {
			log.Printf("warning: failed to save cache: %v", err)
		} else {
			log.Printf("saved manifest cache (%d chunks)", len(manifest.Chunks))
		}
	}

	// Upload manifest to server.
	log.Println("uploading manifest...")
	manifestID, err := u.uploadManifest(&manifest)
	if err != nil {
		return fmt.Errorf("upload manifest: %w", err)
	}

	u.printFinalStats(manifestID)
	return nil
}

// probeAndUpload probes the server for which chunks in the batch are needed,
// then sends only those to the upload channel.
func (u *Uploader) probeAndUpload(batch []pendingChunk, uploadCh chan<- pendingChunk) {
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

	resp, err := u.http.Post(u.cfg.ServerURL+"/api/probe", "application/json", bytes.NewReader(reqBody))
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

	resp, err := u.http.Post(u.cfg.ServerURL+"/api/manifest", "application/json", bytes.NewReader(body))
	if err != nil {
		return "", fmt.Errorf("manifest request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("manifest returned %d: %s", resp.StatusCode, body)
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
	log.Printf("  cache skip:    %d (unchanged since last run)", skippedCache)
	log.Printf("  server skip:   %d (server already had)", skippedServer)
	log.Printf("  uploaded:      %d new chunks", uploaded)
	log.Printf("  failed:        %d", failed)
	log.Printf("  read:          %.2f MiB (%.2f MiB/s)", float64(read)/(1<<20), float64(read)/(1<<20)/elapsed.Seconds())
	log.Printf("  transferred:   %.2f MiB (%.2f MiB/s)", float64(sent)/(1<<20), float64(sent)/(1<<20)/elapsed.Seconds())
	log.Printf("  time:          %s", elapsed.Round(time.Millisecond))
	log.Printf("  manifest ID:   %s", manifestID)

	if total > 0 {
		savingsPct := float64(skippedCache+skippedServer) / float64(total) * 100
		log.Printf("  dedup savings: %.1f%% of chunks skipped", savingsPct)
	}
}
