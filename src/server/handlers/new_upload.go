package handlers

import (
	"github.com/go-openapi/runtime/middleware"
	c "gitlab.computing.dcu.ie/collint9/2021-ca400-collint9-coynemt2/src/server/config"
	"gitlab.computing.dcu.ie/collint9/2021-ca400-collint9-coynemt2/src/server/data"
	"gitlab.computing.dcu.ie/collint9/2021-ca400-collint9-coynemt2/src/server/models"
	"gitlab.computing.dcu.ie/collint9/2021-ca400-collint9-coynemt2/src/server/restapi/operations/files"
	"os"
	"path/filepath"
)

func NewUploadFileHandler(params files.NewUploadParams, dataAccess data.SyncheData) middleware.Responder {
	// TODO: Check the file info here e.g. verify the hash

	// requestUuid := uuid.New().String() could use uuid?
	uploadRequestId := *params.FileInfo.Hash // just use the file hash for the moment

	// Make a directory in .synche/data/received with the hash as the name
	fileChunkDir := filepath.Join(c.Config.Server.UploadDir, uploadRequestId)
	_ = os.MkdirAll(fileChunkDir, os.ModePerm)

	// Store upload request ID, chunk directory, file name, file size, and number of chunks in the data
	err := data.Database.InsertConnectionRequest(dataAccess.Database, uploadRequestId, fileChunkDir, filepath.Base(*params.FileInfo.Name), *params.FileInfo.Size, *params.FileInfo.Chunks)
	if err != nil {
		return files.NewNewUploadBadRequest()
	}

	err = dataAccess.Cache.SetNumberOfChunks(uploadRequestId, *params.FileInfo.Chunks)
	if err != nil {
		return files.NewNewUploadBadRequest()
	}

	contents := models.NewFileUploadRequestAccepted{
		UploadRequestID: uploadRequestId,
	}
	return files.NewNewUploadOK().WithPayload(&contents)
}
