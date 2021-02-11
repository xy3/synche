package upload

import (
	"errors"
	"fmt"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"gitlab.computing.dcu.ie/collint9/2021-ca400-collint9-coynemt2/src/client/apiclient"
	"gitlab.computing.dcu.ie/collint9/2021-ca400-collint9-coynemt2/src/client/data"
	"os"
	"path"
	"sync"
)

type FileUploader interface {
	Upload(filePath string) error
}

type FileUpload struct {
	data.Splitter
	ChunkUploader
}

func NewFileUpload(splitter data.Splitter, chunkUploader ChunkUploader) *FileUpload {
	return &FileUpload{Splitter: splitter, ChunkUploader: chunkUploader}
}

func (fu *FileUpload) Upload(filePath string) error {
	info, err := os.Stat(filePath)
	if err != nil {
		log.Error("Could not stat the file at the specified path")
		return err
	}
	if info.IsDir() {
		log.Error("The path specified is a directory: '%s'", filePath)
		return errors.New(fmt.Sprintf("%s refers to a directory when a file path is required", filePath))
	}

	fileName := path.Base(filePath)
	// Get the file name from the --name flag if it is set
	if viper.IsSet("name") && len(viper.GetString("name")) > 0 {
		fileName = viper.GetString("name")
	}

	chunks, err := fu.Splitter.Split(filePath, viper.GetString("ChunkDir"))
	if err != nil {
		return err
	}

	// Send a new file upload request to the server
	newUploadParams, err := NewFileUploadParams(filePath, fileName, int64(len(chunks)))
	if err != nil {
		return err
	}

	requestAccepted, err := apiclient.Default.Files.NewUpload(newUploadParams)
	if err != nil {
		return err
	}

	log.Info("Uploading chunks to the server")
	err = fu.uploadChunks(chunks, requestAccepted.GetPayload().UploadRequestID)
	if err != nil {
		return err
	}

	// TODO: Show some sort of progress report to the user (and do the uploading in the background?)
	log.Info("Finished uploading all chunks to the server")

	return nil
}

func (fu *FileUpload) uploadChunks(chunks []data.Chunk, uploadRequestID string) error {
	var wg sync.WaitGroup
	uploadErrors := make(chan error, len(chunks))

	for _, chunk := range chunks {
		wg.Add(1)
		params, _ := fu.ChunkUploader.NewParams(chunk.Path, chunk.Hash, uploadRequestID, int64(chunk.Num))
		go fu.ChunkUploader.Upload(&wg, params, uploadErrors)
	}

	log.Debug("Waiting for upload workers to finish.")
	wg.Wait()
	log.Debug("All upload workers completed.")

	close(uploadErrors)

	for err := range uploadErrors {
		if err != nil {
			fmt.Printf("%v\n", err)
			return err
		}
	}

	return nil
}
