// Package server implements the HTTP chunk storage server.
package server

import (
	"encoding/binary"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/theo/synche2/internal/protocol"
	"github.com/theo/synche2/internal/store"
	"github.com/theo/synche2/internal/webdav"
	"github.com/zeebo/blake3"
)

const (
	// maxProbeBodySize limits the probe request body (50 MiB should cover ~500k hashes).
	maxProbeBodySize = 50 << 20
	// maxManifestBodySize limits manifest uploads (256 MiB covers multi-TB devices).
	maxManifestBodySize = 256 << 20
	// maxBatchBodySize limits batch upload requests (1 GiB).
	maxBatchBodySize = 1 << 30
)

// Server is the chunk storage HTTP server.
type Server struct {
	store *store.Store
	mux   *http.ServeMux
	addr  string
}

// New creates a new server that listens on addr and stores chunks in storePath.
func New(addr, storePath string) (*Server, error) {
	s, err := store.New(storePath)
	if err != nil {
		return nil, fmt.Errorf("open store: %w", err)
	}

	srv := &Server{
		store: s,
		mux:   http.NewServeMux(),
		addr:  addr,
	}
	srv.routes()
	return srv, nil
}

// ListenAndServe starts the HTTP server.
func (s *Server) ListenAndServe() error {
	log.Printf("synche server listening on %s", s.addr)
	stats := s.store.Stats()
	log.Printf("store: %d existing chunks", stats.TotalChunks)
	log.Printf("webdav available at http://%s/webdav/", s.addr)
	return http.ListenAndServe(s.addr, s.mux)
}

func (s *Server) routes() {
	s.mux.HandleFunc("POST /api/probe", s.handleProbe)
	s.mux.HandleFunc("PUT /api/chunk/{hash}", s.handleUploadChunk)
	s.mux.HandleFunc("GET /api/chunk/{hash}", s.handleGetChunk)
	s.mux.HandleFunc("POST /api/upload-batch", s.handleBatchUpload)
	s.mux.HandleFunc("POST /api/manifest", s.handleUploadManifest)
	s.mux.HandleFunc("POST /api/reset", s.handleReset)
	s.mux.HandleFunc("GET /api/stats", s.handleStats)
	s.mux.HandleFunc("GET /api/health", s.handleHealth)

	// WebDAV: read-only virtual filesystem of uploaded files.
	dav := webdav.New(s.store, "/webdav")
	s.mux.Handle("/webdav/", dav)
	s.mux.Handle("/webdav", http.RedirectHandler("/webdav/", http.StatusMovedPermanently))
}

