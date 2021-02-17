package data

import (
	log "github.com/sirupsen/logrus"
	"github.com/spf13/afero"
	"math"
	"path"
)

//go:generate mockery --name=Splitter --case underscore
type Splitter interface {
	Split(file afero.File) (*[]Chunk, error)
	NumChunks(fileSize int64) uint64
}

type SplitJob struct {
	ChunkWriter ChunkWriter
	HashFunc    ChunkHashFunc
	ChunkDir    string
	ChunkSize   uint64
}

func NewSplitJob(chunkWriter ChunkWriter, hashFunc ChunkHashFunc, chunkDir string, chunkMBs uint64) *SplitJob {
	if chunkMBs < 1 {
		chunkMBs = 1
	}
	chunkSize := chunkMBs * (1 << 20)
	return &SplitJob{ChunkWriter: chunkWriter, HashFunc: hashFunc, ChunkDir: chunkDir, ChunkSize: chunkSize}
}

func (s *SplitJob) NumChunks(fileSize int64) uint64 {
	return uint64(math.Ceil(float64(fileSize) / float64(s.ChunkSize))) // total number of chunks
}

func (s *SplitJob) Split(file afero.File) (*[]Chunk, error) {
	// Prepare the chunk array
	stats, err := file.Stat()
	if err != nil {
		return nil, err
	}
	fileSize := stats.Size()
	numChunks := s.NumChunks(fileSize)
	log.Infof("Splitting to %d pieces.", numChunks)
	chunks := make([]Chunk, numChunks)

	for i := uint64(0); i < numChunks; i++ {
		numChunkBytes := int(math.Min(float64(s.ChunkSize), float64(fileSize-int64(i*s.ChunkSize))))
		chunkBytes := make([]byte, numChunkBytes)

		_, err := file.Read(chunkBytes)
		if err != nil {
			log.Errorf("Failed to read from the chunk buffer: %v", err)
			return nil, err
		}

		// write to disk
		chunkHash := s.HashFunc(chunkBytes)
		chunkPath := path.Join(s.ChunkDir, chunkHash)
		chunk := NewChunk(chunkPath, chunkHash, i)

		err = s.ChunkWriter(chunk, &chunkBytes)
		if err != nil {
			log.Errorf("Failed to write the chunk data to: '%s'", chunkPath)
			return nil, err
		}

		chunks[i] = *chunk

		log.Debug("Split to : ", chunkPath)
	}

	return &chunks, nil
}
