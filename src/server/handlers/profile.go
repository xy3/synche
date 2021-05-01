package handlers

import (
	"github.com/go-openapi/runtime/middleware"
	"gitlab.computing.dcu.ie/collint9/2021-ca400-collint9-coynemt2/src/database/schema"
	"gitlab.computing.dcu.ie/collint9/2021-ca400-collint9-coynemt2/src/server/models"
	"gitlab.computing.dcu.ie/collint9/2021-ca400-collint9-coynemt2/src/server/restapi/operations/users"
)

func Profile(_ users.ProfileParams, user *schema.User) middleware.Responder {
	userID := uint64(user.ID)
	return users.NewProfileOK().WithPayload(&models.User{
		ID:            &userID,
		Email:         &user.Email,
		EmailVerified: user.EmailVerified,
		Name:          &user.Name,
		Picture:       &user.Picture,
		Role:          &user.Role,
	})
}
