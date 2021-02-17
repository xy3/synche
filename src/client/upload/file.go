package upload

import (
	log "github.com/sirupsen/logrus"
	"github.com/spf13/afero"
	"gitlab.computing.dcu.ie/collint9/2021-ca400-collint9-coynemt2/src/client/data"
)

//go:generate mockery --name=FileUploader --case underscore
type FileUploader interface {
	Upload(file afero.File) error
}

type FileUpload struct {
	Splitter      data.Splitter
	ChunkUploader ChunkUploader
	FileHashFunc data.FileHashFunc
	NewUploadRequester NewUploadRequester
}

func NewFileUpload(
	splitter data.Splitter,
	chunkUploader ChunkUploader,
	fileHashFunc data.FileHashFunc,
	newUploadRequester NewUploadRequester,
) *FileUpload {
	return &FileUpload{
		Splitter: splitter,
		ChunkUploader: chunkUploader,
		FileHashFunc: fileHashFunc,
		NewUploadRequester: newUploadRequester,
	}
}


func (fu *FileUpload) Upload(file afero.File) error {
	chunks, err := fu.Splitter.Split(file)
	if err != nil {
		return err
	}

	fileHash, err := fu.FileHashFunc(file.Name())
	if err != nil {
		return err
	}

	uploadRequestID, err := fu.NewUploadRequester.CreateNewUpload(file, fileHash, int64(len(*chunks)))
	if err != nil {
		return err
	}

	log.Info("Uploading chunks to the server")
	err = fu.ChunkUploader.UploadChunks(*chunks, uploadRequestID)
	if err != nil {
		return err
	}

	// TODO: Show some sort of progress report to the user (and do the uploading in the background?)
	log.Info("Finished uploading all chunks to the server")

	return nil
}
