package scopes

import (
	"gitlab.computing.dcu.ie/collint9/2021-ca400-collint9-coynemt2/src/server/data/schema"
	"gorm.io/gorm"
)

type Scope func(db *gorm.DB) *gorm.DB

func CurrentUser(user *schema.User) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where("user_id = ?", user.ID)
	}
}
