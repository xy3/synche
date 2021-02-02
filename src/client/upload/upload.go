package upload

import (
	"encoding/hex"
	"github.com/go-openapi/runtime"
	"github.com/kalafut/imohash"
	"gitlab.computing.dcu.ie/collint9/2021-ca400-collint9-coynemt2/src/client/apiclient/files"
	"os"
)

func NewChunkUploadParams(chunkPath, chunkName string) (*files.UploadFileParams, error) {
	chunkHash, err := imohash.SumFile(chunkPath)
	if err != nil {
		return nil, err
	}

	file, err := os.Open(chunkPath)
	if err != nil {
		return nil, err
	}

	hash := hex.EncodeToString(chunkHash[:])
	uploadRequestId := "example1"
	readCloser := runtime.NamedReader(chunkName, file)
	params := files.NewUploadFileParams().
		WithChunkNumber(0).
		WithChunkHash(hash).
		WithUploadRequestID(uploadRequestId).
		WithChunkData(readCloser)

	return params, nil
}
