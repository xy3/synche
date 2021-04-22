package schema

import "gorm.io/gorm"

type Upload struct {
	gorm.Model
	ChunkDirectoryID uint
	ChunkDirectory   Directory
	FileID           uint
	File             File
	NumChunks        int64
}
