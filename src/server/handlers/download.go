package handlers

import (
	"github.com/go-openapi/runtime/middleware"
	"gitlab.computing.dcu.ie/collint9/2021-ca400-collint9-coynemt2/src/server/data/schema"
	"gitlab.computing.dcu.ie/collint9/2021-ca400-collint9-coynemt2/src/server/restapi/operations/transfer"
)

func DownloadFile(
	params transfer.DownloadFileParams,
	user *schema.User,
) middleware.Responder {
	return middleware.NotImplemented("operation transfer.DownloadFile has not yet been implemented")
}
