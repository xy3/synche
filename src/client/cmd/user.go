package cmd

import (
	"context"
	"errors"
	"fmt"
	log "github.com/sirupsen/logrus"
	"gitlab.computing.dcu.ie/collint9/2021-ca400-collint9-coynemt2/src/client/apiclient"
	"gitlab.computing.dcu.ie/collint9/2021-ca400-collint9-coynemt2/src/client/apiclient/users"
	"gitlab.computing.dcu.ie/collint9/2021-ca400-collint9-coynemt2/src/client/models"
	"golang.org/x/crypto/ssh/terminal"

	"github.com/spf13/cobra"
)

// userCmd represents the user command
var userCmd = &cobra.Command{
	Use:   "user",
	Short: "Register a new user account",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		_, err := registerUser()
		if err != nil {
			registerError, ok := err.(*users.RegisterDefault)
			if ok && registerError.Code() == 409{
				log.Error("A user already exists for that email")
				return
			}
			log.WithError(err).Error("Failed to register a new account")
			return
		}
		log.Info("New account created successfully")
	},
}

var email, name string

func init() {
	newCmd.AddCommand(userCmd)
	userCmd.Flags().StringVarP(&email, "email", "e", "", "User email address")
	userCmd.Flags().StringVarP(&name, "name", "n", "", "Your name")
}

func registerUser() (*models.User, error) {
	if email == "" {
		fmt.Println("Email address:")
		_, err := fmt.Scanln(&email)
		if err != nil {
			return nil, err
		}
	}
	if name == "" {
		fmt.Println("Your name:")
		_, err := fmt.Scanln(&name)
		if err != nil {
			return nil, err
		}
	}

	password, err := getUserPassword()
	if err != nil {
		return nil, err
	}

	resp, err := apiclient.Client.Users.Register(&users.RegisterParams{
		Email:    email,
		Name:     &name,
		Password: password,
		Context: context.Background(),
	})
	if err != nil {
		return nil, err
	}

	return resp.Payload, nil
}


func getUserPassword() (password string, err error) {
	fmt.Println("Password:")
	passwordInput, err := terminal.ReadPassword(0)
	if err != nil {
		return
	}

	fmt.Println("Confirm Password:")
	confirmPasswordInput, err := terminal.ReadPassword(0)
	if err != nil {
		return
	}

	if string(passwordInput) != string(confirmPasswordInput) {
		err = errors.New("passwords do not match")
		return
	}

	password = string(passwordInput)
	return
}