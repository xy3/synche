package handlers

import (
	"github.com/go-openapi/runtime/middleware"
	"gitlab.computing.dcu.ie/collint9/2021-ca400-collint9-coynemt2/src/database"
	"gitlab.computing.dcu.ie/collint9/2021-ca400-collint9-coynemt2/src/database/repo"
	"gitlab.computing.dcu.ie/collint9/2021-ca400-collint9-coynemt2/src/database/schema"
	"gitlab.computing.dcu.ie/collint9/2021-ca400-collint9-coynemt2/src/files/hash"
	"gitlab.computing.dcu.ie/collint9/2021-ca400-collint9-coynemt2/src/server/auth"
	"gitlab.computing.dcu.ie/collint9/2021-ca400-collint9-coynemt2/src/server/models"
	"gitlab.computing.dcu.ie/collint9/2021-ca400-collint9-coynemt2/src/server/restapi/operations/users"
)

func tokenHash(password string) string {
	return hash.MD5HashString(password)
}

func Register(params users.RegisterParams) middleware.Responder {
	user := schema.User{
		Email:     params.Email,
		EmailHash: hash.MD5HashString(params.Email),
		Name:      *params.Name,
		Password:  params.Password,
		TokenHash: tokenHash(params.Password),
	}

	if err := user.ValidateForRegistration(); err != nil {
		return users.NewRegisterDefault(400).WithPayload(models.Error(err.Error()))
	}

	hashedPassword, err := auth.HashPassword(params.Password)
	if err != nil {
		return users.NewRegisterDefault(500).WithPayload("error hashing password")
	}
	user.Password = hashedPassword

	db := database.DB.Begin()

	// Check the database to see if a user already exists with this email
	if db.Find(&schema.User{}, &schema.User{Email: user.Email}).RowsAffected > 0 {
		return users.NewRegisterDefault(409).WithPayload("a user already exists for this email")
	}

	if db.Create(&user).Error != nil {
		db.Rollback()
		return users.NewRegisterDefault(500).WithPayload("error creating the user")
	}

	_, err = repo.SetupUserHomeDir(&user)
	if err != nil {
		db.Rollback()
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
