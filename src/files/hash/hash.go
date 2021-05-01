package hash

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"github.com/kalafut/imohash"
	"gitlab.computing.dcu.ie/collint9/2021-ca400-collint9-coynemt2/src/files"
	"hash/crc32"
	"strings"
)

type ChunkHashFunc func(chunkBytes []byte) string
type FileHashFunc func(filePath string) (string, error)

var Chunk = MD5Hash
var File = ImoHash

func MD5Hash(bytes []byte) string {
	return fmt.Sprintf("%x", md5.Sum(bytes))
}

func MD5HashString(input string) string {
	return MD5Hash([]byte(input))
}

func PathHash(path string) string {
	return MD5HashString(strings.TrimRight(strings.TrimSpace(path), "/"))
}

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
