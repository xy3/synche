package main

import (
	"flag"
	"log"

	"github.com/theo/synche2/internal/server"
)

func main() {
	log.SetFlags(log.Ltime | log.Lmicroseconds)

	addr := flag.String("addr", ":8420", "listen address")
	storePath := flag.String("store", "./synche-store", "chunk store directory")
	flag.Parse()

	log.Printf("synche server starting")
	log.Printf("  listen: %s", *addr)
	log.Printf("  store:  %s", *storePath)

	srv, err := server.New(*addr, *storePath)
	if err != nil {
		log.Fatalf("failed to create server: %v", err)
	}

	if err := srv.ListenAndServe(); err != nil {
		log.Fatalf("server error: %v", err)
	}
}
