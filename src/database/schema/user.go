package schema

import (
	"errors"
	"fmt"
	"gitlab.computing.dcu.ie/collint9/2021-ca400-collint9-coynemt2/src/server/auth"
	"gorm.io/gorm"
	"net/mail"
)

// User uses an email hash to have a set key length for finding a user by their email address
type User struct {
	gorm.Model
	Email         string `gorm:"not null"`
	EmailVerified bool
	EmailHash     string `gorm:"not null;uniqueIndex;size:32"`
	Password      string `gorm:"not null"`
	Name          string
	Picture       string
	TokenHash     string `gorm:"not null;size:32"`
	Role          string `gorm:"default:user"`
}

func (user *User) Delete(db *gorm.DB) error {
	if err := db.Unscoped().Delete(user).Error; err != nil {
		return err
	}
	return nil
}

func (user *User) ValidateForRegistration() error {
	if !isEmailValid(user.Email) {
		return errors.New("email is invalid")
	}
	if len(user.Name) < 3 {
		return errors.New("name is too short")
	}
	if err := auth.ValidateStrongPassword(user.Password); err != nil {
		return err
	}
	if len(user.TokenHash) != 32 {
		return fmt.Errorf("token hash must be 32 characters long but it is %d", len(user.TokenHash))
	}
	return nil
}

func isEmailValid(email string) bool {
	_, err := mail.ParseAddress(email)
	return err == nil
}
