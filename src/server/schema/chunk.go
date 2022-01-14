package schema

import (
	"github.com/spf13/afero"
	"github.com/xy3/synche/src/client"
	"github.com/xy3/synche/src/files"
	"gorm.io/gorm"
	"path/filepath"
)

// Chunk is a file on the server that has a hash and size to compare it
type Chunk struct {
	gorm.Model
	Hash string `gorm:"uniqueIndex;size:32"`
	Size int64
}

// FileChunk refers to a chunk that makes up part of a file
type FileChunk struct {
	gorm.Model
	Number  int64 `gorm:"uniqueIndex:idx_file_chunk_number"`
	ChunkID uint
	Chunk   Chunk
	FileID  uint `gorm:"uniqueIndex:idx_file_chunk_number"`
	File    File
}

// Reader return a reader with buffer
func (c *Chunk) Reader(rootPath *string) (file afero.File, err error) {
	return files.AppFS.Open(c.Path(rootPath))
}

// Path represent the actual storage path for the chunk
func (c Chunk) Path(rootPath *string) string {
	if rootPath == nil {
		rootPath = &client.Config.Synche.DataDir
	}
	return filepath.Join(*rootPath, c.Hash)
}
