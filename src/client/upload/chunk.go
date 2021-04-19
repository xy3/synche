package upload

import (
	"bytes"
	"context"
	"github.com/go-openapi/runtime"
	log "github.com/sirupsen/logrus"
	"gitlab.computing.dcu.ie/collint9/2021-ca400-collint9-coynemt2/src/client/apiclient"
	"gitlab.computing.dcu.ie/collint9/2021-ca400-collint9-coynemt2/src/client/apiclient/transfer"
	"gitlab.computing.dcu.ie/collint9/2021-ca400-collint9-coynemt2/src/client/data"
	"sync"
)

//go:generate mockery --name=AsyncChunkUploader --case underscore
type AsyncChunkUploader func(wg *sync.WaitGroup, params *transfer.UploadChunkParams, uploadErrors chan error)

func AsyncChunkUpload(wg *sync.WaitGroup, params *transfer.UploadChunkParams, uploadErrors chan error) {
	defer wg.Done()

	// TODO: Have a limit of errors before we consider it "not working"

	resp, err := apiclient.Client.Transfer.UploadChunk(params, apiclient.ClientAuth)
	if err != nil {
		uploadErrors <- err
		log.Error(err)
		return
	}

	chunk := resp.Payload
	log.WithFields(log.Fields{
		"hash":         chunk.Chunk.Hash,
		"file_id":      chunk.FileID,
		"directory_id": chunk.DirectoryID,
	}).Debug("Successfully uploaded chunk")

	// TODO: Do something here with the response payload to check if the chunk was uploaded correctly
	return
}

func NewChunkUploadParams(chunk data.Chunk, uploadID uint64) *transfer.UploadChunkParams {
	chunkData := runtime.NamedReader(chunk.Hash, bytes.NewReader(*chunk.Bytes))
	return &transfer.UploadChunkParams{
		ChunkData:   chunkData,
		ChunkHash:   chunk.Hash,
		ChunkNumber: chunk.Num,
		UploadID:    uploadID,
		Context:     context.TODO(),
	}
}
