package schema

import (
	"gorm.io/gorm"
)

// User uses an email hash to have a set key length for finding a user by their email address
// swagger:model User
type User struct {
	Model
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


