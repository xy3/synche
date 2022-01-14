package handlers

import (
	"fmt"
	"github.com/go-openapi/runtime/middleware"
	"github.com/xy3/synche/src/server"
	"github.com/xy3/synche/src/server/models"
	"github.com/xy3/synche/src/server/repo"
	"github.com/xy3/synche/src/server/restapi/operations/transfer"
	"github.com/xy3/synche/src/server/schema"
)

// createNewUploadAndFile Stores the file data relating to a file that has been requested to be uploaded
func createNewUploadAndFile(directoryID uint, params transfer.NewUploadParams, user *schema.User) middleware.Responder {
	db := server.DB.Begin()

	// TODO: send the chunk size in the upload request
	file := schema.File{
		Name:        params.FileName,
		Size:        params.FileSize,
		Hash:        params.FileHash,
		DirectoryID: directoryID,
		UserID:      user.ID,
		TotalChunks: params.NumChunks,
	}

	if db.Where(file).FirstOrCreate(&file).Error != nil {
		db.Rollback()
		return transfer.NewNewUploadDefault(500).WithPayload("failed to store the file data, file already exists")
	}

	if (file.ChunksReceived == file.TotalChunks) && file.Available {
		return transfer.NewNewUploadDefault(500).WithPayload("failed to store the file data, file already exists")
	}

	db.Commit()

	if err := repo.UpdateDirFileCount(directoryID, server.DB); err != nil {
		return transfer.NewNewUploadDefault(500).WithPayload("failed to update the directory file count")
	}

	return transfer.NewNewUploadOK().WithPayload(file.ConvertToFileModel())
}

// NewUpload Handles requests to upload a new file and responds to the client
func NewUpload(params transfer.NewUploadParams, user *schema.User) middleware.Responder {
	var (
		err         error
		directoryID uint
		directory   *schema.Directory
	)

	if params.DirectoryID != nil && *params.DirectoryID != 0 {
		directoryID = uint(*params.DirectoryID)
		directory, err = repo.GetDirectoryByID(directoryID, server.DB)
		if err != nil {
			return transfer.NewNewUploadDefault(500).WithPayload("directory not found")
		}
	} else {
		directory, err = repo.GetHomeDir(user.ID, server.DB)
		if err != nil {
			return transfer.NewNewUploadDefault(500).WithPayload("home directory not found")
		}
		directoryID = directory.ID
	}

	// prevent users from uploading the same file twice
	var prevFile schema.File
	tx := server.DB.Joins("Upload").Where(&schema.File{
		UserID: user.ID,
		Hash:   params.FileHash,
	}).First(&prevFile)

	if tx.Error != nil {
		return createNewUploadAndFile(directoryID, params, user)
	}

	msg := fmt.Sprintf("you already have this file stored in directory ID: %d", prevFile.DirectoryID)
	errAlreadyExists := transfer.NewNewUploadDefault(400).WithPayload(models.Error(msg))

	if err := prevFile.ValidateHash(server.DB); err == nil {
		return errAlreadyExists
	}

	if prevFile.Available {
		if err := prevFile.Delete(server.DB); err != nil {
			return transfer.NewNewUploadDefault(500).WithPayload("failed to remove old invalid file")
		}
		if err := prevFile.Delete(server.DB); err != nil {
			return transfer.NewNewUploadDefault(500).WithPayload("failed to remove old invalid upload")
		}
		return createNewUploadAndFile(directoryID, params, user)
	}

	// repo.UploadsCache.Set(strconv.Itoa(int(newUpload.ID)), &newUpload, cache.DefaultExpiration)

	return transfer.NewNewUploadOK().WithPayload(prevFile.ConvertToFileModel())
}
