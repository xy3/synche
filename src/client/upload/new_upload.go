package upload

import (
	"encoding/hex"
	"github.com/kalafut/imohash"
	"gitlab.computing.dcu.ie/collint9/2021-ca400-collint9-coynemt2/src/client/apiclient/files"
	"gitlab.computing.dcu.ie/collint9/2021-ca400-collint9-coynemt2/src/client/models"
	"os"
)

func NewFileUploadParams(filePath, fileName string, numChunks int64) (*files.NewUploadParams, error) {
	fileHash, err := imohash.SumFile(filePath)
	if err != nil {
		return nil, err
	}

	fileHashString := hex.EncodeToString(fileHash[:])
	fileInfo, err := os.Stat(filePath)
	if err != nil {
		return nil, err
	}
	fileSize := fileInfo.Size()

	params := files.NewNewUploadParams().WithFileInfo(&models.FileInfo{
		Chunks: &numChunks,
		Hash:   &fileHashString,
		Name:   &fileName,
		Size:   &fileSize,
	})

	return params, nil
}
