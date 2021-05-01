package handlers

import (
	"github.com/go-openapi/runtime/middleware"
	"gitlab.computing.dcu.ie/collint9/2021-ca400-collint9-coynemt2/src/database"
	"gitlab.computing.dcu.ie/collint9/2021-ca400-collint9-coynemt2/src/database/schema"
	"gitlab.computing.dcu.ie/collint9/2021-ca400-collint9-coynemt2/src/server/models"
	"gitlab.computing.dcu.ie/collint9/2021-ca400-collint9-coynemt2/src/server/restapi/operations/files"
)

func FileInfo(params files.GetFileInfoParams, user *schema.User) middleware.Responder {
	var file models.File
	tx := database.DB.Model(&schema.File{}).Where(&schema.File{UserID: user.ID}).Find(&file, params.FileID)
	if tx.Error != nil {
		return files.NewGetFileInfoNotFound()
	}
	return files.NewGetFileInfoOK().WithPayload(&file)
}
