package schema

import "gorm.io/gorm"

type File struct {
	gorm.Model
	Name        string
	Size        int64
	Hash        string
	DirectoryID uint
	Directory   Directory
	UserID      uint
	User        User
}

func NewFile(name string, size int64, hash string, directoryID, userID uint) *File {
	return &File{Name: name, Size: size, Hash: hash, DirectoryID: directoryID, UserID: userID}
}
