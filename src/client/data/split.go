package data

import (
	"encoding/hex"
	"github.com/kalafut/imohash"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"math"
	"os"
)

type Chunk struct {
	Path string
	Hash string
	Num  uint64
}

func NewChunk(path string, hash string, number uint64) *Chunk {
	return &Chunk{Path: path, Hash: hash, Num: number}
}

type Splitter struct {
	chunkWriter ChunkWriter
}

func NewSplitter(chunkWriter ChunkWriter) *Splitter {
	return &Splitter{chunkWriter}
}

func (s *Splitter) Split(filePath, chunkDir string) ([]Chunk, error) {
	file, err := os.Open(filePath)

	if err != nil {
		log.Errorf("Could not open: '%s'", filePath)
		return nil, err
	}

	defer file.Close()

	fileInfo, _ := file.Stat()
	fileSize := fileInfo.Size()
	fileChunk := viper.GetInt("ChunkSize") * (1 << 20)

	// calculate total number of parts the file will be chunked into
	totalPartsNum := uint64(math.Ceil(float64(fileSize) / float64(fileChunk)))
	chunks := make([]Chunk, totalPartsNum)

	log.Printf("Splitting to %d pieces.\n", totalPartsNum)
	for i := uint64(0); i < totalPartsNum; i++ {

		chunkSize := int(math.Min(float64(fileChunk), float64(fileSize-int64(i*uint64(fileChunk)))))
		chunkBytes := make([]byte, chunkSize)

		_, err := file.Read(chunkBytes)
		if err != nil {
			log.Errorf("Failed to read from the chunk buffer: %v", err)
			return nil, err
		}

		// write to disk
		hash := imohash.Sum(chunkBytes)
		chunkHash := hex.EncodeToString(hash[:])
		chunkPath := chunkDir + "/" + chunkHash
		err = s.chunkWriter(chunkPath, &chunkBytes)
		if err != nil {
			log.Errorf("Failed to write the chunk data to: '%s'", chunkPath)
			return nil, err
		}

		chunks[i] = *NewChunk(chunkPath, chunkHash, i)

		log.Debug("Split to : ", chunkPath)
	}

	return chunks, err
}
