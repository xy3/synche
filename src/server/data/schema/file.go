package schema

import "gorm.io/gorm"

type File struct {
	gorm.Model
	Name        string
	Size        int64
	Hash        string
	DirectoryID uint
	Directory   Directory
}