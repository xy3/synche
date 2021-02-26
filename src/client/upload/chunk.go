package upload

import (
	"bytes"
	"github.com/go-openapi/runtime"
	log "github.com/sirupsen/logrus"
	"gitlab.computing.dcu.ie/collint9/2021-ca400-collint9-coynemt2/src/client/apiclient/files"
	"gitlab.computing.dcu.ie/collint9/2021-ca400-collint9-coynemt2/src/client/config"
	"gitlab.computing.dcu.ie/collint9/2021-ca400-collint9-coynemt2/src/client/data"
	"sync"
)

//go:generate mockery --name=ChunkUploader --case underscore
type ChunkUploader interface {
	Upload(wg *sync.WaitGroup, params *files.UploadChunkParams, uploadErrors chan error)
	NewParams(chunk data.Chunk, requestID string) *files.UploadChunkParams
}

type ChunkUpload struct{}

func (cu *ChunkUpload) Upload(wg *sync.WaitGroup, params *files.UploadChunkParams, uploadErrors chan error) {
	defer wg.Done()

	// TODO: Have a limit of errors before we consider it "not working"?
	resp, err := config.Client.Files.UploadChunk(params)
	if err != nil {
		// TODO: Bug - the channel seems to be closing prematurely when an upload fails... see [ch213]
		uploadErrors <- err
		// Attempt to cast the error as an UploadChunkBadRequest
		if err, ok := err.(*files.UploadChunkBadRequest); ok {
			log.Errorf("%v", *err.Payload.Message)
		} else  {
			// otherwise, just print the Error() as supplied by go-swagger
			log.Errorf("%s", err.Error())
		}
		return
	}

	log.Infof("Successfully uploaded chunk: %s | File ID: %s | Server Directory: %s", resp.Payload.Hash, resp.Payload.CompositeFileID, resp.Payload.DirectoryID)
	// TODO: Do something here with the response payload to check if the chunk was uploaded correctly
}

func (cu *ChunkUpload) NewParams(chunk data.Chunk, requestID string) *files.UploadChunkParams {
	readCloser := runtime.NamedReader(chunk.Hash, bytes.NewReader(*chunk.Bytes))
	return files.NewUploadChunkParams().
		WithChunkNumber(chunk.Num).
		WithChunkHash(chunk.Hash).
		WithUploadRequestID(requestID).
		WithChunkData(readCloser)
}
