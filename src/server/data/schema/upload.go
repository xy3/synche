package schema

import "gorm.io/gorm"

type Upload struct {
	gorm.Model
	DirectoryID uint
	Directory   Directory
	FileID      uint
	File        File
	NumChunks   int64
}
