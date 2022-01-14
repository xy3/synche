package handlers

import (
	"github.com/go-openapi/runtime/middleware"
	"github.com/xy3/synche/src/server/models"
	"github.com/xy3/synche/src/server/restapi/operations/users"
	"github.com/xy3/synche/src/server/schema"
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
