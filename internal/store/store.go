// Package store implements content-addressable chunk storage on the server side.
// Chunks are stored in a two-level directory structure based on hash prefix
// to avoid massive flat directories.
package store

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/theo/synche2/internal/protocol"
	"github.com/zeebo/blake3"
)

// ValidHash returns true if hash is a valid hex-encoded BLAKE3-256 hash.
// This MUST be checked before using a hash in any filesystem path to prevent
// path traversal attacks.
func ValidHash(hash string) bool {
	if len(hash) != protocol.HashSize*2 {
		return false
	}
	_, err := hex.DecodeString(hash)
	return err == nil
}

// Store is a content-addressable chunk store backed by the filesystem.
type Store struct {
	root string

	// In-memory set of known hashes for fast existence checks.
	mu     sync.RWMutex
	hashes map[string]struct{}

	// chunkRefs tracks how many manifests reference each chunk hash.
	// Maintained in memory so that orphan detection during delete is O(1)
	// per chunk instead of requiring a full manifest scan.
	chunkRefs map[string]int
}

// New creates or opens a chunk store at the given root directory.
func New(root string) (*Store, error) {
	chunkDir := filepath.Join(root, "chunks")
	manifestDir := filepath.Join(root, "manifests")

	for _, d := range []string{chunkDir, manifestDir} {
		if err := os.MkdirAll(d, 0o755); err != nil {
			return nil, fmt.Errorf("create store dir %s: %w", d, err)
		}
	}

	s := &Store{
		root:      root,
		hashes:    make(map[string]struct{}),
		chunkRefs: make(map[string]int),
	}

	// Scan existing chunks into memory for fast lookups.
	if err := s.loadIndex(); err != nil {
		return nil, fmt.Errorf("load index: %w", err)
	}

	// Build chunk reference counts from manifests.
	if err := s.loadChunkRefs(); err != nil {
		return nil, fmt.Errorf("load chunk refs: %w", err)
	}

	return s, nil
}

// Has returns true if a chunk with the given hash already exists.
func (s *Store) Has(hash string) bool {
	s.mu.RLock()
	defer s.mu.RUnlock()
	_, ok := s.hashes[hash]
	return ok
}

// HasBatch checks multiple hashes at once and returns the set of hashes
// that the store does NOT have (i.e., the ones the client needs to upload).
func (s *Store) HasBatch(hashes []string) []string {
	s.mu.RLock()
	defer s.mu.RUnlock()

	var needed []string
	for _, h := range hashes {
		if _, ok := s.hashes[h]; !ok {
			needed = append(needed, h)
		}
	}
	return needed
}

// Put stores a chunk. Returns (accepted=true) if the chunk was new,
// (accepted=false) if it was a duplicate (no error in either case).
func (s *Store) Put(hash string, data []byte) (bool, error) {
	if !ValidHash(hash) {
		return false, fmt.Errorf("invalid hash: %q", hash)
	}

	// Fast path: already have it.
	if s.Has(hash) {
		return false, nil
	}

	path := s.chunkPath(hash)

	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return false, fmt.Errorf("mkdir for chunk: %w", err)
	}

	// Write atomically: write to tmp, then rename.
	tmp := path + ".tmp"
	if err := os.WriteFile(tmp, data, 0o644); err != nil {
		return false, fmt.Errorf("write chunk: %w", err)
	}
	if err := os.Rename(tmp, path); err != nil {
		os.Remove(tmp) // best-effort cleanup
		return false, fmt.Errorf("rename chunk: %w", err)
	}

	// Mark present in memory AFTER successful disk write.
	s.mu.Lock()
	s.hashes[hash] = struct{}{}
	s.mu.Unlock()

	return true, nil
}

// Get retrieves chunk data by hash. Returns os.ErrNotExist if not found.
func (s *Store) Get(hash string) ([]byte, error) {
	if !ValidHash(hash) {
		return nil, fmt.Errorf("invalid hash: %q", hash)
	}
	path := s.chunkPath(hash)
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	return data, nil
}

