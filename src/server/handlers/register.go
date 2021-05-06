package handlers

import (
	"github.com/go-openapi/runtime/middleware"
	"gitlab.computing.dcu.ie/collint9/2021-ca400-collint9-coynemt2/src/database"
	"gitlab.computing.dcu.ie/collint9/2021-ca400-collint9-coynemt2/src/database/repo"
	"gitlab.computing.dcu.ie/collint9/2021-ca400-collint9-coynemt2/src/server/models"
	"gitlab.computing.dcu.ie/collint9/2021-ca400-collint9-coynemt2/src/server/restapi/operations/users"
)

func Register(params users.RegisterParams) middleware.Responder {
	db := database.DB.Begin()

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
