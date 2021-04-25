package schema

import "gorm.io/gorm"

type File struct {
	gorm.Model
	Name               string `gorm:"not null"`
	Size               int64 `gorm:"not null"`
	Hash               string `gorm:"index;size:32"`
	ChunkDirectoryID   uint
	ChunkDirectory     ChunkDirectory
	StorageDirectoryID uint
	StorageDirectory   Directory `gorm:"constraint:OnDelete:CASCADE;"`
	UserID             uint
	User               User
	Available          bool
}