// manifestContentID generates a deterministic manifest ID from the source
// filename and the ordered list of chunk hashes. Re-uploading the same file
// with the same content always produces the same ID.
func manifestContentID(m *protocol.Manifest) string {
	h := blake3.New()
	h.Write([]byte(filepath.Base(m.SourceDevice)))
	for _, c := range m.Chunks {
		h.Write([]byte(c.Hash))
	}
	sum := hex.EncodeToString(h.Sum(nil))
	return fmt.Sprintf("%s_%s", filepath.Base(m.SourceDevice), sum[:16])
}

// SaveManifest persists a manifest and returns a generated ID.
// The ID is content-based: re-uploading the same file with identical chunks
// returns the existing manifest ID without creating a duplicate.
// The second return value is true if the manifest was newly created.
func (s *Store) SaveManifest(m *protocol.Manifest) (string, bool, error) {
	id := manifestContentID(m)
	path := filepath.Join(s.root, "manifests", id+".json")

	// If this exact manifest already exists, return the existing ID.
	if _, err := os.Stat(path); err == nil {
		return id, false, nil
	}

	data, err := json.MarshalIndent(m, "", "  ")
	if err != nil {
		return "", false, err
	}
	if err := os.WriteFile(path, data, 0o644); err != nil {
		return "", false, err
	}

	// Update in-memory chunk reference counts.
	s.mu.Lock()
	for _, c := range m.Chunks {
		s.chunkRefs[c.Hash]++
	}
	s.mu.Unlock()

	return id, true, nil
}

// ManifestInfo holds metadata about a stored manifest.
type ManifestInfo struct {
	ID         string    // manifest ID (filename without .json)
	Name       string    // source device/file basename
	TotalBytes uint64    // total size of the original file
	NumChunks  int       // number of chunks
	ModTime    time.Time // when the manifest was saved
}

// ListManifests returns info about all stored manifests.
func (s *Store) ListManifests() ([]ManifestInfo, error) {
	manifestDir := filepath.Join(s.root, "manifests")
	entries, err := os.ReadDir(manifestDir)
	if err != nil {
		return nil, err
	}

	var results []ManifestInfo
	for _, e := range entries {
		if e.IsDir() || filepath.Ext(e.Name()) != ".json" {
			continue
		}
		id := strings.TrimSuffix(e.Name(), ".json")
		m, err := s.LoadManifest(id)
		if err != nil {
			continue
		}
		info, _ := e.Info()
		modTime := time.Time{}
		if info != nil {
			modTime = info.ModTime()
		}
		results = append(results, ManifestInfo{
			ID:         id,
			Name:       filepath.Base(m.SourceDevice),
			TotalBytes: m.TotalBytes,
			NumChunks:  len(m.Chunks),
			ModTime:    modTime,
		})
	}
	return results, nil
}

// ValidManifestID returns true if the ID is safe for use in a filepath
// (no path separators, no .., not empty).
func ValidManifestID(id string) bool {
	if id == "" || id == "." || id == ".." {
		return false
	}
	if strings.ContainsAny(id, "/\\") {
		return false
	}
	return true
}

// LoadManifest reads a manifest by ID.
func (s *Store) LoadManifest(id string) (*protocol.Manifest, error) {
	if !ValidManifestID(id) {
		return nil, fmt.Errorf("invalid manifest ID: %q", id)
	}
	path := filepath.Join(s.root, "manifests", id+".json")
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var m protocol.Manifest
	if err := json.Unmarshal(data, &m); err != nil {
		return nil, err
	}
	return &m, nil
}

// ReadChunk reads a chunk and writes it to the given writer. This is used
// for streaming reassembly without loading the entire chunk into memory.
func (s *Store) ReadChunk(hash string) ([]byte, error) {
	return s.Get(hash)
}

