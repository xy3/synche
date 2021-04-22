package handlers

import (
	"github.com/go-openapi/runtime/middleware"
	"github.com/patrickmn/go-cache"
	"gitlab.computing.dcu.ie/collint9/2021-ca400-collint9-coynemt2/src/files"
	c "gitlab.computing.dcu.ie/collint9/2021-ca400-collint9-coynemt2/src/server/config"
	"gitlab.computing.dcu.ie/collint9/2021-ca400-collint9-coynemt2/src/server/data"
	"gitlab.computing.dcu.ie/collint9/2021-ca400-collint9-coynemt2/src/server/data/schema"
	"gitlab.computing.dcu.ie/collint9/2021-ca400-collint9-coynemt2/src/server/models"
	"gitlab.computing.dcu.ie/collint9/2021-ca400-collint9-coynemt2/src/server/restapi/operations/transfer"
	"path/filepath"
	"strconv"
)

func NewUpload(
	params transfer.NewUploadParams,
	user *schema.User,
) middleware.Responder {
	// Make a directory in the upload dir with the hash as the name
	fileChunkDir := filepath.Join(c.Config.Server.UploadDir, *params.UploadInfo.FileHash)
	err := files.AppFS.MkdirAll(fileChunkDir, 0755)
	if err != nil {
		return transfer.NewNewUploadDefault(500).WithPayload("Failed to create a directory for the file")
	}

	tx := data.DB.Begin()

	fileDir := schema.Directory{
		Path: fileChunkDir,
	}

	tx.Where(fileDir).FirstOrCreate(&fileDir)
	if tx.Error != nil {
		tx.Rollback()
		return transfer.NewNewUploadDefault(500).WithPayload("Failed to create the chunk directory")
	}

	newUpload := schema.Upload{
		ChunkDirectory: fileDir,
		File: schema.File{
			Name:               *params.UploadInfo.FileName,
			Size:               *params.UploadInfo.FileSize,
			Hash:               *params.UploadInfo.FileHash,
			ChunkDirectoryID:   fileDir.ID,
			StorageDirectoryID: uint(1), // default to home (config) storage directory
			UserID:             user.ID,
		},
		NumChunks: *params.UploadInfo.NumChunks,
	}

	tx.Create(&newUpload)

	if tx.Error != nil {
		tx.Rollback()
		return transfer.NewNewUploadDefault(500).WithPayload("Failed to add the file info to the database")
	}

	tx.Commit()

	data.Cache.Uploads.Set(strconv.Itoa(int(newUpload.ID)), &newUpload, cache.DefaultExpiration)

	return transfer.NewNewUploadOK().WithPayload(&models.Upload{
		ChunkDirectoryID: uint64(newUpload.ChunkDirectoryID),
		FileID:           uint64(newUpload.FileID),
		ID:               uint64(newUpload.ID),
		NumChunks:        newUpload.NumChunks,
	})
}
