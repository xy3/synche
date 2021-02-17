package upload

import (
	"github.com/spf13/afero"
	"gitlab.computing.dcu.ie/collint9/2021-ca400-collint9-coynemt2/src/client/apiclient"
	"gitlab.computing.dcu.ie/collint9/2021-ca400-collint9-coynemt2/src/client/apiclient/files"
	"gitlab.computing.dcu.ie/collint9/2021-ca400-collint9-coynemt2/src/client/models"
)

//go:generate mockery --name=NewUploadRequester --case underscore
type NewUploadRequester interface {
	NewUpload(params *files.NewUploadParams) (accepted *models.NewFileUploadRequestAccepted, err error)
	NewUploadParams(fileSize, numChunks int64, fileHash, fileName string,) *files.NewUploadParams
	CreateNewUpload(file afero.File, fileHash string, numChunks int64) (uploadRequestID string, err error)
}

var DefaultNewUploadRequester = new(NewUploadRequest)

type NewUploadRequest struct {
}

func (nu *NewUploadRequest) CreateNewUpload(file afero.File, fileHash string, numChunks int64) (
	uploadRequestID string,
	err error,
) {
	stat, err := file.Stat()
	if err != nil {
		return uploadRequestID, err
	}
	params := nu.NewUploadParams(stat.Size(), numChunks, fileHash, file.Name())
	requestAccepted, err := nu.NewUpload(params)
	if err != nil {
		return uploadRequestID, err
	}
	return requestAccepted.UploadRequestID, nil
}

func (nu *NewUploadRequest) NewUpload(params *files.NewUploadParams) (*models.NewFileUploadRequestAccepted, error) {
	requestAccepted, err := apiclient.Client.Files.NewUpload(params)
	if err != nil {
		return nil, err
	}
	return requestAccepted.Payload, nil
}

func (nu *NewUploadRequest) NewUploadParams(
	fileSize, numChunks int64,
	fileHash, fileName string,
) *files.NewUploadParams {
	return files.NewNewUploadParams().WithFileInfo(&models.FileInfo{
		Chunks: &numChunks,
		Hash:   &fileHash,
		Name:   &fileName,
		Size:   &fileSize,
	})
}
