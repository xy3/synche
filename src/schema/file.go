package schema

import (
	"errors"
	"github.com/go-openapi/strfmt"
	"github.com/xy3/synche/src/files"
	"github.com/xy3/synche/src/models"
	"gorm.io/gorm"
	"io"
	"math"
	"os"
	"path/filepath"
)

var (
	ErrInvalidHash = errors.New("hashes do not match")
)

// swagger:model File
type File struct {
	Model
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

// Validate validates this file
func (m *File) Validate(formats strfmt.Registry) error {
	return nil
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

func (f *File) ValidateHash(db *gorm.DB) error {
	path, err := f.Path(db)
	if err != nil {
		return err
	}

	fileHash, err := files.FileHash(path)
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

func (f *File) ConvertToFileModel() *models.File {
	return &models.File{
		ID:             uint64(f.ID),
		Size:           f.Size,
		Hash:           f.Hash,
		Name:           f.Name,
		Available:      f.Available,
		TotalChunks:    f.TotalChunks,
		ChunksReceived: f.ChunksReceived,
		DirectoryID:    uint64(f.DirectoryID),
	}
}

func (f *File) LastChunkNumber() (int, error) {
	// TODO
	return 0, nil
}

func (f *File) ChunkByNumber(num int) (*Chunk, error) {
	// TODO
	return nil, nil
}

var (
	// ErrInvalidFile represent a invalid file
	ErrInvalidFile = errors.New("invalid file")
	// ErrFileNoChunks represent that a file has no any chunks
	ErrFileNoChunks = errors.New("file has no any chunks")
	// ErrInvalidSeekWhence represent invalid seek whence
	ErrInvalidSeekWhence = errors.New("invalid seek whence")
	// ErrNegativePosition represent negative position
	ErrNegativePosition = errors.New("negative read position")
)

// The code below can be used when we implement reassembling on the fly, and store only chunks.
// For now, files do not need to be read using this code - it is obsolete. The can simply be read
// using the file system, since they are already reassembled by read-stage.

type fileReader struct {
	file               *File
	rootPath           *string
	currentChunkReader io.ReadSeekCloser
	totalChunkNumber   int
	currentChunkNumber int
	alreadyReadCount   int
}

func NewFileReader(file *File, rootPath *string) (*fileReader, error) {
	if file == nil {
		return nil, ErrInvalidFile
	}

	var (
		err              error
		firstChunk       *Chunk
		chunkReader      io.ReadSeekCloser
		totalChunkNumber int
	)

	if totalChunkNumber, err = file.LastChunkNumber(); err != nil {
		return nil, err
	}

	if totalChunkNumber == 0 {
		return nil, ErrFileNoChunks
	}

	if firstChunk, err = file.ChunkByNumber(1); err != nil {
		return nil, err
	}

	if chunkReader, err = firstChunk.Reader(*rootPath); err != nil {
		return nil, err
	}

	return &fileReader{
		file:               file,
		currentChunkReader: chunkReader,
		rootPath:           rootPath,
		currentChunkNumber: 1,
		totalChunkNumber:   totalChunkNumber,
	}, nil
}

func (fr *fileReader) Read(p []byte) (readCount int, err error) {
	if fr.alreadyReadCount >= int(fr.file.Size) {
		_ = fr.currentChunkReader
		return 0, io.EOF
	}
	defer func() { fr.alreadyReadCount += readCount }()
	readCount, err = fr.currentChunkReader.Read(p)
	if err != nil && err == io.EOF {
		_ = fr.currentChunkReader.Close()
		fr.currentChunkNumber++
		var nextChunk *Chunk
		if nextChunk, err = fr.file.ChunkByNumber(fr.currentChunkNumber); err != nil {
			return
		}
		if fr.currentChunkReader, err = nextChunk.Reader(*fr.rootPath); err != nil {
			return readCount, err
		}
		return readCount, nil
	}
	return readCount, err
}

func (fr *fileReader) Seek(offset int64, whence int) (abs int64, err error) {
	switch whence {
	case io.SeekStart:
		abs = offset
	case io.SeekCurrent:
		abs = int64(fr.alreadyReadCount) + offset
	case io.SeekEnd:
		abs = fr.file.Size + offset
	default:
		return 0, ErrInvalidSeekWhence
	}
	if abs < 0 {
		return 0, ErrNegativePosition
	}
	if abs >= fr.file.Size {
		fr.alreadyReadCount = int(abs)
		fr.currentChunkNumber = fr.totalChunkNumber
		return abs, nil
	}
	var (
		currentChunk       *Chunk
		currentChunkReader io.ReadSeekCloser
		currentChunkNumber = int(math.Ceil(float64(abs) / float64(fr.file.ChunkSize)))
	)

	if abs%fr.file.ChunkSize == 0 {
		currentChunkNumber++
	}

	if currentChunkNumber == fr.currentChunkNumber {
		currentChunkReader = fr.currentChunkReader
	} else {
		if currentChunk, err = fr.file.ChunkByNumber(currentChunkNumber); err != nil {
			return 0, nil
		}
		if currentChunkReader, err = currentChunk.Reader(*fr.rootPath); err != nil {
			return 0, err
		}
	}
	if _, err = currentChunkReader.Seek(abs%fr.file.ChunkSize, io.SeekStart); err != nil {
		return 0, err
	}
	if currentChunkNumber != fr.currentChunkNumber {
		_ = fr.currentChunkReader.Close()
	}
	fr.currentChunkReader = currentChunkReader
	fr.currentChunkNumber = currentChunkNumber
	fr.alreadyReadCount = int(abs)
	return abs, nil
}
