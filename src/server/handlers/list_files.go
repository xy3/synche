package handlers

import (
	"github.com/go-openapi/runtime/middleware"
	"gitlab.computing.dcu.ie/collint9/2021-ca400-collint9-coynemt2/src/server/models"
	"gitlab.computing.dcu.ie/collint9/2021-ca400-collint9-coynemt2/src/server/restapi/operations/files"
)

func ListFilesHandler(params files.ListFilesParams) middleware.Responder {
	contents := models.DirectoryContents{
		DirectoryID: "d290f1ee-6c54-4b01-90e6-d701748f0851",
		Listings:    []string{},
	}
	return files.NewListFilesOK().WithPayload(&contents)
}
