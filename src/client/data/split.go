package data

import (
	log "github.com/sirupsen/logrus"
	"math"
	"os"
	"path"
)

type SplitJob struct {
	chunkWriter ChunkWriter
	hashFunc 	ChunkHashFunc
	chunkDir    string
	chunkSize   uint64
}


func NewSplitJob(chunkWriter ChunkWriter, hashFunc ChunkHashFunc, chunkDir string, chunkMBs uint64) *SplitJob {
	if chunkMBs < 1 {
		chunkMBs = 1
	}
	chunkSize := chunkMBs * (1 << 20)
	return &SplitJob{chunkWriter: chunkWriter, hashFunc: hashFunc, chunkDir: chunkDir, chunkSize: chunkSize}
}

func (s *SplitJob) Split(file *os.File) (*[]Chunk, error) {
	// Prepare the chunk array
	stats, err := file.Stat()
	if err != nil {
		return nil, err
	}
	fileSize := stats.Size()
	numChunks := uint64(math.Ceil(float64(fileSize) / float64(s.chunkSize))) // total number of chunks
	log.Printf("Splitting to %d pieces.\n", numChunks)
	chunks := make([]Chunk, numChunks)

	for i := uint64(0); i < numChunks; i++ {
		numChunkBytes := int(math.Min(float64(s.chunkSize), float64(fileSize-int64(i*s.chunkSize))))
		chunkBytes := make([]byte, numChunkBytes)

		_, err := file.Read(chunkBytes)
		if err != nil {
			log.Errorf("Failed to read from the chunk buffer: %v", err)
			return nil, err
		}

		// write to disk
		chunkHash := s.hashFunc(chunkBytes)
		chunkPath := path.Join(s.chunkDir, chunkHash)
		chunk := NewChunk(chunkPath, chunkHash, i)

		err = s.chunkWriter(chunk, &chunkBytes)
		if err != nil {
			log.Errorf("Failed to write the chunk data to: '%s'", chunkPath)
			return nil, err
		}

		chunks[i] = *chunk

		log.Debug("Split to : ", chunkPath)
	}

	return &chunks, nil
}
