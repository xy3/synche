// Package protocol defines the shared types and constants between client and server.
package protocol

const (
	// ChunkSize is the fixed chunk size: 1 MiB.
	ChunkSize = 1 << 20 // 1,048,576 bytes

	// HashSize is the length of a BLAKE3-256 hash in bytes.
	HashSize = 32
)

// ChunkMeta describes a single chunk in a transfer manifest.
type ChunkMeta struct {
	Index  uint64 `json:"index"`            // positional index on the source device
	Hash   string `json:"hash"`             // hex-encoded BLAKE3 hash
	Size   int    `json:"size"`             // actual byte count (last chunk may be < ChunkSize)
	Offset uint64 `json:"offset,omitempty"` // byte offset on source device
}

// Manifest describes an entire block device snapshot.
type Manifest struct {
	SourceDevice string      `json:"source_device"`
	TotalBytes   uint64      `json:"total_bytes"`
	ChunkSize    int         `json:"chunk_size"`
	Chunks       []ChunkMeta `json:"chunks"`
}

// ProbeRequest is sent by the client to ask which chunks the server already has.
type ProbeRequest struct {
	Hashes []string `json:"hashes"`
}

// ProbeResponse tells the client which hashes are needed (not already stored).
type ProbeResponse struct {
	Needed []string `json:"needed"` // hashes the server does NOT have
}

// UploadResponse is returned after a chunk upload attempt.
type UploadResponse struct {
	Hash     string `json:"hash"`
	Accepted bool   `json:"accepted"` // false = server already had it (rejected duplicate)
	Error    string `json:"error,omitempty"`
}

// ManifestUploadResponse is returned after storing a manifest.
type ManifestUploadResponse struct {
	ID    string `json:"id"` // server-assigned manifest ID
	Error string `json:"error,omitempty"`
}
