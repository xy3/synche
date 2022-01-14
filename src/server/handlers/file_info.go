package handlers

import (
	"github.com/go-openapi/runtime/middleware"
	"github.com/xy3/synche/src/server"
	"github.com/xy3/synche/src/server/models"
	"github.com/xy3/synche/src/server/repo"
	"github.com/xy3/synche/src/server/restapi/operations/files"
	"github.com/xy3/synche/src/server/schema"
)

// FileInfo Retrieves and returns to the client information regarding a file specified by ID
func FileInfo(params files.GetFileInfoParams, user *schema.User) middleware.Responder {
	var file models.File
	tx := server.DB.Model(&schema.File{}).Where(&schema.File{UserID: user.ID}).Find(&file, params.FileID)
	if tx.Error != nil {
		return files.NewGetFileInfoNotFound()
	}
	return files.NewGetFileInfoOK().WithPayload(&file)
}

// FilePathInfo Retrieves and returns to the client  information regarding a file specified by path
func FilePathInfo(params files.GetFilePathInfoParams, user *schema.User) middleware.Responder {
	fullPath, err := repo.BuildFullPath(params.FilePath, user, server.DB)

	if err != nil {
		return files.NewGetFilePathInfoNotFound()
	}

	file, err := repo.FindFileByFullPath(fullPath, server.DB)
	if err != nil {
		return files.NewGetFilePathInfoNotFound()
	}

	return files.NewGetFilePathInfoOK().WithPayload(file.ConvertToFileModel())
}
