package repo

import (
	"errors"
	"github.com/patrickmn/go-cache"
	"gitlab.computing.dcu.ie/collint9/2021-ca400-collint9-coynemt2/src/database/schema"
	"gitlab.computing.dcu.ie/collint9/2021-ca400-collint9-coynemt2/src/files/hash"
	"gitlab.computing.dcu.ie/collint9/2021-ca400-collint9-coynemt2/src/server/auth"
	"gorm.io/gorm"
	"strings"
)

// GetUserByEmail retrieves user details relating to an email address
func GetUserByEmail(email string, db *gorm.DB) (*schema.User, error) {
	// The email hash should be used instead of the plaintext email for performance
	emailHash := hash.MD5Hash([]byte(strings.TrimSpace(email)))
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
		EmailHash: hash.MD5HashString(email),
		Password:  password,
		TokenHash: hash.MD5HashString(password),
	}

	if name != nil {
		user.Name = *name
	}
	if picture != nil {
		user.Password = *picture
	}

	if err = user.ValidateForRegistration(); err != nil {
		return nil, err
	}

	hashedPassword, err := auth.HashPassword(password)
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
