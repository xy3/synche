package handlers

import (
	"github.com/go-openapi/runtime/middleware"
	"gitlab.computing.dcu.ie/collint9/2021-ca400-collint9-coynemt2/src/server/data/schema"
	"gitlab.computing.dcu.ie/collint9/2021-ca400-collint9-coynemt2/src/server/restapi/operations/files"
)

func DeleteFile(
	params files.DeleteFileParams,
	user *schema.User,
) middleware.Responder {
	return middleware.NotImplemented("operation files.DeleteFile has not yet been implemented")
}
