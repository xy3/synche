package handlers

import (
	"github.com/go-openapi/runtime/middleware"
	"github.com/xy3/synche/src/schema"
	"github.com/xy3/synche/src/server"
	"github.com/xy3/synche/src/server/models"
	"github.com/xy3/synche/src/server/repo"
	"github.com/xy3/synche/src/server/restapi/operations/users"
	"gorm.io/gorm"
	"time"
)

// Login Calls Login() to log in a user, and generates access tokens
func Login(params users.LoginParams, authService server.Service) middleware.Responder {
	user, err := LoginUser(params.Email, params.Password, server.DB)
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

	err = server.CheckPassword(user.Password, password)
	if err != nil {
		return nil, err
	}
	return user, nil
}
