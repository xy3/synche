package upload

import (
	"context"
	"errors"
	"gitlab.computing.dcu.ie/collint9/2021-ca400-collint9-coynemt2/src/client/apiclient"
	"gitlab.computing.dcu.ie/collint9/2021-ca400-collint9-coynemt2/src/client/apiclient/transfer"
	"gitlab.computing.dcu.ie/collint9/2021-ca400-collint9-coynemt2/src/client/data"
	"gitlab.computing.dcu.ie/collint9/2021-ca400-collint9-coynemt2/src/client/models"
)

//go:generate mockery --name=NewUploadFunc --case underscore
type NewUploadFunc func(params *transfer.NewUploadParams) (*models.File, error)

// NewUpload Sends a request to the server that checks if it can begin uploading a file
func NewUpload(params *transfer.NewUploadParams) (*models.File, error) {
	if params.NumChunks == 0 {
		return nil, errors.New("cannot upload with: NumChunks = 0")
	}
	newUploadFile, err := apiclient.Client.Transfer.NewUpload(params, apiclient.ClientAuth)
	if err != nil {
		return nil, err
	}

	return newUploadFile.GetPayload(), nil
}

// NewUploadParams Creates an upload object to send to the server
func NewUploadParams(
	fileSize, numChunks int64,
	fileHash, fileName string,
	uploadDirID uint64,
) *transfer.NewUploadParams {
	return &transfer.NewUploadParams{
		Context:     context.Background(),
		DirectoryID: &uploadDirID,
		FileHash:    fileHash,
		FileName:    fileName,
		FileSize:    fileSize,
		NumChunks:   numChunks,
	}
}

// NewUploadParamsFromSplitter Creates an upload object to send to the server
func NewUploadParamsFromSplitter(splitter data.Splitter, uploadDirID uint) *transfer.NewUploadParams {
	return NewUploadParams(
		splitter.File().FileSize,
		splitter.NumChunks(),
		splitter.File().Hash,
		splitter.File().Name,
		uint64(uploadDirID),
	)
}
