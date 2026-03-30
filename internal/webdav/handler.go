// Package webdav provides a read-only WebDAV handler that exposes uploaded
// manifests as virtual files. When a client downloads a file, the server
// reassembles it on the fly by streaming chunks from the store in order.
package webdav

import (
	"encoding/xml"
	"fmt"
	"io"
	"log"
	"net/http"
	"path"
	"strconv"
	"strings"
	"time"

	"github.com/theo/synche2/internal/store"
)

// Handler serves WebDAV requests backed by a chunk store.
type Handler struct {
	store  *store.Store
	prefix string // URL path prefix, e.g. "/webdav"
}

// New creates a WebDAV handler. prefix is the URL path prefix (e.g. "/webdav").
func New(s *store.Store, prefix string) *Handler {
	return &Handler{store: s, prefix: strings.TrimRight(prefix, "/")}
}

// ServeHTTP dispatches WebDAV methods.
func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// Strip the prefix to get the virtual path.
	vpath := strings.TrimPrefix(r.URL.Path, h.prefix)
	if vpath == "" {
		vpath = "/"
	}
	if !strings.HasPrefix(vpath, "/") {
		vpath = "/" + vpath
	}

	switch r.Method {
	case "OPTIONS":
		h.handleOptions(w, r)
	case "PROPFIND":
		h.handlePropfind(w, r, vpath)
	case "GET", "HEAD":
		h.handleGet(w, r, vpath)
	default:
		// Read-only: reject writes.
		w.Header().Set("Allow", "OPTIONS, PROPFIND, GET, HEAD")
		http.Error(w, "method not allowed (read-only WebDAV)", http.StatusMethodNotAllowed)
	}
}

// handleOptions returns DAV capabilities.
func (h *Handler) handleOptions(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Allow", "OPTIONS, PROPFIND, GET, HEAD")
	w.Header().Set("DAV", "1")
	w.WriteHeader(http.StatusOK)
}

// handlePropfind returns directory listings or file properties.
func (h *Handler) handlePropfind(w http.ResponseWriter, r *http.Request, vpath string) {
	// Parse Depth header (default 1 for collections, 0 for resources).
	depth := r.Header.Get("Depth")
	if depth == "" {
		depth = "1"
	}

	manifests, err := h.store.ListManifests()
	if err != nil {
		log.Printf("webdav: list manifests: %v", err)
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}

	// Build file map: manifest ID -> info. Use ID as filename to avoid collisions.
	// Filename format: <source_name>.<id> but we keep it simple: just use ID.
	// Actually, let's use: <source_basename>--<id> so it's human-readable.

	vpath = path.Clean(vpath)

	if vpath == "/" || vpath == "." {
		// Collection listing.
		h.propfindCollection(w, manifests, depth)
		return
	}

	// Single resource lookup.
	filename := strings.TrimPrefix(vpath, "/")
	info := h.findManifest(manifests, filename)
	if info == nil {
		http.Error(w, "not found", http.StatusNotFound)
		return
	}

	h.propfindFile(w, info)
}

// handleGet streams the reassembled file from chunks.
func (h *Handler) handleGet(w http.ResponseWriter, r *http.Request, vpath string) {
	manifests, err := h.store.ListManifests()
	if err != nil {
		log.Printf("webdav: list manifests: %v", err)
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}

	filename := strings.TrimPrefix(path.Clean(vpath), "/")
	if filename == "" || filename == "." {
		// GET on root — return a simple HTML directory listing.
		h.handleDirectoryListing(w, manifests)
		return
	}

	info := h.findManifest(manifests, filename)
	if info == nil {
		http.Error(w, "not found", http.StatusNotFound)
		return
	}

	manifest, err := h.store.LoadManifest(info.ID)
	if err != nil {
		log.Printf("webdav: load manifest %s: %v", info.ID, err)
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/octet-stream")
	w.Header().Set("Content-Length", strconv.FormatUint(manifest.TotalBytes, 10))
	w.Header().Set("Content-Disposition", fmt.Sprintf(`attachment; filename="%s"`, info.Name))
	w.Header().Set("Last-Modified", info.ModTime.UTC().Format(http.TimeFormat))

	if r.Method == "HEAD" {
		return
	}

	// Stream chunks in order.
	for i, chunk := range manifest.Chunks {
		data, err := h.store.ReadChunk(chunk.Hash)
		if err != nil {
			log.Printf("webdav: read chunk %d (%s) of %s: %v", i, chunk.Hash[:8], info.ID, err)
			// Can't send error after headers are written; just stop.
			return
		}
		if _, err := w.Write(data); err != nil {
			return // client disconnected
		}
	}
}

// handleDirectoryListing returns a simple HTML page listing files.
func (h *Handler) handleDirectoryListing(w http.ResponseWriter, manifests []store.ManifestInfo) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	fmt.Fprintf(w, "<!DOCTYPE html><html><head><title>synche files</title></head><body>\n")
	fmt.Fprintf(w, "<h1>synche files</h1><table><tr><th>File</th><th>Size</th><th>Chunks</th><th>Date</th></tr>\n")
	for _, m := range manifests {
		fname := manifestFilename(m)
		fmt.Fprintf(w, "<tr><td><a href=\"%s/%s\">%s</a></td><td>%s</td><td>%d</td><td>%s</td></tr>\n",
			h.prefix, fname, fname, humanSize(m.TotalBytes), m.NumChunks, m.ModTime.Format(time.RFC3339))
	}
	fmt.Fprintf(w, "</table></body></html>\n")
}

