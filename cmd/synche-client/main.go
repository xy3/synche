package main

import (
	"flag"
	"log"
	"os"
	"runtime"

	"github.com/theo/synche2/internal/client"
)

func main() {
	log.SetFlags(log.Ltime | log.Lmicroseconds)

	serverURL := flag.String("server", "http://localhost:8420", "server URL")
	source := flag.String("source", "", "block device or file path to upload (e.g. /dev/sda, /path/to/file)")
	concurrency := flag.Int("concurrency", runtime.NumCPU(), "number of parallel upload/hash workers")
	probeBatch := flag.Int("probe-batch", 512, "number of hashes per probe request")
	cacheDir := flag.String("cache-dir", "", "manifest cache directory (default: ~/.cache/synche)")
	noCache := flag.Bool("no-cache", false, "disable local manifest caching")

	// Keep --device as an alias for backwards compat.
	device := flag.String("device", "", "alias for --source")
	flag.Parse()

	src := *source
	if src == "" {
		src = *device
	}
	if src == "" {
		log.Fatal("--source is required (e.g. --source /dev/sda or --source /path/to/file)")
	}

	// Verify source is readable.
	f, err := os.Open(src)
	if err != nil {
		log.Fatalf("cannot open %s: %v (try running with sudo for block devices)", src, err)
	}
	f.Close()

	cfg := client.Config{
		ServerURL:   *serverURL,
		DevicePath:  src,
		Concurrency: *concurrency,
		ProbeBatch:  *probeBatch,
		CacheDir:    *cacheDir,
		NoCache:     *noCache,
	}

	log.Printf("synche client starting")
	log.Printf("  source:      %s", cfg.DevicePath)
	log.Printf("  server:      %s", cfg.ServerURL)
	log.Printf("  concurrency: %d", cfg.Concurrency)
	log.Printf("  probe batch: %d", cfg.ProbeBatch)
	if cfg.NoCache {
		log.Printf("  cache:       disabled")
	}

	uploader := client.NewUploader(cfg)
	if err := uploader.Run(); err != nil {
		log.Fatalf("upload failed: %v", err)
	}
}
