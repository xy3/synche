package schema

import "gorm.io/gorm"

// Chunk is a file on the server that has a hash and size to compare it
type Chunk struct {
	gorm.Model
	Hash string
	Size int64
}

// FileChunk refers to a chunk that makes up part of a file
type FileChunk struct {
	gorm.Model
	Number      int64
	ChunkID     uint
	Chunk       Chunk
	DirectoryID uint
	Directory   Directory
	FileID      uint
	File        File
	UploadID    uint
	Upload      Upload
}
