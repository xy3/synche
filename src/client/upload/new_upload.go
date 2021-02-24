package upload

import (
	"gitlab.computing.dcu.ie/collint9/2021-ca400-collint9-coynemt2/src/client/apiclient"
	"gitlab.computing.dcu.ie/collint9/2021-ca400-collint9-coynemt2/src/client/apiclient/files"
	"gitlab.computing.dcu.ie/collint9/2021-ca400-collint9-coynemt2/src/client/data"
	"gitlab.computing.dcu.ie/collint9/2021-ca400-collint9-coynemt2/src/client/models"
)

//go:generate mockery --name=NewUploadRequester --case underscore
type NewUploadRequester interface {
	CreateNewUpload(splitter data.Splitter) (uploadRequestID string, err error)
	NewUpload(params *files.NewUploadParams) (accepted *models.NewFileUploadRequestAccepted, err error)
	NewUploadParams(fileSize, numChunks int64, fileHash, fileName string,) *files.NewUploadParams
}

type NewUploadRequest struct {}

var DefaultNewUploadRequester = new(NewUploadRequest)

func (nu *NewUploadRequest) CreateNewUpload(splitter data.Splitter) (uploadRequestID string, err error) {
	params := nu.NewUploadParams(splitter.File().FileSize, splitter.NumChunks(), splitter.File().Hash, splitter.File().Name)
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

func (nu *NewUploadRequest) NewUploadParams(fileSize, numChunks int64, fileHash, fileName string,) *files.NewUploadParams {
	return files.NewNewUploadParams().WithFileInfo(
		&models.FileInfo{
			Chunks: &numChunks,
			Hash:   &fileHash,
			Name:   &fileName,
			Size:   &fileSize,
		},
	)
}
