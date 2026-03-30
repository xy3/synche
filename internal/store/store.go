// Package store implements content-addressable chunk storage on the server side.
// Chunks are stored in a two-level directory structure based on hash prefix
// to avoid massive flat directories.
package store

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/theo/synche2/internal/protocol"
)

// Store is a content-addressable chunk store backed by the filesystem.
type Store struct {
	root string

	// In-memory set of known hashes for fast existence checks.
	mu     sync.RWMutex
	hashes map[string]struct{}
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
		root:   root,
		hashes: make(map[string]struct{}),
	}

	// Scan existing chunks into memory for fast lookups.
	if err := s.loadIndex(); err != nil {
		return nil, fmt.Errorf("load index: %w", err)
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
	// Fast path: already have it.
	if s.Has(hash) {
		return false, nil
	}

	path := s.chunkPath(hash)

	// Double-check with lock held to avoid races.
	s.mu.Lock()
	if _, ok := s.hashes[hash]; ok {
		s.mu.Unlock()
		return false, nil
	}
	s.hashes[hash] = struct{}{}
	s.mu.Unlock()

	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return false, fmt.Errorf("mkdir for chunk: %w", err)
	}

	// Write atomically: write to tmp, then rename.
	tmp := path + ".tmp"
	if err := os.WriteFile(tmp, data, 0o644); err != nil {
		// Roll back the in-memory state.
		s.mu.Lock()
		delete(s.hashes, hash)
		s.mu.Unlock()
		return false, fmt.Errorf("write chunk: %w", err)
	}
	if err := os.Rename(tmp, path); err != nil {
		s.mu.Lock()
		delete(s.hashes, hash)
		s.mu.Unlock()
		return false, fmt.Errorf("rename chunk: %w", err)
	}

	return true, nil
}

// Get retrieves chunk data by hash. Returns os.ErrNotExist if not found.
func (s *Store) Get(hash string) ([]byte, error) {
	path := s.chunkPath(hash)
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	return data, nil
}

// SaveManifest persists a manifest and returns a generated ID.
func (s *Store) SaveManifest(m *protocol.Manifest) (string, error) {
	id := fmt.Sprintf("%s_%d", filepath.Base(m.SourceDevice), time.Now().UnixNano())
	path := filepath.Join(s.root, "manifests", id+".json")

	data, err := json.MarshalIndent(m, "", "  ")
	if err != nil {
		return "", err
	}
	if err := os.WriteFile(path, data, 0o644); err != nil {
		return "", err
	}
	return id, nil
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

// LoadManifest reads a manifest by ID.
func (s *Store) LoadManifest(id string) (*protocol.Manifest, error) {
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
	TotalChunks int    `json:"total_chunks"`
	StorePath   string `json:"store_path"`
}

// Stats returns store statistics.
func (s *Store) Stats() Stats {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return Stats{
		TotalChunks: len(s.hashes),
		StorePath:   s.root,
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

// Reset wipes all chunks and manifests from the store.
func (s *Store) Reset() error {
	s.mu.Lock()
	defer s.mu.Unlock()

	chunkDir := filepath.Join(s.root, "chunks")
	manifestDir := filepath.Join(s.root, "manifests")

	if err := os.RemoveAll(chunkDir); err != nil {
		return fmt.Errorf("remove chunks: %w", err)
	}
	if err := os.RemoveAll(manifestDir); err != nil {
		return fmt.Errorf("remove manifests: %w", err)
	}
	if err := os.MkdirAll(chunkDir, 0o755); err != nil {
		return fmt.Errorf("recreate chunks dir: %w", err)
	}
	if err := os.MkdirAll(manifestDir, 0o755); err != nil {
		return fmt.Errorf("recreate manifests dir: %w", err)
	}

	s.hashes = make(map[string]struct{})
	return nil
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
