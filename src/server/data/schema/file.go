package schema

import "gorm.io/gorm"

type File struct {
	gorm.Model
	Name               string
	Size               int64
	Hash               string
	ChunkDirectoryID   uint
	ChunkDirectory     Directory
	StorageDirectoryID uint
	StorageDirectory   Directory
	UserID             uint
	User               User
}

func NewFile(name string, size int64, hash string, chunkDirectoryID uint, storageDirectoryID uint, userID uint) *File {
	return &File{Name: name, Size: size, Hash: hash, ChunkDirectoryID: chunkDirectoryID, StorageDirectoryID: storageDirectoryID, UserID: userID}
}
