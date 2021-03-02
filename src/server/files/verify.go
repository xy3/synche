package files

import (
	"encoding/hex"
	"github.com/kalafut/imohash"
	"os"
)

func VerifyChunk(chunkFilePath string) (bool, error) {
	file, err := os.Open(chunkFilePath)
	if err != nil {
		return false, err
	}
	filename := file.Name()

	hash, err := imohash.SumFile(chunkFilePath)
	if err != nil {
		return false, err
	}

	hashHex := hex.EncodeToString(hash[:])

	return filename == hashHex, nil
}


func VerifyComposite(filename string) (bool, error) {
	// do something like "select * from data where file == filename
	return false, nil
}