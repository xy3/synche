package upload

import (
	"gitlab.computing.dcu.ie/collint9/2021-ca400-collint9-coynemt2/src/client/apiclient/transfer"
	"gitlab.computing.dcu.ie/collint9/2021-ca400-collint9-coynemt2/src/client/config"
	"gitlab.computing.dcu.ie/collint9/2021-ca400-collint9-coynemt2/src/client/data"
	"gitlab.computing.dcu.ie/collint9/2021-ca400-collint9-coynemt2/src/client/models"
)

//go:generate mockery --name=NewUploadRequester --case underscore
type NewUploadRequester interface {
	CreateNewUpload(splitter data.Splitter) (uploadRequestID string, err error)
	NewUpload(params *transfer.NewUploadParams) (accepted *models.NewFileUploadRequestAccepted, err error)
	NewUploadParams(fileSize, numChunks int64, fileHash, fileName string,) *transfer.NewUploadParams
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

func (nu *NewUploadRequest) NewUpload(params *transfer.NewUploadParams) (*models.NewFileUploadRequestAccepted, error) {
	requestAccepted, err := config.Client.Transfer.NewUpload(params)
	if err != nil {
		return nil, err
	}
	return requestAccepted.GetPayload(), nil
}

func (nu *NewUploadRequest) NewUploadParams(fileSize, numChunks int64, fileHash, fileName string,) *transfer.NewUploadParams {
	return transfer.NewNewUploadParams().WithFileInfo(
		&models.FileInfo{
			Chunks: &numChunks,
			Hash:   &fileHash,
			Name:   &fileName,
			Size:   &fileSize,
		},
	)
}
