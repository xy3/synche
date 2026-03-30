// Package server implements the HTTP chunk storage server.
package server

import (
	"crypto/subtle"
	"crypto/tls"
	"encoding/binary"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"

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

// Config holds server configuration.
type Config struct {
	Addr      string // listen address (e.g. ":8420")
	StorePath string // path to chunk store directory
	APIKey    string // required API key (empty = no auth)
	TLSCert   string // path to TLS certificate file (empty = plain HTTP)
	TLSKey    string // path to TLS private key file
}

// Server is the chunk storage HTTP server.
type Server struct {
	store  *store.Store
	mux    *http.ServeMux
	config Config
}

// New creates a new server with the given configuration.
func New(cfg Config) (*Server, error) {
	s, err := store.New(cfg.StorePath)
	if err != nil {
		return nil, fmt.Errorf("open store: %w", err)
	}

	srv := &Server{
		store:  s,
		mux:    http.NewServeMux(),
		config: cfg,
	}
	srv.routes()
	return srv, nil
}

// ListenAndServe starts the HTTP or HTTPS server.
func (s *Server) ListenAndServe() error {
	log.Printf("synche server listening on %s", s.config.Addr)
	stats := s.store.Stats()
	log.Printf("store: %d existing chunks", stats.TotalChunks)
	if s.config.APIKey != "" {
		log.Printf("api key: enabled")
	} else {
		log.Printf("api key: disabled (WARNING: no authentication)")
	}

	if s.config.TLSCert != "" && s.config.TLSKey != "" {
		log.Printf("tls: enabled")
		log.Printf("webdav available at https://%s/webdav/", s.config.Addr)
		tlsConfig := &tls.Config{
			MinVersion: tls.VersionTLS12,
		}
		server := &http.Server{
			Addr:      s.config.Addr,
			Handler:   s.mux,
			TLSConfig: tlsConfig,
		}
		return server.ListenAndServeTLS(s.config.TLSCert, s.config.TLSKey)
	}

	log.Printf("tls: disabled (WARNING: traffic is unencrypted)")
	log.Printf("webdav available at http://%s/webdav/", s.config.Addr)
	return http.ListenAndServe(s.config.Addr, s.mux)
}

func (s *Server) routes() {
	// Health check — always unauthenticated.
	s.mux.HandleFunc("GET /api/health", s.handleHealth)

	// All other API endpoints require authentication.
	s.mux.HandleFunc("POST /api/probe", s.requireAuth(s.handleProbe))
	s.mux.HandleFunc("PUT /api/chunk/{hash}", s.requireAuth(s.handleUploadChunk))
	s.mux.HandleFunc("GET /api/chunk/{hash}", s.requireAuth(s.handleGetChunk))
	s.mux.HandleFunc("POST /api/upload-batch", s.requireAuth(s.handleBatchUpload))
	s.mux.HandleFunc("POST /api/manifest", s.requireAuth(s.handleUploadManifest))
	s.mux.HandleFunc("DELETE /api/manifest/{id}", s.requireAuth(s.handleDeleteManifest))
	s.mux.HandleFunc("GET /api/stats", s.requireAuth(s.handleStats))

	// WebDAV: read-only virtual filesystem of uploaded files.
	// Authenticated so files aren't publicly browsable.
	dav := webdav.New(s.store, "/webdav")
	s.mux.Handle("/webdav/", s.requireAuthHandler(dav))
	s.mux.Handle("/webdav", http.RedirectHandler("/webdav/", http.StatusMovedPermanently))
}

// requireAuth wraps a handler function with API key authentication.
// If no API key is configured, all requests are allowed.
func (s *Server) requireAuth(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if !s.checkAuth(w, r) {
			return
		}
		next(w, r)
	}
}

// requireAuthHandler wraps an http.Handler with API key authentication.
func (s *Server) requireAuthHandler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !s.checkAuth(w, r) {
			return
		}
		next.ServeHTTP(w, r)
	})
}

// checkAuth validates the Authorization header. Returns true if the request
// is authorized, false if it was rejected (and a 401 was written).
func (s *Server) checkAuth(w http.ResponseWriter, r *http.Request) bool {
	if s.config.APIKey == "" {
		return true // no auth configured
	}

	auth := r.Header.Get("Authorization")
	if auth == "" {
		w.Header().Set("WWW-Authenticate", `Bearer realm="synche"`)
		http.Error(w, "missing authorization header", http.StatusUnauthorized)
		return false
	}

	// Expect "Bearer <key>"
	const prefix = "Bearer "
	if !strings.HasPrefix(auth, prefix) {
		http.Error(w, "invalid authorization format (expected: Bearer <key>)", http.StatusUnauthorized)
		return false
	}
	token := auth[len(prefix):]

	if subtle.ConstantTimeCompare([]byte(token), []byte(s.config.APIKey)) != 1 {
		http.Error(w, "invalid api key", http.StatusUnauthorized)
		return false
	}

	return true
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

// handleDeleteManifest deletes a manifest by ID and removes any chunks that
// are no longer referenced by any remaining manifest.
func (s *Server) handleDeleteManifest(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if !store.ValidManifestID(id) {
		http.Error(w, "invalid manifest ID", http.StatusBadRequest)
		return
	}

	deleted, err := s.store.DeleteManifestAndOrphanedChunks(id)
	if err != nil {
		http.Error(w, "delete manifest: "+err.Error(), http.StatusNotFound)
		return
	}

	log.Printf("deleted manifest %s (%d orphaned chunks removed)", id, deleted)
	w.WriteHeader(http.StatusNoContent)
}
