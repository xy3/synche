package schema

import (
	"errors"
	"gitlab.computing.dcu.ie/collint9/2021-ca400-collint9-coynemt2/src/files"
	"gitlab.computing.dcu.ie/collint9/2021-ca400-collint9-coynemt2/src/files/hash"
	"gorm.io/gorm"
	"io"
	"os"
	"path/filepath"
)

var (
	ErrInvalidHash = errors.New("hashes do not match")
)

type File struct {
	gorm.Model
	Name        string `gorm:"not null;uniqueIndex:idx_directory_filename;size:255"`
	Size        int64  `gorm:"not null"`
	Hash        string `gorm:"index;size:32;uniqueIndex:idx_user_file_hash"`
	ChunkSize   int64
	DirectoryID uint `gorm:"not null;uniqueIndex:idx_directory_filename"`
	Directory   *Directory
	UserID      uint `gorm:"uniqueIndex:idx_user_file_hash"`
	User        User
	Available   bool

	TotalChunks    int64
	ChunksReceived int64

	Chunks []FileChunk `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;foreignKey:file_id;association_autoupdate:false;association_autocreate:false"`
}

func (f *File) Reader(db *gorm.DB) (io.ReadSeeker, error) {
	path, err := f.Path(db)
	if err != nil {
		return nil, err
	}

	file, err := files.AppFS.Open(path)
	if err != nil {
		return nil, err
	}

	return file, nil
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
	if f.Directory == nil {
		if err = db.Preload("Directory").Find(f).Error; err != nil {
			return err
		}
	}

	filePath := filepath.Join(f.Directory.Path, f.Name)
	if err = f.executeDelete(filePath); err != nil {
		return err
	}

	if _, err = f.Directory.UpdateFileCount(db); err != nil {
		return err
	}

	if err = db.Unscoped().Delete(f).Error; err != nil {
		return err
	}
	return nil
}

// Path is used to get the complete path of a file
func (f *File) Path(db *gorm.DB) (string, error) {
	if f.Directory != nil {
		return filepath.Join(f.Directory.Path, f.Name), nil
	}

	var directory Directory
	if err := db.First(&directory, f.DirectoryID).Error; err != nil {
		return "", err
	}

	return filepath.Join(directory.Path, f.Name), nil
}

// executeMove moves a file
func executeMove(oldPath, newFilePath string) error {
	if exists, _ := files.Afs.Exists(oldPath); !exists {
		return errors.New("a file does not exist at that location")
	}

	if exists, _ := files.Afs.Exists(newFilePath); exists {
		return errors.New("a file already exists at the destination path")
	}

	return files.Afs.Rename(oldPath, newFilePath)
}

// Move only renames the file on the disk, it does not update the database
func (f *File) Move(newPath string, db *gorm.DB) error {
	path, err := f.Path(db)
	if err != nil {
		return err
	}
	return executeMove(path, newPath)
}

func appendToFile(path string, reader io.Reader) (int64, error) {
	if exists, _ := files.Afs.Exists(path); !exists {
		return 0, errors.New("file does not exist")
	}
	file, err := files.AppFS.OpenFile(path, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0644)

	if err != nil {
		return 0, err
	}

	return io.Copy(file, reader)
}

func (f *File) AppendFromReader(reader io.Reader, userID uint, db *gorm.DB) (err error) {
	var (
		size int64
		path string
	)

	if userID != f.UserID {
		return errors.New("you do not have permission to modify this file")
	}

	if path, err = f.Path(db); err != nil {
		return err
	}

	// Write the file data to the end of the existing file
	if size, err = appendToFile(path, reader); err != nil {
		return err
	}

	f.Size += size

	return db.Save(f).Error
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
