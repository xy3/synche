package handlers

import (
	"github.com/go-openapi/runtime/middleware"
	"gitlab.computing.dcu.ie/collint9/2021-ca400-collint9-coynemt2/src/server/models"
	"gitlab.computing.dcu.ie/collint9/2021-ca400-collint9-coynemt2/src/server/restapi/operations/files"
	"os"
	"path/filepath"
)

func NewUploadFileHandler(params files.NewUploadParams) middleware.Responder {
	// TODO: Check the file info here e.g. verify the hash
	// store info in the database such as the number of expected chunks

	// TODO: IMPORTANT Add this ID to the database
	//requestUuid := uuid.New().String() could use uuid?
	uploadRequestId := *params.FileInfo.Hash // just use the file hash for the moment

	// Make a directory in /data/received with the hash as the name
	fileChunkDir := filepath.Join(UploadDirectory, uploadRequestId)
	_ = os.MkdirAll(fileChunkDir, os.ModePerm)

	contents := models.NewFileUploadRequestAccepted{
		UploadRequestID: uploadRequestId,
	}
	return files.NewNewUploadOK().WithPayload(&contents)
}
