package handlers

import (
	"github.com/go-openapi/runtime/middleware"
	"gitlab.computing.dcu.ie/collint9/2021-ca400-collint9-coynemt2/src/database"
	"gitlab.computing.dcu.ie/collint9/2021-ca400-collint9-coynemt2/src/database/repo"
	"gitlab.computing.dcu.ie/collint9/2021-ca400-collint9-coynemt2/src/database/schema"
	"gitlab.computing.dcu.ie/collint9/2021-ca400-collint9-coynemt2/src/server/models"
	"gitlab.computing.dcu.ie/collint9/2021-ca400-collint9-coynemt2/src/server/restapi/operations/files"
)

// FileInfo Retrieves and returns to the client information regarding a file specified by ID
func FileInfo(params files.GetFileInfoParams, user *schema.User) middleware.Responder {
	var file models.File
	tx := database.DB.Model(&schema.File{}).Where(&schema.File{UserID: user.ID}).Find(&file, params.FileID)
	if tx.Error != nil {
		return files.NewGetFileInfoNotFound()
	}
	return files.NewGetFileInfoOK().WithPayload(&file)
}

// FilePathInfo Retrieves and returns to the client  information regarding a file specified by path
func FilePathInfo(params files.GetFilePathInfoParams, user *schema.User) middleware.Responder {
	fullPath, err := repo.BuildFullPath(params.FilePath, user, database.DB)

	if err != nil {
		return files.NewGetFilePathInfoNotFound()
	}

	file, err := repo.FindFileByFullPath(fullPath, database.DB)
	if err != nil {
		return files.NewGetFilePathInfoNotFound()
	}

	return files.NewGetFilePathInfoOK().WithPayload(ConvertToFileModel(file))
}
