package repo

import (
	"errors"
	"github.com/patrickmn/go-cache"
	"github.com/xy3/synche/src/files"
	schema2 "github.com/xy3/synche/src/schema"
	"gorm.io/gorm"
	"io"
	"path/filepath"
	"strconv"
	"strings"
)

// GetFileByID returns the file information relating to the file with the specified ID
func GetFileByID(fileID uint, db *gorm.DB) (file *schema2.File, err error) {
	strFileID := strconv.Itoa(int(fileID))
	if fileValue, ok := idFileCache.Get(strFileID); ok {
		return fileValue.(*schema2.File), nil
	}

	err = db.First(&file, fileID).Error
	if err != nil {
		return nil, err
	}

	idFileCache.Set(strFileID, file, cache.DefaultExpiration)
	return file, nil
}

// FindFileByFullPath finds a file by its full path on the disk. It assumes the path given
// is a file and not a directory, so if it is given the path to a directory it will treat it
// like a path to a file.
func FindFileByFullPath(path string, db *gorm.DB) (*schema2.File, error) {
	if fileValue, ok := pathFileCache.Get(path); ok {
		return fileValue.(*schema2.File), nil
	}

	// Get the md5 hash of only the directory part of the path
	path = strings.TrimRight(strings.TrimSpace(path), "/")
	filename := filepath.Base(path)
	dirPathHash := files.PathHash(filepath.Dir(path))

	file := &schema2.File{}
	tx := db.Model(file).Preload("Directory", "path_hash = ?", dirPathHash).Where("name = ?", filename).First(file)

	if tx.Error != nil {
		return nil, tx.Error
	}

	_ = pathFileCache.Add(path, file, cache.DefaultExpiration)
	return file, nil
}

func writeFileData(path string, reader io.Reader) (int64, error) {
	// for ftp it should allow file overwriting, but what about other cases?
	// e.g.
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
	file *schema2.File,
	err error,
) {
	var (
		writtenBytes int64
		fileHash     string
		parentDir    *schema2.Directory
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

	if fileHash, err = files.FileHash(path); err != nil {
		return nil, err
	}

	file = &schema2.File{
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

func RenameFile(fileID uint, newName string, db *gorm.DB) (newFile *schema2.File, err error) {
	tx := db.Model(&schema2.File{}).Preload("Directory").Where("id = ?", fileID).First(&newFile)
	if tx.Error != nil {
		return nil, tx.Error
	}

	newPath := filepath.Join(newFile.Directory.Path, newName)
	return newFile, MoveFile(newFile, newPath, db)
}

// MoveFile moves the file to a new directory, both in the database and on the disk
func MoveFile(file *schema2.File, newFullPath string, db *gorm.DB) (err error) {
	var directory *schema2.Directory

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

	tx := db.Model(file).Where("id = ?", file.ID).Updates(schema2.File{
		Name:        newFileName,
		DirectoryID: directory.ID,
	})

	file.DirectoryID = directory.ID
	file.Name = newFileName

	if _, err := directory.UpdateFileCount(db); err != nil {
		return err
	}
	return tx.Error
}

func GetTotalFileChunks(fileID uint64, db *gorm.DB) (uint64, error) {
	var file schema2.File
	if tx := db.First(&file, fileID); tx.Error != nil {
		return 0, tx.Error
	}
	return uint64(file.TotalChunks), nil
}
