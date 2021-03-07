package data

import (
	"encoding/hex"
	"fmt"
	"github.com/kalafut/imohash"
	"gitlab.computing.dcu.ie/collint9/2021-ca400-collint9-coynemt2/src/client/files"
	"hash/crc32"
)

type ChunkHashFunc func(chunkBytes []byte) string
type FileHashFunc func(filePath string) (string, error)

var DefaultChunkHashFunc = CRC32Hash
var DefaultFileHashFunc = ImoHash

func CRC32Hash(bytes []byte) string {
	checksum := crc32.ChecksumIEEE(bytes)
	return fmt.Sprintf("%08x", checksum)
}

func ImoHash(filePath string) (hash string, err error) {
	fileData, err := files.Afs.ReadFile(filePath)
	if err != nil {
		return hash, err
	}
	fileHash := imohash.Sum(fileData)
	return hex.EncodeToString(fileHash[:]), nil
}