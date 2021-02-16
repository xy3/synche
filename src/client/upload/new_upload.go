package upload

import (
	"gitlab.computing.dcu.ie/collint9/2021-ca400-collint9-coynemt2/src/client/apiclient/files"
	"gitlab.computing.dcu.ie/collint9/2021-ca400-collint9-coynemt2/src/client/models"
)

func NewFileUploadParams(fileHash, fileName string, numChunks, fileSize int64) *files.NewUploadParams {
	return files.NewNewUploadParams().WithFileInfo(&models.FileInfo{
		Chunks: &numChunks,
		Hash:   &fileHash,
		Name:   &fileName,
		Size:   &fileSize,
	})
}