package server

import (
	validator "github.com/wagslane/go-password-validator"
	"golang.org/x/crypto/bcrypt"
)

// minPasswordEntropy is the minimum strength of a password that is accepted
// See https://github.com/wagslane/go-password-validator
const minPasswordEntropy = 50

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}

func CheckPassword(storedPassword, inputPassword string) error {
	err := bcrypt.CompareHashAndPassword([]byte(storedPassword), []byte(inputPassword))
	if err != nil {
		return err
	}
	return nil
}

func ValidateStrongPassword(password string) error {
	if err := validator.Validate(password, minPasswordEntropy); err != nil {
		return err
	}
	return nil
}
