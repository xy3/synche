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
	"time"
)

//go:generate mockery --name=AsyncChunkUploader --case underscore
type AsyncChunkUploader func(wg *sync.WaitGroup, params *transfer.UploadChunkParams, uploadErrors chan error)

func AsyncChunkUpload(wg *sync.WaitGroup, params *transfer.UploadChunkParams, uploadErrors chan error) {
	defer wg.Done()

	timeout := make(chan bool, 1)
	go func() {
		time.Sleep(10 * time.Second)
		timeout <- true
	}()

	resp, err := apiclient.Client.Transfer.UploadChunk(params, apiclient.ClientAuth)
	if err != nil {
		uploadErrors <- err
		log.Error(err)
		return
	}

	chunk := resp.Payload
	log.WithFields(log.Fields{
		"Hash":     chunk.Chunk.Hash,
		"FileID":   chunk.FileID,
		"ChunkNum": chunk.Number,
		"ChunkID":  chunk.Chunk.ID,
	}).Debug("Successfully uploaded chunk")

	return
}

func NewChunkUploadParams(chunk data.Chunk, fileID uint64) *transfer.UploadChunkParams {
	chunkData := runtime.NamedReader(chunk.Hash, bytes.NewReader(*chunk.Bytes))
	return &transfer.UploadChunkParams{
		ChunkData:   chunkData,
		ChunkHash:   chunk.Hash,
		ChunkNumber: chunk.Num,
		FileID:      fileID,
		Context:     context.TODO(),
	}
}
