package handlers

import (
	"fmt"
	"github.com/go-openapi/runtime/middleware"
	"github.com/patrickmn/go-cache"
	"gitlab.computing.dcu.ie/collint9/2021-ca400-collint9-coynemt2/src/files"
	c "gitlab.computing.dcu.ie/collint9/2021-ca400-collint9-coynemt2/src/server/config"
	"gitlab.computing.dcu.ie/collint9/2021-ca400-collint9-coynemt2/src/server/data"
	"gitlab.computing.dcu.ie/collint9/2021-ca400-collint9-coynemt2/src/server/data/repo"
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

	db := data.DB.Begin()

	chunkDir := schema.ChunkDirectory{
		Path:   fileChunkDir,
		UserID: user.ID,
	}

	if db.Where(chunkDir).FirstOrCreate(&chunkDir).Error != nil {
		db.Rollback()
		return transfer.NewNewUploadDefault(500).WithPayload("Failed to create the chunk directory")
	}

	var storageDirectoryID uint
	if params.UploadInfo.DirectoryID != 0 {
		storageDirectoryID = uint(params.UploadInfo.DirectoryID)
	} else {
		homeDir, err := repo.GetHomeDir(user)
		if err != nil {
			return nil
		}
		storageDirectoryID = homeDir.ID
	}

	file := schema.File{
		Name:               *params.UploadInfo.FileName,
		Size:               *params.UploadInfo.FileSize,
		Hash:               *params.UploadInfo.FileHash,
		ChunkDirectoryID:   chunkDir.ID,
		StorageDirectoryID: storageDirectoryID,
		UserID:             user.ID,
	}

	// prevent users from uploading the same file twice
	var fileWithSameHash schema.File
	tx := db.Where(&schema.File{Hash: file.Hash, UserID: user.ID}).First(&fileWithSameHash)
	if tx.RowsAffected > 0 {
		db.Rollback()
		msg := fmt.Sprintf("you already have this file stored in directory ID: %d", fileWithSameHash.StorageDirectoryID)
		return transfer.NewNewUploadDefault(400).WithPayload(models.Error(msg))
	}

	if db.Create(&file).Error != nil {
		db.Rollback()
		return transfer.NewNewUploadDefault(500).WithPayload("failed to store the file data")
	}

	newUpload := schema.Upload{
		ChunkDirectoryID: chunkDir.ID,
		FileID:           file.ID,
		NumChunks:        *params.UploadInfo.NumChunks,
		UserID:           user.ID,
	}

	if db.Create(&newUpload).Error != nil {
		db.Rollback()
		return transfer.NewNewUploadDefault(500).WithPayload("Failed to add the file info to the database")
	}

	db.Commit()

	data.Cache.Uploads.Set(strconv.Itoa(int(newUpload.ID)), &newUpload, cache.DefaultExpiration)

	return transfer.NewNewUploadOK().WithPayload(&models.Upload{
		ChunkDirectoryID: uint64(newUpload.ChunkDirectoryID),
		FileID:           uint64(newUpload.FileID),
		ID:               uint64(newUpload.ID),
		NumChunks:        newUpload.NumChunks,
	})
}
