package repo

import (
	"errors"
	"fmt"
	"github.com/patrickmn/go-cache"
	"github.com/xy3/synche/src/files"
	"github.com/xy3/synche/src/schema"
	"github.com/xy3/synche/src/server"
	"gorm.io/gorm"
	"net/mail"
	"strings"
)

// GetUserByEmail finds a user either in the cache or database using their email address
func GetUserByEmail(email string, db *gorm.DB) (*schema.User, error) {
	// The email hash should be used instead of the plaintext email for performance
	emailHash := files.MD5Hash([]byte(strings.TrimSpace(email)))
	// Check the cache for the user data
	if res, found := EmailHashUserCache.Get(emailHash); found {
		return res.(*schema.User), nil
	}

	// Otherwise, get it from the database
	var user schema.User
	res := db.Where(&schema.User{EmailHash: emailHash}).First(&user)
	if res.Error != nil {
		return nil, res.Error
	}

	EmailHashUserCache.Set(emailHash, &user, cache.DefaultExpiration)
	return &user, nil
}

// NewUser validates a user and enters their details into the database
func NewUser(email, password string, name, picture *string, db *gorm.DB) (user *schema.User, err error) {
	user = &schema.User{
		Email:     email,
		EmailHash: files.MD5HashString(email),
		Password:  password,
		TokenHash: files.MD5HashString(password),
	}

	if name != nil {
		user.Name = *name
	}
	if picture != nil {
		user.Picture = *picture
	}

	if err = validateForRegistration(*user, server.ValidateStrongPassword); err != nil {
		return nil, err
	}

	hashedPassword, err := server.HashPassword(password)
	if err != nil {
		return nil, err
	}

	user.Password = hashedPassword

	// Check the database to see if a user already exists with this email
	if db.Where(&schema.User{Email: user.Email}).Find(&schema.User{}).RowsAffected > 0 {
		return nil, errors.New("user already exists")
	}

	tx := db.Create(user)
	if tx.Error != nil {
		return nil, tx.Error
	}

	return user, nil
}

func validateForRegistration(user schema.User, passwordValidator func(string) error) error {
	if !isEmailValid(user.Email) {
		return errors.New("email is invalid")
	}
	if len(user.Name) < 3 {
		return errors.New("name is too short")
	}
	if err := passwordValidator(user.Password); err != nil {
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