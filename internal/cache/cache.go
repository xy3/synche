// Package cache implements local manifest caching so that subsequent syncs
// can skip chunks that haven't changed since the last run.
//
// The cache key is derived from the absolute path of the source. For each
// source, the cache stores the previous manifest (index → hash mapping).
// On the next run, after hashing a chunk, we compare its hash to the cached
// hash at the same index. If they match, the chunk is unchanged — we skip
// both the server probe and the upload entirely.
package cache

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/theo/synche2/internal/protocol"
)

// Cache provides local manifest caching.
type Cache struct {
	dir string
}

// New creates a cache rooted at the given directory.
// If dir is empty, defaults to ~/.cache/synche.
func New(dir string) (*Cache, error) {
	if dir == "" {
		home, err := os.UserHomeDir()
		if err != nil {
			return nil, fmt.Errorf("get home dir: %w", err)
		}
		dir = filepath.Join(home, ".cache", "synche")
	}
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return nil, fmt.Errorf("create cache dir: %w", err)
	}
	return &Cache{dir: dir}, nil
}

// cacheKey generates a deterministic filename from a source path.
func cacheKey(sourcePath string) string {
	abs, err := filepath.Abs(sourcePath)
	if err != nil {
		abs = sourcePath
	}
	h := sha256.Sum256([]byte(abs))
	return hex.EncodeToString(h[:16]) // 128 bits is plenty for a cache key
}

// CachedManifest is the on-disk cache format.
type CachedManifest struct {
	SourcePath string            `json:"source_path"`
	TotalBytes uint64            `json:"total_bytes"`
	ChunkSize  int               `json:"chunk_size"`
	Hashes     map[uint64]string `json:"hashes"` // index → hash
}

// Load reads the cached manifest for the given source path.
// Returns nil (no error) if no cache exists.
func (c *Cache) Load(sourcePath string) (*CachedManifest, error) {
	path := filepath.Join(c.dir, cacheKey(sourcePath)+".json")
	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, fmt.Errorf("read cache: %w", err)
	}

	var cm CachedManifest
	if err := json.Unmarshal(data, &cm); err != nil {
		// Corrupt cache, just ignore it.
		return nil, nil
	}
	return &cm, nil
}

// Save writes a manifest to the cache for the given source path.
func (c *Cache) Save(sourcePath string, manifest *protocol.Manifest) error {
	cm := CachedManifest{
		SourcePath: sourcePath,
		TotalBytes: manifest.TotalBytes,
		ChunkSize:  manifest.ChunkSize,
		Hashes:     make(map[uint64]string, len(manifest.Chunks)),
	}
	for _, chunk := range manifest.Chunks {
		cm.Hashes[chunk.Index] = chunk.Hash
	}

	data, err := json.Marshal(&cm)
	if err != nil {
		return fmt.Errorf("marshal cache: %w", err)
	}

	path := filepath.Join(c.dir, cacheKey(sourcePath)+".json")
	tmp := path + ".tmp"
	if err := os.WriteFile(tmp, data, 0o644); err != nil {
		return fmt.Errorf("write cache: %w", err)
	}
	if err := os.Rename(tmp, path); err != nil {
		return fmt.Errorf("rename cache: %w", err)
	}
	return nil
}

// Dir returns the cache directory path.
func (c *Cache) Dir() string {
	return c.dir
}
