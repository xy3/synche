package hash

import (
	"crypto/md5"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"github.com/kalafut/imohash"
	"github.com/xy3/synche/src/files"
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

func SHA256Hash(bytes []byte) string {
	return fmt.Sprintf("%x", sha256.Sum256(bytes))
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

// Random is used to generate random bytes
func Random(length uint) []byte {
	var r = make([]byte, length)
	_, _ = rand.Reader.Read(r)
	return r
}

// RandomMD5Hash is used to generate random and hashed by md5
func RandomMD5Hash(bytes uint) string {
	hash := md5.New()
	_, _ = hash.Write(Random(bytes))
	return hex.EncodeToString(hash.Sum(nil))
}
