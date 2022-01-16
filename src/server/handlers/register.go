package handlers

import (
	"github.com/go-openapi/runtime/middleware"
	"github.com/xy3/synche/src/server"
	"github.com/xy3/synche/src/server/models"
	"github.com/xy3/synche/src/server/repo"
	"github.com/xy3/synche/src/server/restapi/operations/users"
)

// Register Handles new user registration and responds to the client
func Register(params users.RegisterParams) middleware.Responder {
	db := server.DB.Begin()

	user, err := repo.NewUser(params.Email, params.Password, params.Name, params.Picture, db)
	if err != nil {
		return users.NewRegisterDefault(500).WithPayload(models.Error("failed to register the user: " + err.Error()))
	}

	if _, err = repo.SetupUserHomeDir(user, db); err != nil {
		return users.NewRegisterDefault(500).WithPayload("could not create the user's home directory")
	}

	db.Commit()

	userId := uint64(user.ID)
	return users.NewRegisterOK().WithPayload(&models.User{
		Email:         &user.Email,
		EmailVerified: user.EmailVerified,
		ID:            &userId,
		Name:          &user.Name,
		Picture:       &user.Picture,
		Role:          &user.Role,
	})
}
