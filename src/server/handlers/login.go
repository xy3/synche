package handlers

import (
	"github.com/go-openapi/runtime/middleware"
	"gitlab.computing.dcu.ie/collint9/2021-ca400-collint9-coynemt2/src/server/auth"
	"gitlab.computing.dcu.ie/collint9/2021-ca400-collint9-coynemt2/src/server/data"
	"gitlab.computing.dcu.ie/collint9/2021-ca400-collint9-coynemt2/src/server/data/schema"
	"gitlab.computing.dcu.ie/collint9/2021-ca400-collint9-coynemt2/src/server/models"
	"gitlab.computing.dcu.ie/collint9/2021-ca400-collint9-coynemt2/src/server/restapi/operations/users"
	"time"
)

func Login(params users.LoginParams, authService auth.Service) middleware.Responder {
	user, err := login(params.Email, params.Password)
	if err != nil {
		return users.NewLoginDefault(401).WithPayload("invalid user credentials")
	}

	// Generate a new access token for the user
	accessToken, err := authService.GenerateAccessToken(user.Email)
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

func login(email, password string) (*schema.User, error) {
	var user schema.User
	result := data.DB.Where(schema.User{Email: email}).First(&user)
	if result.Error != nil {
		return nil, result.Error
	}

	err := auth.CheckPassword(user.Password, password)
	if err != nil {
		return nil, err
	}
	return &user, nil
}
