// Package server implements the HTTP chunk storage server.
package server

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/theo/synche2/internal/protocol"
	"github.com/theo/synche2/internal/store"
	"github.com/zeebo/blake3"
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
	log.Printf("store: %s (%d existing chunks)", stats.StorePath, stats.TotalChunks)
	return http.ListenAndServe(s.addr, s.mux)
}

func (s *Server) routes() {
	s.mux.HandleFunc("POST /api/probe", s.handleProbe)
	s.mux.HandleFunc("PUT /api/chunk/{hash}", s.handleUploadChunk)
	s.mux.HandleFunc("GET /api/chunk/{hash}", s.handleGetChunk)
	s.mux.HandleFunc("POST /api/manifest", s.handleUploadManifest)
	s.mux.HandleFunc("POST /api/reset", s.handleReset)
	s.mux.HandleFunc("GET /api/stats", s.handleStats)
	s.mux.HandleFunc("GET /api/health", s.handleHealth)
}

// handleProbe accepts a list of hashes and returns which ones the server needs.
func (s *Server) handleProbe(w http.ResponseWriter, r *http.Request) {
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
	if len(hash) != protocol.HashSize*2 {
		http.Error(w, "invalid hash length", http.StatusBadRequest)
		return
	}

	// Fast reject: already have it.
	if s.store.Has(hash) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(protocol.UploadResponse{
			Hash:     hash,
			Accepted: false,
		})
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

	accepted, err := s.store.Put(hash, data)
	if err != nil {
		http.Error(w, "store chunk: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(protocol.UploadResponse{
		Hash:     hash,
		Accepted: accepted,
	})
}

// handleGetChunk retrieves a chunk by hash.
func (s *Server) handleGetChunk(w http.ResponseWriter, r *http.Request) {
	hash := r.PathValue("hash")
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
