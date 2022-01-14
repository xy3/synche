package client

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/go-openapi/runtime"
	httptransport "github.com/go-openapi/runtime/client"
	log "github.com/sirupsen/logrus"
	"github.com/xy3/synche/src/client/models"
	"github.com/xy3/synche/src/client/tokens"
	users2 "github.com/xy3/synche/src/client/users"
	"github.com/xy3/synche/src/files"
	"golang.org/x/crypto/ssh/terminal"
	"os"
	"time"
)

var Authenticator = AuthenticateClient
var ClientAuth runtime.ClientAuthInfoWriter

//go:generate mockery --name=AuthenticatorFunc --case=underscore
type AuthenticatorFunc func(string) error

// AuthenticateClient Authenticates the client using the pre-generated access token
func AuthenticateClient(tokenFile string) error {
	token, err := getToken(tokenFile)
	if err != nil {
		return err
	}

	if err = checkTokenWorks(token.AccessToken); err != nil {
		log.Warn("Current access token is invalid. Please log in again:")
		token, err = TerminalLogin()
		if err != nil {
			return err
		}
		if err = checkTokenWorks(token.AccessToken); err != nil {
			return err
		}
		_ = saveToken(token, tokenFile)
	}

	ClientAuth = httptransport.APIKeyAuth("X-Token", "header", token.AccessToken)
	return nil
}

// Login Logs in user with their email and password
func Login(email, password string) (*models.AccessAndRefreshToken, error) {
	resp, err := Client.Users.Login(&users2.LoginParams{
		Email:    email,
		Password: password,
		Context:  context.Background(),
	})
	if err != nil {
		log.Warn("Login details are invalid")
		return nil, err
	}

	log.Info("Logged in successfully")
	return resp.GetPayload(), nil
}

// TerminalLogin Reads login details from command line and log in user
func TerminalLogin() (*models.AccessAndRefreshToken, error) {
	var email string
	log.Info("Email: ")
	_, _ = fmt.Scanln(&email)

	log.Info("Password: ")
	password, _ := terminal.ReadPassword(0)

	return Login(email, string(password))
}

// checkTokenWorks Validates token
func checkTokenWorks(accessToken string) error {
	tempAuth := httptransport.APIKeyAuth("X-Token", "header", accessToken)
	_, err := Client.Users.Profile(&users2.ProfileParams{Context: context.Background()}, tempAuth)

	return err
}

// getToken Retrieves stored token from token.json and ensure it is valid
func getToken(tokenFile string) (*models.AccessAndRefreshToken, error) {
	saveNeeded := true
	token, err := getSavedToken(tokenFile)
	if err != nil {
		log.Info("No stored credentials found. Please login to authenticate this client:")
		token, err = TerminalLogin()
		if err != nil {
			return nil, err
		}
	} else {
		saveNeeded = false
	}

	if token.AccessTokenExpiry < time.Now().Local().Unix() {
		log.Debug("AccessToken expired, attempting to refresh")
		accessToken, err := refreshAccessToken(token.RefreshToken)
		if err != nil {
			log.Info("Token refresh failed - please log in or run the 'new user' command to create a new account.")
			saveNeeded = true
			token, err = TerminalLogin()
			if err != nil {
				return nil, err
			}
		} else {
			token.AccessToken = accessToken.AccessToken
		}
	}

	if saveNeeded {
		_ = saveToken(token, tokenFile)
	}

	return token, nil
}

// refreshAccessToken Generates refresh token to allow new token to be acquired
func refreshAccessToken(refreshToken string) (*models.AccessToken, error) {
	refreshTokenAuth := httptransport.APIKeyAuth("X-Refresh-Token", "header", refreshToken)
	response, err := Client.Tokens.RefreshToken(&tokens.RefreshTokenParams{Context: context.Background()}, refreshTokenAuth)
	if err != nil {
		return nil, err
	}

	return response.GetPayload(), nil
}

// getSavedToken Retrieves token from token.json
func getSavedToken(tokenFile string) (*models.AccessAndRefreshToken, error) {
	f, err := files.AppFS.Open(tokenFile)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	tok := &models.AccessAndRefreshToken{}
	err = json.NewDecoder(f).Decode(tok)

	return tok, err
}

// saveToken Write generated token to token.json
func saveToken(token *models.AccessAndRefreshToken, tokenFile string) error {
	log.Infof("Saving credentials to: %s", tokenFile)
	f, err := files.AppFS.OpenFile(tokenFile, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		log.WithError(err).Warn("Failed to locally save the token file")
		return err
	}

	defer f.Close()
	_ = json.NewEncoder(f).Encode(token)

	return nil
}
