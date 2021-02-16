package upload

import (
	"github.com/go-openapi/runtime"
	log "github.com/sirupsen/logrus"
	"gitlab.computing.dcu.ie/collint9/2021-ca400-collint9-coynemt2/src/client/apiclient"
	"gitlab.computing.dcu.ie/collint9/2021-ca400-collint9-coynemt2/src/client/apiclient/files"
	"gitlab.computing.dcu.ie/collint9/2021-ca400-collint9-coynemt2/src/client/data"
	"os"
	"sync"
)

type ChunkUploader interface {
	Upload(wg *sync.WaitGroup, params *files.UploadChunkParams, uploadErrors chan error)
	NewParams(chunk data.Chunk, requestID string) (*files.UploadChunkParams, error)
}

type ChunkUpload struct {
	hashFunc data.ChunkHashFunc
}

func NewChunkUpload(hashFunc data.ChunkHashFunc) *ChunkUpload {
	return &ChunkUpload{hashFunc: hashFunc}
}

func (cu *ChunkUpload) Upload(wg *sync.WaitGroup, params *files.UploadChunkParams, uploadErrors chan error) {
	defer wg.Done()

	// TODO: Have a limit of errors before we consider it "not working"?
	resp, err := apiclient.Default.Files.UploadChunk(params)
	if err != nil {
		// TODO: Bug - the channel seems to be closing prematurely when an upload fails... see [ch213]
		uploadErrors <- err
		return
	}

	log.Infof("Response received: %#v", resp.Payload)

	// TODO: Do something here with the response payload to check if the chunk was uploaded correctly
}

func (cu *ChunkUpload) NewParams(chunk data.Chunk, requestID string) (*files.UploadChunkParams, error) {
	file, err := os.Open(chunk.Path)
	if err != nil {
		return nil, err
	}

	readCloser := runtime.NamedReader(chunk.Hash, file)
	params := files.NewUploadChunkParams().
		WithChunkNumber(int64(chunk.Num)).
		WithChunkHash(chunk.Hash).
		WithUploadRequestID(requestID).
		WithChunkData(readCloser)

	return params, nil
}
