package handlers

import (
	"github.com/go-openapi/runtime/middleware"
	"gitlab.computing.dcu.ie/collint9/2021-ca400-collint9-coynemt2/src/server/data"
	"gitlab.computing.dcu.ie/collint9/2021-ca400-collint9-coynemt2/src/server/data/schema"
	"gitlab.computing.dcu.ie/collint9/2021-ca400-collint9-coynemt2/src/server/data/scopes"
	"gitlab.computing.dcu.ie/collint9/2021-ca400-collint9-coynemt2/src/server/models"
	"gitlab.computing.dcu.ie/collint9/2021-ca400-collint9-coynemt2/src/server/restapi/operations/files"
)

func FileInfo(
	params files.GetFileInfoParams,
	user *schema.User,
) middleware.Responder {
	var file models.File
	tx := data.DB.Model(&schema.File{}).Scopes(scopes.CurrentUser(user)).Find(&file, params.FileID)
	if tx.Error != nil {
		return files.NewGetFileInfoNotFound()
	}
	return files.NewGetFileInfoOK().WithPayload(&file)
}
