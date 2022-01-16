package handlers

import (
	"github.com/go-openapi/runtime/middleware"
	schema2 "github.com/xy3/synche/src/schema"
	"github.com/xy3/synche/src/server"
	"github.com/xy3/synche/src/server/models"
	"github.com/xy3/synche/src/server/repo"
	"github.com/xy3/synche/src/server/restapi/operations/files"
)

// FileInfo Retrieves and returns to the client information regarding a file specified by ID
func FileInfo(params files.GetFileInfoParams, user *schema2.User) middleware.Responder {
	var file models.File
	tx := server.DB.Model(&schema2.File{}).Where(&schema2.File{UserID: user.ID}).Find(&file, params.FileID)
	if tx.Error != nil {
		return files.NewGetFileInfoNotFound()
	}
	return files.NewGetFileInfoOK().WithPayload(&file)
}

// FilePathInfo Retrieves and returns to the client  information regarding a file specified by path
func FilePathInfo(params files.GetFilePathInfoParams, user *schema2.User) middleware.Responder {
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
