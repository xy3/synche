package schema

import "gorm.io/gorm"

type Upload struct {
	gorm.Model
	ChunkDirectoryID uint
	ChunkDirectory   ChunkDirectory
	FileID           uint
	File             File
	NumChunks        int64
	UserID           uint
	User             User
}
