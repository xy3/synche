package handlers

import (
	"github.com/go-openapi/runtime/middleware"
	"gitlab.computing.dcu.ie/collint9/2021-ca400-collint9-coynemt2/src/database"
	"gitlab.computing.dcu.ie/collint9/2021-ca400-collint9-coynemt2/src/database/repo"
	"gitlab.computing.dcu.ie/collint9/2021-ca400-collint9-coynemt2/src/database/schema"
	"gitlab.computing.dcu.ie/collint9/2021-ca400-collint9-coynemt2/src/server/auth"
	"gitlab.computing.dcu.ie/collint9/2021-ca400-collint9-coynemt2/src/server/models"
	"gitlab.computing.dcu.ie/collint9/2021-ca400-collint9-coynemt2/src/server/restapi/operations/users"
	"gorm.io/gorm"
	"time"
)

func Login(params users.LoginParams, authService auth.Service) middleware.Responder {
	user, err := LoginUser(params.Email, params.Password, database.DB)
	if err != nil {
		return users.NewLoginDefault(401).WithPayload("invalid user credentials")
	}

	// Generate a new access token for the user
	accessToken, err := authService.GenerateAccessToken(user.ID, user.Email, user.Name, user.Picture, user.Role)
	if err != nil {
		return users.NewLoginDefault(500).WithPayload("error signing the token")
	}

	// Generate a new refresh token for the user
	refreshToken, err := authService.GenerateRefreshToken(user.Email, user.TokenHash)
	if err != nil {
		return users.NewLoginDefault(500).WithPayload("failed to generate a refresh token")
	}

	return users.NewLoginOK().WithPayload(&models.AccessAndRefreshToken{
		AccessToken:       accessToken,
		AccessTokenExpiry: time.Now().Local().Add(time.Hour * time.Duration(authService.ExpirationHours)).Unix(),
		RefreshToken:      refreshToken,
	})
}

func LoginUser(email, password string, db *gorm.DB) (*schema.User, error) {
	user, err := repo.GetUserByEmail(email, db)
	if err != nil {
		return nil, err
	}

	err = auth.CheckPassword(user.Password, password)
	if err != nil {
		return nil, err
	}
	return user, nil
}
