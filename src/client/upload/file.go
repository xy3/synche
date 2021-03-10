package upload

import (
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"gitlab.computing.dcu.ie/collint9/2021-ca400-collint9-coynemt2/src/client/data"
	"sync"
)

//go:generate mockery --name=FileUploader --case underscore
type FileUploader interface {
	AsyncUpload(splitter data.Splitter) error
	SyncUpload(splitter data.Splitter) error
}

type FileUpload struct {
	ChunkUploader      ChunkUploader
	NewUploadRequester NewUploadRequester
}

func NewFileUpload(chunkUploader ChunkUploader, newUploadRequester NewUploadRequester) *FileUpload {
	return &FileUpload{ChunkUploader: chunkUploader, NewUploadRequester: newUploadRequester}
}


func (fu *FileUpload) SyncUpload(splitter data.Splitter) error {
	uploadRequestID, err := fu.NewUploadRequester.CreateNewUpload(splitter)
	if err != nil {
		return err
	}

	// The anonymous func here is called everytime a new chunk is read from the file
	err = splitter.Split(
		func(chunk *data.Chunk) error {
			params := fu.ChunkUploader.NewParams(*chunk, uploadRequestID)
			return fu.ChunkUploader.SyncUpload(params)
		},
	)
	if err != nil {
		return err
	}
	log.Infof("Finished uploading all %d chunks to the server", splitter.NumChunks())
	return nil
}


func (fu *FileUpload) AsyncUpload(splitter data.Splitter) error {
	var wg sync.WaitGroup
	uploadErrors := make(chan error)

	uploadRequestID, err := fu.NewUploadRequester.CreateNewUpload(splitter)
	if err != nil {
		return err
	}

	// The anonymous func here is called everytime a new chunk is read from the file
	workers := viper.GetInt("config.chunks.workers")
	chunkSize := viper.GetInt("config.chunks.size")
	log.WithFields(log.Fields{"workers": workers, "chunksize": chunkSize}).Info("Chunk config")
	// splitFile := splitter.File()
	err = splitter.Split(
		func(chunk *data.Chunk) error {
			idx := splitter.File().CurrentIndex
			if idx % int64(workers) == 0 {
				log.Infof("%d - Waiting for %d workers...", idx, workers)
				wg.Wait()
			}
			params := fu.ChunkUploader.NewParams(*chunk, uploadRequestID)
			wg.Add(1)
			go fu.ChunkUploader.AsyncUpload(&wg, params, uploadErrors)
			return nil
		},
	)
	if err != nil {
		return err
	}

	log.Info("Waiting for upload workers to finish.")
	wg.Wait()
	log.Infof("Finished uploading all %d chunks to the server", splitter.NumChunks())
	close(uploadErrors)

	// Here we could attempt to cast the error as an UploadChunkBadRequest or other relevant error types
	for err := range uploadErrors {
		if err != nil {
			log.Error(err)
			return err
		}
	}

	return nil
}
