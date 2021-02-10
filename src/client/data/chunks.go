package data

import (
	"io/ioutil"
	log "github.com/sirupsen/logrus"
	"os"
)

type ChunkWriter func(chunkPath string, chunkBytes *[]byte) error

func DefaultChunkWriter(chunkPath string, chunkBytes *[]byte) error {
	_, err := os.Create(chunkPath)
	if err != nil {
		log.Printf("Failed to create a new chunk file: %v", err)
		return err
	}

	// write/save buffer to disk
	err = ioutil.WriteFile(chunkPath, *chunkBytes, os.ModeAppend)
	if err != nil {
		log.Printf("Failed to write the chunk data to a new file: %v", err)
		return err
	}

	return nil
}
