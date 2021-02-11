package handlers

import (
	"database/sql"
	"github.com/go-openapi/runtime/middleware"
	"github.com/spf13/viper"
	"gitlab.computing.dcu.ie/collint9/2021-ca400-collint9-coynemt2/src/server/database"
	"gitlab.computing.dcu.ie/collint9/2021-ca400-collint9-coynemt2/src/server/models"
	"gitlab.computing.dcu.ie/collint9/2021-ca400-collint9-coynemt2/src/server/restapi/operations/files"
	"log"
	"os"
	"path/filepath"
)

func NewUploadFileHandler(params files.NewUploadParams, db *sql.DB) middleware.Responder {
	// TODO: Check the file info here e.g. verify the hash

	//requestUuid := uuid.New().String() could use uuid?
	uploadRequestId := *params.FileInfo.Hash // just use the file hash for the moment

	// Make a directory in /data/received with the hash as the name
	fileChunkDir := filepath.Join(viper.GetString("server.uploadDirectory"), uploadRequestId)
	_ = os.MkdirAll(fileChunkDir, os.ModePerm)

	// Store upload request ID, chunk directory, file name, file size, and number of chunks in the database
	err := database.InsertConnectionRequest(db, uploadRequestId, fileChunkDir, *params.FileInfo.Name, *params.FileInfo.Size, *params.FileInfo.Chunks)
	if err != nil {
		log.Printf("Could not insert connection request into database with ID: %s", uploadRequestId)
	} else {
		log.Printf("Inserted into database connection request information with ID: %s", uploadRequestId)
	}

	contents := models.NewFileUploadRequestAccepted{
		UploadRequestID: uploadRequestId,
	}
	return files.NewNewUploadOK().WithPayload(&contents)
}