// handleProbe accepts a list of hashes and returns which ones the server needs.
func (s *Server) handleProbe(w http.ResponseWriter, r *http.Request) {
	r.Body = http.MaxBytesReader(w, r.Body, maxProbeBodySize)
	var req protocol.ProbeRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	needed := s.store.HasBatch(req.Hashes)
	resp := protocol.ProbeResponse{Needed: needed}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

// handleUploadChunk receives raw chunk data and stores it.
// The hash is provided in the URL. The server verifies it matches the data.
func (s *Server) handleUploadChunk(w http.ResponseWriter, r *http.Request) {
	hash := r.PathValue("hash")
	if !store.ValidHash(hash) {
		http.Error(w, "invalid hash", http.StatusBadRequest)
		return
	}

	// Fast reject: already have it.
	if s.store.Has(hash) {
		w.WriteHeader(http.StatusNoContent)
		return
	}

	data, err := io.ReadAll(io.LimitReader(r.Body, int64(protocol.ChunkSize)+1))
	if err != nil {
		http.Error(w, "read body: "+err.Error(), http.StatusInternalServerError)
		return
	}
	if len(data) > protocol.ChunkSize {
		http.Error(w, "chunk too large", http.StatusRequestEntityTooLarge)
		return
	}

	// Verify hash matches data.
	hasher := blake3.New()
	hasher.Write(data)
	computed := hex.EncodeToString(hasher.Sum(nil))
	if computed != hash {
		http.Error(w, fmt.Sprintf("hash mismatch: expected %s, got %s", hash, computed), http.StatusBadRequest)
		return
	}

	if _, err := s.store.Put(hash, data); err != nil {
		http.Error(w, "store chunk: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// handleGetChunk retrieves a chunk by hash.
func (s *Server) handleGetChunk(w http.ResponseWriter, r *http.Request) {
	hash := r.PathValue("hash")
	if !store.ValidHash(hash) {
		http.Error(w, "invalid hash", http.StatusBadRequest)
		return
	}
	data, err := s.store.Get(hash)
	if err != nil {
		http.Error(w, "not found", http.StatusNotFound)
		return
	}
	w.Header().Set("Content-Type", "application/octet-stream")
	w.Write(data)
}

// handleUploadManifest stores a manifest describing a complete device snapshot.
func (s *Server) handleUploadManifest(w http.ResponseWriter, r *http.Request) {
	r.Body = http.MaxBytesReader(w, r.Body, maxManifestBodySize)
	var manifest protocol.Manifest
	if err := json.NewDecoder(r.Body).Decode(&manifest); err != nil {
		http.Error(w, "invalid manifest", http.StatusBadRequest)
		return
	}

	id, err := s.store.SaveManifest(&manifest)
	if err != nil {
		http.Error(w, "save manifest: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(protocol.ManifestUploadResponse{ID: id})
}

// handleStats returns store statistics.
func (s *Server) handleStats(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(s.store.Stats())
}

// handleHealth is a simple health check endpoint.
func (s *Server) handleHealth(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(`{"status":"ok"}`))
}

// handleBatchUpload receives multiple chunks in a single request.
// Binary wire format: [64-byte hex hash][4-byte big-endian size][data bytes]...
// The server checks each chunk, stores new ones, skips duplicates.
func (s *Server) handleBatchUpload(w http.ResponseWriter, r *http.Request) {
	r.Body = http.MaxBytesReader(w, r.Body, maxBatchBodySize)
	var resp protocol.BatchUploadResponse
	hashBuf := make([]byte, protocol.HashSize*2) // 64 hex chars
	sizeBuf := make([]byte, 4)
	hasher := blake3.New()

	for {
		// Read hash (64 hex bytes).
		if _, err := io.ReadFull(r.Body, hashBuf); err != nil {
			if err == io.EOF || err == io.ErrUnexpectedEOF {
				break // end of stream
			}
			http.Error(w, "read hash: "+err.Error(), http.StatusBadRequest)
			return
		}
		hash := string(hashBuf)

		// Read size (4 bytes big-endian).
		if _, err := io.ReadFull(r.Body, sizeBuf); err != nil {
			http.Error(w, "read size: "+err.Error(), http.StatusBadRequest)
			return
		}
		size := binary.BigEndian.Uint32(sizeBuf)
		if size > uint32(protocol.ChunkSize) {
			http.Error(w, "chunk too large", http.StatusRequestEntityTooLarge)
			return
		}

		// If server already has it, skip the data bytes.
		if s.store.Has(hash) {
			if _, err := io.CopyN(io.Discard, r.Body, int64(size)); err != nil {
				http.Error(w, "skip data: "+err.Error(), http.StatusBadRequest)
				return
			}
			resp.Skipped++
			continue
		}

		// Read data.
		data := make([]byte, size)
		if _, err := io.ReadFull(r.Body, data); err != nil {
			http.Error(w, "read data: "+err.Error(), http.StatusBadRequest)
			return
		}

		// Verify hash.
		hasher.Reset()
		hasher.Write(data)
		computed := hex.EncodeToString(hasher.Sum(nil))
		if computed != hash {
			resp.Failed++
			continue
		}

		// Store.
		accepted, err := s.store.Put(hash, data)
		if err != nil {
			resp.Failed++
			continue
		}
		if accepted {
			resp.Accepted++
		} else {
			resp.Skipped++
		}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

// handleReset wipes all chunks and manifests from the store.
func (s *Server) handleReset(w http.ResponseWriter, r *http.Request) {
	log.Println("resetting store...")
	if err := s.store.Reset(); err != nil {
		http.Error(w, "reset failed: "+err.Error(), http.StatusInternalServerError)
		return
	}
	log.Println("store reset complete")
	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(`{"status":"reset"}`))
}
