package handlers

import (
	"github.com/go-openapi/runtime/middleware"
	"gitlab.computing.dcu.ie/collint9/2021-ca400-collint9-coynemt2/src/server/models"
	"gitlab.computing.dcu.ie/collint9/2021-ca400-collint9-coynemt2/src/server/restapi/operations/testing"
)

func CheckGetHandler(params testing.CheckGetParams) middleware.Responder {
	return testing.NewCheckGetOK().WithPayload(&models.Message{
		Code:    200,
		Message: "Success",
	})
}
