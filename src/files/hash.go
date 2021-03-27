package files

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"github.com/kalafut/imohash"
	"hash/crc32"
)

type ChunkHashFunc func(chunkBytes []byte) string
type FileHashFunc func(filePath string) (string, error)

var HashChunk = MD5Hash
var HashFile = ImoHash

func MD5Hash(bytes []byte) string {
	return fmt.Sprintf("%x", md5.Sum(bytes))
}

func CRC32Hash(bytes []byte) string {
	checksum := crc32.ChecksumIEEE(bytes)
	return fmt.Sprintf("%08x", checksum)
}

func ImoHash(filePath string) (hash string, err error) {
	fileData, err := Afs.ReadFile(filePath)
	if err != nil {
		return hash, err
	}
	fileHash := imohash.Sum(fileData)
	return hex.EncodeToString(fileHash[:]), nil
}
