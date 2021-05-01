package schema

import (
	"errors"
	"gitlab.computing.dcu.ie/collint9/2021-ca400-collint9-coynemt2/src/files"
	"gitlab.computing.dcu.ie/collint9/2021-ca400-collint9-coynemt2/src/files/hash"
	"gorm.io/gorm"
	"io"
	"path/filepath"
)

var (
	ErrInvalidHash = errors.New("hashes do not match")
)

type File struct {
	gorm.Model
	Name        string `gorm:"not null"`
	Size        int64  `gorm:"not null"`
	Hash        string `gorm:"index;size:32"`
	ChunkSize   int64
	DirectoryID uint `gorm:"not null"`
	Directory   *Directory
	UserID      uint
	User        User
	Available   bool

	TotalChunks    int64
	ChunksReceived int64

	Chunks []FileChunk `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;foreignKey:file_id;association_autoupdate:false;association_autocreate:false"`
}

func (f *File) Reader(rootPath *string) (io.ReadSeeker, error) {
	return NewFileReader(f, rootPath)
}

// executeDelete deletes the file data on the disk
func (f *File) executeDelete(filePath string) error {
	if _, err := files.Afs.Stat(filePath); err != nil {
		return nil
	}
	return files.Afs.Remove(filePath)
}

// Delete deletes the file from both the disk and the database.
// It updates the containing directories for size and file count
func (f *File) Delete(db *gorm.DB) (err error) {
	tx := db.Begin()

	if f.Directory == nil {
		if err = tx.Preload("Directory").Find(f).Error; err != nil {
			return err
		}
	}

	filePath := filepath.Join(f.Directory.Path, f.Name)
	if err = f.executeDelete(filePath); err != nil {
		return err
	}

	if _, err = f.Directory.UpdateFileCount(db); err != nil {
		tx.Rollback()
		return err
	}

	if err = db.Unscoped().Delete(f).Error; err != nil {
		tx.Rollback()
		return err
	}

	tx.Commit()
	return nil
}

// Path is used to get the complete path of a file
func (f *File) Path(db *gorm.DB) (string, error) {
	var directory Directory
	if err := db.First(&directory, f.DirectoryID).Error; err != nil {
		return "", err
	}

	return filepath.Join(directory.Path, f.Name), nil
}

func (f *File) MoveToDir(newPath string) (err error) {
	// TODO
	// newDir, err := repo.GetDirByPath(newPath)
	// if err != nil {
	// 	return err
	// }
	//
	// _, err = repo.UpdateFileDirectory(nil, f.ID, newDir.ID)
	return err
}

func (f *File) AppendFromReader(reader io.Reader, num int, rootChunkPath *string) error {
	// TODO
	return nil
}

func (f *File) LastChunkNumber() (int, error) {
	// TODO
	return 0, nil
}

func (f *File) ChunkByNumber(num int) (*Chunk, error) {
	// TODO
	return nil, nil
}

func (f *File) ValidateHash(db *gorm.DB) error {
	path, err := f.Path(db)
	if err != nil {
		return err
	}

	fileHash, err := hash.File(path)
	if err != nil {
		return err
	}

	if fileHash != f.Hash {
		return ErrInvalidHash
	}

	return nil
}

func (f *File) SetAvailable(db *gorm.DB) error {
	f.Available = true
	return db.Save(f).Error
}

func (f *File) SetUnavailable(db *gorm.DB) error {
	f.Available = false
	return db.Save(f).Error
}

func (f *File) Rename(newName string, db *gorm.DB) error {
	f.Name = newName
	return db.Save(f).Error
}