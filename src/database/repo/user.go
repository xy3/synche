package repo

import (
	"github.com/patrickmn/go-cache"
	"gitlab.computing.dcu.ie/collint9/2021-ca400-collint9-coynemt2/src/database"
	"gitlab.computing.dcu.ie/collint9/2021-ca400-collint9-coynemt2/src/database/schema"
	"gitlab.computing.dcu.ie/collint9/2021-ca400-collint9-coynemt2/src/files/hash"
	"strings"
)

func GetUserByEmail(email string) (*schema.User, error) {
	// The email hash should be used instead of the plaintext email for performance
	emailHash := hash.MD5Hash([]byte(strings.TrimSpace(email)))
	// Check the cache for the user data
	if res, found := EmailHashUserCache.Get(emailHash); found {
		return res.(*schema.User), nil
	}

	// Otherwise, get it from the database
	var user schema.User
	res := database.DB.Where(&schema.User{EmailHash: emailHash}).First(&user)
	if res.Error != nil {
		return nil, res.Error
	}

	EmailHashUserCache.Set(emailHash, &user, cache.DefaultExpiration)
	return &user, nil
}
