package handlers

import (
	"github.com/go-openapi/runtime/middleware"
	"gitlab.computing.dcu.ie/collint9/2021-ca400-collint9-coynemt2/src/server/data/schema"
	"gitlab.computing.dcu.ie/collint9/2021-ca400-collint9-coynemt2/src/server/restapi/operations/files"
)

func ListFiles(
	params files.ListParams,
	user *schema.User,
) middleware.Responder {
	return middleware.NotImplemented("operation transfer.ListFiles has not yet been implemented")
}

// TODO: Implement