// Stats returns basic statistics about the store.
type Stats struct {
	TotalChunks int `json:"total_chunks"`
}

// Stats returns store statistics.
func (s *Store) Stats() Stats {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return Stats{
		TotalChunks: len(s.hashes),
	}
}

// chunkPath returns the filesystem path for a given chunk hash.
// Uses two-level directory: <root>/chunks/ab/cd/<full-hash>
func (s *Store) chunkPath(hash string) string {
	if len(hash) < 4 {
		return filepath.Join(s.root, "chunks", hash)
	}
	return filepath.Join(s.root, "chunks", hash[:2], hash[2:4], hash)
}

// DeleteManifest removes a manifest by ID. The chunks it referenced are NOT
// deleted because they may be shared by other manifests.
func (s *Store) DeleteManifest(id string) error {
	if !ValidManifestID(id) {
		return fmt.Errorf("invalid manifest ID: %q", id)
	}
	path := filepath.Join(s.root, "manifests", id+".json")
	if err := os.Remove(path); err != nil {
		return fmt.Errorf("remove manifest: %w", err)
	}
	return nil
}

// DeleteManifestAndOrphanedChunks removes a manifest and any chunks that are
// no longer referenced by any remaining manifest. Returns the number of
// orphaned chunks deleted.
//
// Uses in-memory chunk reference counts so this is O(chunks_in_manifest)
// rather than requiring a scan of all manifests.
func (s *Store) DeleteManifestAndOrphanedChunks(id string) (int, error) {
	if !ValidManifestID(id) {
		return 0, fmt.Errorf("invalid manifest ID: %q", id)
	}

	// Load the target manifest to know which chunks it references.
	target, err := s.LoadManifest(id)
	if err != nil {
		return 0, fmt.Errorf("load manifest to delete: %w", err)
	}

	// Delete the manifest file first.
	manifestPath := filepath.Join(s.root, "manifests", id+".json")
	if err := os.Remove(manifestPath); err != nil {
		return 0, fmt.Errorf("remove manifest: %w", err)
	}

	// Decrement reference counts and delete orphaned chunks.
	deleted := 0
	s.mu.Lock()
	defer s.mu.Unlock()

	for _, c := range target.Chunks {
		s.chunkRefs[c.Hash]--
		if s.chunkRefs[c.Hash] <= 0 {
			// No more manifests reference this chunk — delete it.
			delete(s.chunkRefs, c.Hash)
			path := s.chunkPath(c.Hash)
			if err := os.Remove(path); err == nil {
				delete(s.hashes, c.Hash)
				deleted++
			}
		}
	}

	return deleted, nil
}

// loadIndex scans the chunk directory and populates the in-memory hash set.
func (s *Store) loadIndex() error {
	chunkDir := filepath.Join(s.root, "chunks")
	return filepath.Walk(chunkDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return nil // skip errors
		}
		if info.IsDir() {
			return nil
		}
		hash := info.Name()
		// Skip tmp files.
		if filepath.Ext(hash) == ".tmp" {
			return nil
		}
		s.mu.Lock()
		s.hashes[hash] = struct{}{}
		s.mu.Unlock()
		return nil
	})
}

// loadChunkRefs scans all manifests and builds the in-memory chunk reference
// count map. Called once at startup.
func (s *Store) loadChunkRefs() error {
	manifestDir := filepath.Join(s.root, "manifests")
	entries, err := os.ReadDir(manifestDir)
	if err != nil {
		return nil // no manifests dir is fine
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	for _, e := range entries {
		if e.IsDir() || filepath.Ext(e.Name()) != ".json" {
			continue
		}
		id := strings.TrimSuffix(e.Name(), ".json")
		m, err := s.LoadManifest(id)
		if err != nil {
			continue
		}
		for _, c := range m.Chunks {
			s.chunkRefs[c.Hash]++
		}
	}
	return nil
}
