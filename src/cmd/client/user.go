package main

import (
	"bufio"
	"context"
	"fmt"
	log "github.com/sirupsen/logrus"
	"github.com/xy3/synche/src/client"
	"github.com/xy3/synche/src/models"
	users2 "github.com/xy3/synche/src/client/users"
	"golang.org/x/crypto/ssh/terminal"
	"os"
	"strings"

	"github.com/spf13/cobra"
)

var email, name string

// getUserPassword Reads the user password input via the command line and ensures it's valid
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
		log.Error("Passwords do not match")
		return getUserPassword()
	}

	password = string(passwordInput)
	return
}

// registerUser Prompts the user to input all the details required to create a user
func registerUser() (*models.User, error) {
	scanner := bufio.NewScanner(os.Stdin)

	if email == "" {
		fmt.Println("Email address:")
		if scanner.Scan() {
			email = strings.TrimSpace(scanner.Text())
		}
	}

	if name == "" {
		fmt.Println("Your name:")
		if scanner.Scan() {
			name = strings.TrimSpace(scanner.Text())
		}
	}

	password, err := getUserPassword()
	if err != nil {
		return nil, err
	}

	resp, err := client.Client.Users.Register(&users2.RegisterParams{
		Email:    email,
		Name:     &name,
		Password: password,
		Context:  context.Background(),
	})
	if err != nil {
		return nil, err
	}

	return resp.Payload, nil
}

// userCmd Handles everything to do with the user command
var userCmd = &cobra.Command{
	Use:   "user",
	Short: "Register a new user account",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		_, err := registerUser()
		if err != nil {
			registerError, ok := err.(*users2.RegisterDefault)
			if ok && registerError.Code() == 409 {
				log.Error("A user already exists for that email")
				return
			}
			log.WithError(err).Error("Failed to register a new account")
			return
		}

		log.Info("New account created successfully")
	},
}

func init() {
	newCmd.AddCommand(userCmd)
	userCmd.Flags().StringVarP(&email, "email", "e", "", "User email address")
	userCmd.Flags().StringVarP(&name, "name", "n", "", "Your name")
}
