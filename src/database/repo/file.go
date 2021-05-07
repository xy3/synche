package repo

import (
	"errors"
	"github.com/patrickmn/go-cache"
	log "github.com/sirupsen/logrus"
	"gitlab.computing.dcu.ie/collint9/2021-ca400-collint9-coynemt2/src/database/schema"
	"gitlab.computing.dcu.ie/collint9/2021-ca400-collint9-coynemt2/src/files"
	"gitlab.computing.dcu.ie/collint9/2021-ca400-collint9-coynemt2/src/files/hash"
	"gorm.io/gorm"
	"io"
	"path/filepath"
	"strconv"
	"strings"
)

func GetFileByID(fileID uint, db *gorm.DB) (file *schema.File, err error) {
	strFileID := strconv.Itoa(int(fileID))
	if fileValue, ok := idFileCache.Get(strFileID); ok {
		return fileValue.(*schema.File), nil
	}

	err = db.First(&file, fileID).Error
	if err != nil {
		return nil, err
	}

	idFileCache.Set(strFileID, &file, cache.DefaultExpiration)
	return file, nil
}

// FindFileByFullPath finds a file by its full path on the disk. It assumes the path given
// is a file and not a directory, so if it is given the path to a directory it will treat it
// like a path to a file.
func FindFileByFullPath(path string, db *gorm.DB) (*schema.File, error) {
	if fileValue, ok := pathFileCache.Get(path); ok {
		return fileValue.(*schema.File), nil
	}

	// Get the md5 hash of only the directory part of the path
	path = strings.TrimRight(strings.TrimSpace(path), "/")
	filename := filepath.Base(path)
	dirPathHash := hash.PathHash(filepath.Dir(path))

	file := &schema.File{}
	tx := db.Model(file).Preload("Directory", "path_hash = ?", dirPathHash).Where("name = ?", filename).First(file)

	if tx.Error != nil {
		return nil, tx.Error
	}

	_ = pathFileCache.Add(path, file, cache.DefaultExpiration)
	return file, nil
}

func writeFileData(path string, reader io.Reader) (int64, error) {
	// for ftp it should allow file overwriting, but what about other cases?
	// if exists, _ := files.Afs.Exists(path); exists {
	// 	return 0, errors.New("file already exists")
	// }

	newFile, err := files.Afs.Create(path)
	if err != nil {
		return 0, err
	}

	writtenBytes, err := io.Copy(newFile, reader)
	if err != nil {
		return 0, err
	}

	return writtenBytes, nil
}

func CreateFileFromReader(path string, reader io.Reader, userID uint, db *gorm.DB) (
	file *schema.File,
	err error,
) {
	var (
		writtenBytes int64
		fileHash     string
		parentDir    *schema.Directory
		parentPath   = filepath.Dir(path)
		fileName     = filepath.Base(path)
	)

	if parentDir, err = GetDirByPath(parentPath, db); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			parentDir, err = CreateDirectoryFromPath(parentPath, db)
			if err != nil {
				return nil, err
			}
		} else {
			return nil, err
		}
	}

	// Write the data to the disk
	if writtenBytes, err = writeFileData(path, reader); err != nil {
		return nil, err
	}

	if fileHash, err = hash.File(path); err != nil {
		return nil, err
	}

	file = &schema.File{
		Name:        fileName,
		Size:        writtenBytes,
		Hash:        fileHash,
		DirectoryID: parentDir.ID,
		UserID:      userID,
		Available:   true,
	}

	if err = db.Where(file).FirstOrCreate(file).Error; err != nil {
		return nil, err
	}

	return file, nil
}

func RenameFile(fileID uint, newName string, db *gorm.DB) error {
	var newFile *schema.File
	tx := db.Joins("Directory").Where("id = ?", fileID).First(&newFile)
	if tx.Error != nil {
		return tx.Error
	}

	newPath := filepath.Join(newFile.Directory.Path, newName)
	return MoveFile(newFile, newPath, db)
}

// MoveFile moves the file to a new directory, both in the database and on the disk
func MoveFile(file *schema.File, newFullPath string, db *gorm.DB) (err error) {
	var directory *schema.Directory

	log.Infof("Received move request for: %s to %s", file.Name, newFullPath)

	newDirPath := filepath.Dir(newFullPath)
	newFileName := filepath.Base(newFullPath)

	if len(newDirPath) < 2 {
		return errors.New("invalid directory in path")
	}

	// Find the directory or create it
	if directory, err = GetOrCreateDirectory(newDirPath, db); err != nil {
		return err
	}

	// It is mandatory that this is called before updating the database record
	// as file.Move will use the files current path to move it to the new path
	if err = file.Move(newFullPath, db); err != nil {
		return err
	}

	if newFileName == "" {
		newFileName = file.Name
	}

	tx := db.Model(file).Where("id = ?", file.ID).Updates(schema.File{
		Name:        newFileName,
		DirectoryID: directory.ID,
	})

	file.DirectoryID = directory.ID
	file.Name = newFileName

	return tx.Error
}
