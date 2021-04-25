package schema

import "gorm.io/gorm"

type ChunkDirectory struct {
	gorm.Model
	Path      string `gorm:"not null"`
	UserID    uint
	User      User
}

type Directory struct {
	gorm.Model
	Name      string `gorm:"not null"`
	Path      string `gorm:"not null"`
	PathHash  string `gorm:"size:32;uniqueIndex"`
	UserID    uint
	User      User
	FileCount uint
}
