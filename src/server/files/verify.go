package files

import (
	"encoding/hex"
	"github.com/kalafut/imohash"
	"gitlab.computing.dcu.ie/collint9/2021-ca400-collint9-coynemt2/src/files"
)

// Todo: this function should use the correct hashing algorithm to verify the chunks (instead of imohash)
func VerifyChunk(chunkFilePath string) (bool, error) {
	file, err := files.Afs.Open(chunkFilePath)
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
