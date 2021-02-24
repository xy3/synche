package upload

import (
	log "github.com/sirupsen/logrus"
	"gitlab.computing.dcu.ie/collint9/2021-ca400-collint9-coynemt2/src/client/data"
	"sync"
)

//go:generate mockery --name=FileUploader --case underscore
type FileUploader interface {
	Upload(splitter data.Splitter) error
}

type FileUpload struct {
	ChunkUploader      ChunkUploader
	NewUploadRequester NewUploadRequester
}

func NewFileUpload(chunkUploader ChunkUploader, newUploadRequester NewUploadRequester) *FileUpload {
	return &FileUpload{ChunkUploader: chunkUploader, NewUploadRequester: newUploadRequester}
}

func (fu *FileUpload) Upload(splitter data.Splitter) error {
	var wg sync.WaitGroup
	uploadErrors := make(chan error, splitter.NumChunks())

	uploadRequestID, err := fu.NewUploadRequester.CreateNewUpload(splitter)
	if err != nil {
		return err
	}

	// The anonymous func here is called everytime a new chunk is read from the file
	err = splitter.Split(
		func(chunk *data.Chunk) error {
			params := fu.ChunkUploader.NewParams(*chunk, uploadRequestID)
			wg.Add(1)
			go fu.ChunkUploader.Upload(&wg, params, uploadErrors)
			return nil
		},
	)
	if err != nil {
		return err
	}

	log.Info("Waiting for upload workers to finish.")
	wg.Wait()
	log.Info("Finished uploading all chunks to the server")
	close(uploadErrors)

	// Here we could attempt to cast the error as an UploadChunkBadRequest or other relevant error types
	for err := range uploadErrors {
		if err != nil {
			log.Errorf("%v\n", err)
			return err
		}
	}

	return nil
}
