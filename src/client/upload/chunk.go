package upload

import (
	"bytes"
	"github.com/go-openapi/runtime"
	log "github.com/sirupsen/logrus"
	"gitlab.computing.dcu.ie/collint9/2021-ca400-collint9-coynemt2/src/client/apiclient/transfer"
	"gitlab.computing.dcu.ie/collint9/2021-ca400-collint9-coynemt2/src/client/config"
	"gitlab.computing.dcu.ie/collint9/2021-ca400-collint9-coynemt2/src/client/data"
	"sync"
)

//go:generate mockery --name=ChunkUploader --case underscore
type ChunkUploader interface {
	AsyncUpload(wg *sync.WaitGroup, params *transfer.UploadChunkParams, uploadErrors chan error)
	NewParams(chunk data.Chunk, requestID int64) *transfer.UploadChunkParams
	SyncUpload(params *transfer.UploadChunkParams) error
}

type ChunkUpload struct{}

func (cu *ChunkUpload) SyncUpload(params *transfer.UploadChunkParams) error {
	resp, err := config.Client.Transfer.UploadChunk(params)
	if err != nil {
		log.Error(err)
		return err
	}
	chunk := resp.Payload
	log.WithFields(log.Fields{
		"hash": chunk.Hash,
		"file_id": chunk.FileID,
		"directory_id": chunk.DirectoryID,
	}).Debug("Successfully uploaded chunk")
	return nil
}

func (cu *ChunkUpload) AsyncUpload(wg *sync.WaitGroup, params *transfer.UploadChunkParams, uploadErrors chan error) {
	defer wg.Done()

	// TODO: Have a limit of errors before we consider it "not working"
	resp, err := config.Client.Transfer.UploadChunk(params)
	if err != nil {
		uploadErrors <- err
		log.Error(err)
		return
	}

	chunk := resp.Payload
	log.WithFields(log.Fields{
		"hash":         chunk.Hash,
		"file_id":      chunk.FileID,
		"directory_id": chunk.DirectoryID,
	}).Debug("Successfully uploaded chunk")

	// TODO: Do something here with the response payload to check if the chunk was uploaded correctly
}

func (cu *ChunkUpload) NewParams(chunk data.Chunk, uploadId int64) *transfer.UploadChunkParams {
	readCloser := runtime.NamedReader(chunk.Hash, bytes.NewReader(*chunk.Bytes))
	return transfer.NewUploadChunkParams().
		WithChunkNumber(chunk.Num).
		WithChunkHash(chunk.Hash).
		WithUploadRequestID(uploadId).
		WithChunkData(readCloser)
}
