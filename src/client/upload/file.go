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

//go:generate mockery --name=FileUploader --case underscore
type FileUploader interface {
	Upload(file *os.File) error
	uploadChunks(chunks []data.Chunk, uploadRequestID string) error
}

type FileUpload struct {
	SplitJob      data.SplitJob
	ChunkUploader ChunkUploader
	FileHashFunc  data.FileHashFunc
}

func NewFileUpload(splitJob data.SplitJob, chunkUploader ChunkUploader, fileHashFunc data.FileHashFunc) *FileUpload {
	return &FileUpload{SplitJob: splitJob, ChunkUploader: chunkUploader, FileHashFunc: fileHashFunc}
}

func (fu *FileUpload) Upload(file *os.File) error {
	stat, err := file.Stat()
	if err != nil {
		log.Error("Could not stat the file at the specified path")
		return err
	}
	if stat.IsDir() {
		log.Error("The path specified is a directory: '%s'", file.Name())
		return errors.New(fmt.Sprintf("%s refers to a directory when a file path is required", file.Name()))
	}
	fileName := path.Base(file.Name())
	// Get the file name from the --name flag if it is set
	if viper.IsSet("name") && len(viper.GetString("name")) > 0 {
		fileName = viper.GetString("name")
	}

	chunks, err := fu.SplitJob.Split(file)
	if err != nil {
		return err
	}

	// Send a new file upload request to the server
	fileHash, err := fu.FileHashFunc(file.Name())
	if err != nil {
		return err
	}

	numChunks := len(*chunks)
	fileSize := stat.Size()
	newUploadParams := NewFileUploadParams(fileHash, fileName, int64(numChunks), fileSize)

	requestAccepted, err := apiclient.Default.Files.NewUpload(newUploadParams)
	if err != nil {
		return err
	}

	log.Info("Uploading chunks to the server")
	err = fu.uploadChunks(*chunks, requestAccepted.GetPayload().UploadRequestID)
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
		params, err := fu.ChunkUploader.NewParams(chunk, uploadRequestID)
		if err != nil {
			panic(err)
		}
		go fu.ChunkUploader.Upload(&wg, params, uploadErrors)
	}

	// return when an error is encountered
	err := <-uploadErrors
	if err != nil {
		wg.Wait()
		close(uploadErrors)
		return err
	}

	//// Attempt to cast the error as an UploadChunkBadRequest
	//if err, ok := err.(*files.UploadChunkBadRequest); ok {
	//	log.Errorf("%v", *err.Payload.Message)
	//} else {
	//	// otherwise, just print the Error() as supplied by go-swagger
	//	log.Errorf("%s", err.Error())
	//}

	log.Infof("Waiting for upload workers to finish.")
	wg.Wait()
	log.Infof("All upload workers completed.")
	close(uploadErrors)

	for err := range uploadErrors {
		if err != nil {
			log.Errorf("%v\n", err)
			return err
		}
	}

	return nil
}
