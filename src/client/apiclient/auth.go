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
	ClientAuth = httptransport.APIKeyAuth("X-Token", "header", token.AccessToken)
	log.Infof("Using access token: %s", token.AccessToken)
	return nil
}

func getToken(tokenFile string) (*models.AccessAndRefreshToken, error) {
	token, err := getSavedToken(tokenFile)
	if err != nil {
		token, err = getTokenFromWeb()
		if err != nil {
			return nil, err
		}
	}
	if token.AccessTokenExpiry < time.Now().Local().Unix() {
		log.Debug("AccessToken expired, attempting to refresh")
		accessToken, err := refreshAccessToken(token.RefreshToken)
		if err != nil {
			return token, err
		}
		token.AccessToken = accessToken.AccessToken
	}
	err = saveToken(token, tokenFile)
	if err != nil {
		log.WithError(err).Warn("Failed to locally save the token file")
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

// Request a token from the web, then returns the retrieved token.
func getTokenFromWeb() (*models.AccessAndRefreshToken, error) {
	log.Info("No stored credentials found. Please login to authenticate this client.")
	log.Info("Email: ")
	var email string
	_, _ = fmt.Scanln(&email)
	log.Info("Password: ")
	password, err := terminal.ReadPassword(0)

	resp, err := Client.Users.Login(&users.LoginParams{
		Email:    email,
		Password: string(password),
		Context:  context.Background(),
	})
	if err != nil {
		return nil, err
	}
	log.Info("Logged in successfully")
	return resp.GetPayload(), nil
}

func saveToken(token *models.AccessAndRefreshToken, tokenFile string) error {
	log.Infof("Saving credentials to: %s", tokenFile)
	f, err := files.AppFS.OpenFile(tokenFile, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		return err
	}
	defer f.Close()
	_ = json.NewEncoder(f).Encode(token)
	return nil
}
