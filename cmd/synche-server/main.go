package main

import (
	"flag"
	"log"
	"os"

	"github.com/theo/synche2/internal/server"
)

func main() {
	log.SetFlags(log.Ltime | log.Lmicroseconds)

	addr := flag.String("addr", ":8420", "listen address")
	storePath := flag.String("store", "./synche-store", "chunk store directory")
	apiKey := flag.String("api-key", "", "API key for authentication (or set SYNCHE_API_KEY env var)")
	tlsCert := flag.String("tls-cert", "", "path to TLS certificate file")
	tlsKey := flag.String("tls-key", "", "path to TLS private key file")
	flag.Parse()

	// API key from flag takes precedence, then env var.
	key := *apiKey
	if key == "" {
		key = os.Getenv("SYNCHE_API_KEY")
	}

	cfg := server.Config{
		Addr:      *addr,
		StorePath: *storePath,
		APIKey:    key,
		TLSCert:   *tlsCert,
		TLSKey:    *tlsKey,
	}

	log.Printf("synche server starting")
	log.Printf("  listen: %s", cfg.Addr)
	log.Printf("  store:  %s", cfg.StorePath)

	srv, err := server.New(cfg)
	if err != nil {
		log.Fatalf("failed to create server: %v", err)
	}

	if err := srv.ListenAndServe(); err != nil {
		log.Fatalf("server error: %v", err)
	}
}
