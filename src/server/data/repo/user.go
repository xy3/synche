package repo

import (
	"github.com/patrickmn/go-cache"
	"gitlab.computing.dcu.ie/collint9/2021-ca400-collint9-coynemt2/src/server/data"
	"gitlab.computing.dcu.ie/collint9/2021-ca400-collint9-coynemt2/src/server/data/schema"
)

func GetUserByEmail(email string) (*schema.User, error) {
	// Check the cache for the user data
	if res, found := data.Cache.Users.Get(email); found {
		return res.(*schema.User), nil
	}

	// Otherwise, get it from the database
	var user schema.User
	res := data.DB.Where(&schema.User{Email: email}).First(&user)
	if res.Error != nil {
		return nil, res.Error
	}

	data.Cache.Users.Set(email, &user, cache.DefaultExpiration)
	return &user, nil
}
