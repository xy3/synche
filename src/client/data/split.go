package data

import (
	log "github.com/sirupsen/logrus"
	"io"
	"math"
)

type Splitter interface {
	NumChunks() int64
	NextChunk() ([]byte, error)
	Split(handleChunk func(*Chunk) error) error
	File() *SplitFile
}

type SplitFile struct {
	FileSize     int64
	ChunkSize    int64
	CurrentIndex int64
	Path         string
	Name         string
	Hash         string
	Reader       io.Reader
}

func (sf *SplitFile) File() *SplitFile {
	return sf
}

func NewSplitFile(fileSize, chunkMBs int64, path, name, hash string, reader io.Reader) *SplitFile {
	if chunkMBs < 1 {
		chunkMBs = 1
	}
	chunkSize := chunkMBs * KB
	return &SplitFile{
		FileSize:     fileSize,
		ChunkSize:    chunkSize,
		CurrentIndex: 0,
		Path:         path,
		Name:         name,
		Hash:         hash,
		Reader:       reader,
	}
}

func (sf SplitFile) NumChunks() int64 {
	return int64(math.Ceil(float64(sf.FileSize) / float64(sf.ChunkSize)))
}

func (sf *SplitFile) NextChunk() ([]byte, error) {
	if sf.CurrentIndex >= sf.NumChunks() {
		return nil, nil
	}
	// Use ChunkSize only if it is smaller than the rest of the file
	numChunkBytes := sf.FileSize - (sf.CurrentIndex * sf.ChunkSize)
	if sf.ChunkSize < numChunkBytes {
		numChunkBytes = sf.ChunkSize
	}

	chunkBytes := make([]byte, numChunkBytes)
	_, err := sf.Reader.Read(chunkBytes)
	if err != nil {
		log.Error("Couldn't read from reader")
		return chunkBytes, err
	}

	sf.CurrentIndex++
	return chunkBytes, nil
}

func (sf *SplitFile) Split(handleChunk func(*Chunk) error) error {
	for sf.CurrentIndex < sf.NumChunks() {
		chunkBytes, err := sf.NextChunk()
		if err != nil {
			return err
		}
		chunk := NewChunk(sf.CurrentIndex, &chunkBytes)
		if err = handleChunk(chunk); err != nil {
			return err
		}
	}
	return nil
}
