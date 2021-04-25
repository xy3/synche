package upload

import (
	"gitlab.computing.dcu.ie/collint9/2021-ca400-collint9-coynemt2/src/client/apiclient"
	"gitlab.computing.dcu.ie/collint9/2021-ca400-collint9-coynemt2/src/client/apiclient/transfer"
	"gitlab.computing.dcu.ie/collint9/2021-ca400-collint9-coynemt2/src/client/data"
	"gitlab.computing.dcu.ie/collint9/2021-ca400-collint9-coynemt2/src/client/models"
)

//go:generate mockery --name=NewUploadFunc --case underscore
type NewUploadFunc func(params *transfer.NewUploadParams) (*models.Upload, error)

func NewUpload(params *transfer.NewUploadParams) (*models.Upload, error) {
	newUpload, err := apiclient.Client.Transfer.NewUpload(params, apiclient.ClientAuth)
	if err != nil {
		return nil, err
	}
	return newUpload.GetPayload(), nil
}

func NewUploadParams(
	fileSize, numChunks int64,
	fileHash, fileName string,
	uploadDirID uint64,
) *transfer.NewUploadParams {
	return transfer.NewNewUploadParams().WithUploadInfo(
		&models.NewFileUpload{
			DirectoryID: uploadDirID,
			FileHash:    &fileHash,
			FileName:    &fileName,
			FileSize:    &fileSize,
			NumChunks:   &numChunks,
		},
	)
}

func NewUploadParamsFromSplitter(splitter data.Splitter, uploadDirID uint) *transfer.NewUploadParams {
	return NewUploadParams(splitter.File().FileSize, splitter.NumChunks(), splitter.File().Hash, splitter.File().Name, uint64(uploadDirID))
}
