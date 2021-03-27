package handlers

import (
	"github.com/go-openapi/runtime/middleware"
	log "github.com/sirupsen/logrus"
	"gitlab.computing.dcu.ie/collint9/2021-ca400-collint9-coynemt2/src/files"
	c "gitlab.computing.dcu.ie/collint9/2021-ca400-collint9-coynemt2/src/server/config"
	"gitlab.computing.dcu.ie/collint9/2021-ca400-collint9-coynemt2/src/server/data"
	"gitlab.computing.dcu.ie/collint9/2021-ca400-collint9-coynemt2/src/server/data/schema"
	"gitlab.computing.dcu.ie/collint9/2021-ca400-collint9-coynemt2/src/server/models"
	"gitlab.computing.dcu.ie/collint9/2021-ca400-collint9-coynemt2/src/server/restapi/operations/transfer"
	"path/filepath"
)

func NewUploadFileHandler(syncheData *data.SyncheData, params transfer.NewUploadParams) middleware.Responder {
	// Make a directory in the upload dir with the hash as the name
	fileChunkDir := filepath.Join(c.Config.Server.UploadDir, *params.FileInfo.Hash)
	err := files.AppFS.MkdirAll(fileChunkDir, 0755)
	if err != nil {
		return transfer.NewNewUploadBadRequest().WithPayload("failed to create a directory for the file")
	}

	fileDir := schema.Directory{
		Path: fileChunkDir,
	}

	tx := syncheData.DB.Begin()
	newUpload := schema.Upload{
		Directory: fileDir,
		File: schema.File{
			Name:      *params.FileInfo.Name,
			Size:      *params.FileInfo.Size,
			Hash:      *params.FileInfo.Hash,
			Directory: fileDir,
		},
		NumChunks: *params.FileInfo.Chunks,
	}

	tx.Create(&newUpload)

	if tx.Error != nil {
		tx.Rollback()
		return transfer.NewNewUploadBadRequest().WithPayload("failed to add the file info to the database")
	}

	tx.Commit()

	err = syncheData.Cache.UploadCache.SetUpload(syncheData.Cache, newUpload.ID, newUpload)
	if err != nil {
		log.WithError(err).Error("Failed to cache the upload data")
	}

	return transfer.NewNewUploadOK().WithPayload(&models.NewFileUploadRequestAccepted{
		UploadRequestID: int64(newUpload.ID),
	})
}
