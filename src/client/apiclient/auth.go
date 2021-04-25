package apiclient

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/go-openapi/runtime"
	httptransport "github.com/go-openapi/runtime/client"
	log "github.com/sirupsen/logrus"
	"gitlab.computing.dcu.ie/collint9/2021-ca400-collint9-coynemt2/src/client/apiclient/tokens"
	"gitlab.computing.dcu.ie/collint9/2021-ca400-collint9-coynemt2/src/client/apiclient/users"
	"gitlab.computing.dcu.ie/collint9/2021-ca400-collint9-coynemt2/src/client/models"
	"gitlab.computing.dcu.ie/collint9/2021-ca400-collint9-coynemt2/src/files"
	"golang.org/x/crypto/ssh/terminal"
	"os"
	"time"
)

var Authenticator = AuthenticateClient
var ClientAuth runtime.ClientAuthInfoWriter

//go:generate mockery --name=AuthenticatorFunc --case=underscore
type AuthenticatorFunc func(string) error

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
	log.Debugf("Using access token: %s", token.AccessToken)
	return nil
}

func Login(email, password string) (*models.AccessAndRefreshToken, error) {
	resp, err := Client.Users.Login(&users.LoginParams{
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

func TerminalLogin() (*models.AccessAndRefreshToken, error) {
	log.Info("Email: ")
	var email string
	_, _ = fmt.Scanln(&email)
	log.Info("Password: ")
	password, _ := terminal.ReadPassword(0)
	return Login(email, string(password))
}

func checkTokenWorks(accessToken string) error {
	tempAuth := httptransport.APIKeyAuth("X-Token", "header", accessToken)
	_, err := Client.Users.Profile(&users.ProfileParams{Context: context.Background()}, tempAuth)
	return err
}

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

func refreshAccessToken(refreshToken string) (*models.AccessToken, error) {
	refreshTokenAuth := httptransport.APIKeyAuth("X-Refresh-Token", "header", refreshToken)
	response, err := Client.Tokens.RefreshToken(&tokens.RefreshTokenParams{Context: context.Background()}, refreshTokenAuth)
	if err != nil {
		return nil, err
	}
	return response.GetPayload(), nil
}

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
