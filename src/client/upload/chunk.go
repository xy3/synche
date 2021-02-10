package upload

import (
	"encoding/hex"
	"github.com/go-openapi/runtime"
	"github.com/kalafut/imohash"
	log "github.com/sirupsen/logrus"
	"gitlab.computing.dcu.ie/collint9/2021-ca400-collint9-coynemt2/src/client/apiclient/files"
	"os"
	"sync"
)

type ChunkUploadApiRequest func(params *files.UploadChunkParams) (*files.UploadChunkCreated, error)

type ChunkUploader struct {
	Send ChunkUploadApiRequest
}

func NewChunkUploader(sender ChunkUploadApiRequest) *ChunkUploader {
	return &ChunkUploader{Send: sender}
}

func (cu *ChunkUploader) Upload(wg *sync.WaitGroup, params *files.UploadChunkParams, uploadErrors chan error) {
	defer wg.Done()

	resp, err := cu.Send(params)
	if err != nil {
		uploadErrors <- err
		return
	}

	log.Debugf("%#v\n", resp.Payload)

	// TODO: Do something here with the response payload to check if the chunk was uploaded correctly
}

func (cu *ChunkUploader) NewParams(
	chunkPath, chunkName, uploadRequestId string,
	chunkNum int64,
) (*files.UploadChunkParams, error) {
	chunkHash, err := imohash.SumFile(chunkPath)
	if err != nil {
		return nil, err
	}

	file, err := os.Open(chunkPath)
	if err != nil {
		return nil, err
	}

	hash := hex.EncodeToString(chunkHash[:])
	readCloser := runtime.NamedReader(chunkName, file)
	params := files.NewUploadChunkParams().
		WithChunkNumber(chunkNum).
		WithChunkHash(hash).
		WithUploadRequestID(uploadRequestId).
		WithChunkData(readCloser)

	return params, nil
}