// findManifest looks up a manifest by its virtual filename.
func (h *Handler) findManifest(manifests []store.ManifestInfo, filename string) *store.ManifestInfo {
	for i := range manifests {
		if manifestFilename(manifests[i]) == filename {
			return &manifests[i]
		}
	}
	return nil
}

// manifestFilename returns the virtual filename for a manifest.
// Format: <source_name>--<id>  (e.g., "myfile.img--myfile.img_1774869887324831679")
// But since ID already contains source name, just use ID directly.
func manifestFilename(m store.ManifestInfo) string {
	return m.ID
}

func humanSize(b uint64) string {
	const (
		KB = 1024
		MB = 1024 * KB
		GB = 1024 * MB
	)
	switch {
	case b >= GB:
		return fmt.Sprintf("%.1f GB", float64(b)/float64(GB))
	case b >= MB:
		return fmt.Sprintf("%.1f MB", float64(b)/float64(MB))
	case b >= KB:
		return fmt.Sprintf("%.1f KB", float64(b)/float64(KB))
	default:
		return fmt.Sprintf("%d B", b)
	}
}

// ---- WebDAV XML responses ----

// multistatus is the top-level DAV:multistatus response.
type multistatus struct {
	XMLName   xml.Name   `xml:"DAV: multistatus"`
	Responses []response `xml:"response"`
}

type response struct {
	Href     string   `xml:"href"`
	Propstat propstat `xml:"propstat"`
}

type propstat struct {
	Prop   prop   `xml:"prop"`
	Status string `xml:"status"`
}

type prop struct {
	DisplayName  string        `xml:"displayname,omitempty"`
	ResourceType *resourceType `xml:"resourcetype"`
	ContentLen   string        `xml:"getcontentlength,omitempty"`
	ContentType  string        `xml:"getcontenttype,omitempty"`
	LastModified string        `xml:"getlastmodified,omitempty"`
	ETag         string        `xml:"getetag,omitempty"`
}

type resourceType struct {
	Collection *struct{} `xml:"collection,omitempty"`
}

func (h *Handler) propfindCollection(w http.ResponseWriter, manifests []store.ManifestInfo, depth string) {
	ms := multistatus{}

	// The collection itself.
	ms.Responses = append(ms.Responses, response{
		Href: h.prefix + "/",
		Propstat: propstat{
			Prop: prop{
				DisplayName:  "synche",
				ResourceType: &resourceType{Collection: &struct{}{}},
			},
			Status: "HTTP/1.1 200 OK",
		},
	})

	if depth != "0" {
		// Children.
		for _, m := range manifests {
			fname := manifestFilename(m)
			ms.Responses = append(ms.Responses, response{
				Href: h.prefix + "/" + fname,
				Propstat: propstat{
					Prop: prop{
						DisplayName:  fname,
						ResourceType: &resourceType{},
						ContentLen:   strconv.FormatUint(m.TotalBytes, 10),
						ContentType:  "application/octet-stream",
						LastModified: m.ModTime.UTC().Format(http.TimeFormat),
					},
					Status: "HTTP/1.1 200 OK",
				},
			})
		}
	}

	w.Header().Set("Content-Type", "application/xml; charset=utf-8")
	w.WriteHeader(http.StatusMultiStatus) // 207
	io.WriteString(w, xml.Header)
	enc := xml.NewEncoder(w)
	enc.Indent("", "  ")
	enc.Encode(ms)
}

func (h *Handler) propfindFile(w http.ResponseWriter, info *store.ManifestInfo) {
	fname := manifestFilename(*info)
	ms := multistatus{
		Responses: []response{{
			Href: h.prefix + "/" + fname,
			Propstat: propstat{
				Prop: prop{
					DisplayName:  fname,
					ResourceType: &resourceType{},
					ContentLen:   strconv.FormatUint(info.TotalBytes, 10),
					ContentType:  "application/octet-stream",
					LastModified: info.ModTime.UTC().Format(http.TimeFormat),
				},
				Status: "HTTP/1.1 200 OK",
			},
		}},
	}

	w.Header().Set("Content-Type", "application/xml; charset=utf-8")
	w.WriteHeader(http.StatusMultiStatus)
	io.WriteString(w, xml.Header)
	enc := xml.NewEncoder(w)
	enc.Indent("", "  ")
	enc.Encode(ms)
}
