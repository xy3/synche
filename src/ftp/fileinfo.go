package ftp

import (
	"os"
	"time"
)

// FileInfo is used to represent the information of file or directory
type FileInfo struct {
	name     string
	size     int64
	isDir    bool
	modeTime time.Time
}

// Sys always return nil
func (f *FileInfo) Sys() interface{} {
	return nil
}

// IsDir represent whether the object is a directory
func (f *FileInfo) IsDir() bool {
	return f.isDir
}

// ModTime is used to return the modify time of file
func (f *FileInfo) ModTime() time.Time {
	return f.modeTime
}

// Mode returns a file's mode and permission bits.
func (f *FileInfo) Mode() os.FileMode {
	if f.isDir {
		return os.ModePerm | os.ModeDir
	}
	return os.ModePerm
}

// Size return the size of file or directory
func (f *FileInfo) Size() int64 {
	return f.size
}

// Name return the name of file or directory
func (f *FileInfo) Name() string {
	return f.name
}

// Owner return the owner of file
func (f *FileInfo) Owner() string {
	return "synche"
}

// Group return the group of file
func (f *FileInfo) Group() string {
	return "synche"
}
