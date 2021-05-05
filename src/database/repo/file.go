package repo

import (
	"errors"
	"github.com/patrickmn/go-cache"
	log "github.com/sirupsen/logrus"
	"gitlab.computing.dcu.ie/collint9/2021-ca400-collint9-coynemt2/src/database"
	"gitlab.computing.dcu.ie/collint9/2021-ca400-collint9-coynemt2/src/database/schema"
	"gitlab.computing.dcu.ie/collint9/2021-ca400-collint9-coynemt2/src/files"
	f "gitlab.computing.dcu.ie/collint9/2021-ca400-collint9-coynemt2/src/files"
	"gitlab.computing.dcu.ie/collint9/2021-ca400-collint9-coynemt2/src/files/hash"
	"gitlab.computing.dcu.ie/collint9/2021-ca400-collint9-coynemt2/src/server/models"
	"gorm.io/gorm"
	"io"
	"path/filepath"
	"strings"
)

func ConvertFileToModelsFile(file schema.File) *models.File {
	return &models.File{
		ID:             uint64(file.ID),
		ChunksReceived: file.ChunksReceived,
		DirectoryID:    uint64(file.DirectoryID),
		Hash:           file.Hash,
		Name:           file.Name,
		Size:           file.Size,
		TotalChunks:    file.TotalChunks,
		Available:      file.Available,
	}
}

func CreateFileFromReader(path string, reader io.Reader, userID uint, db *gorm.DB) (
	file *schema.File,
	err error,
) {
	var (
		writtenBytes int64
		fileHash     string
		parentDir    *schema.Directory
		fileName     = filepath.Base(path)
	)

	log.Infof("trying to create file: %s", path)

	if parentDir, err = GetDirByPath(filepath.Dir(path)); err != nil {
		return nil, err
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

func GetFileByID(fileID uint) (*schema.File, error) {
	var file schema.File
	if err := database.DB.First(&file, fileID).Error; err != nil {
		return nil, err
	}
	return &file, nil
}

func GetFileByPath(user *schema.User, path string) (*schema.File, error) {
	var file *schema.File
	homeDir, err := GetHomeDir(user.ID)
	if err != nil {
		return file, err
	}

	fileDir, err := GetDirByPath(filepath.Clean(filepath.Join(homeDir.Path, filepath.Dir(path))))
	if err != nil {
		return file, err
	}

	file, err = GetFileByNameAndDirID(filepath.Base(path), uint64(fileDir.ID))
	if err != nil {
		return file, err
	}

	return file, nil
}

func GetFileByNameAndDirID(filename string, dirID uint64) (*schema.File, error) {
	var file *schema.File

	tx := database.DB.Where(&schema.File{Name: filename, DirectoryID: uint(dirID)}).First(&file)
	if tx.Error != nil {
		return nil, tx.Error
	}
	return file, nil
}

// FindFileByFullPath finds a file by its full path on the disk. It assumes the path given
// is a file and not a directory, so if it is given the path to a directory it will treat it
// like a path to a file.
func FindFileByFullPath(path string) (*schema.File, error) {
	if fileValue, ok := pathFileCache.Get(path); ok {
		return fileValue.(*schema.File), nil
	}

	// Get the md5 hash of only the directory part of the path
	filename := filepath.Base(path)
	dirPath := filepath.Dir(path)
	dirPathHash := hash.PathHash(dirPath)
	log.Infof("Received FindFileByFullPath request for file: %s", path)
	log.Infof("Received FindFileByFullPath request with dir as: %s", dirPath)

	file := &schema.File{}
	tx := database.DB.Model(schema.File{}).Joins("Directory").Where(schema.File{
		Name:      filename,
		Directory: &schema.Directory{PathHash: dirPathHash},
	}).First(file)

	if tx.Error != nil {
		return nil, tx.Error
	}

	log.Info(file)
	_ = pathFileCache.Add(path, file, cache.DefaultExpiration)
	return file, nil
}

func MoveFileToNewDir(file *schema.File, newDirID uint) error {
	var directory schema.Directory
	tx := database.DB.Where("id = ?", newDirID).Find(&directory)
	if tx.Error != nil {
		return tx.Error
	}

	filePath, err := file.Path(database.DB)
	if err != nil {
		return err
	}
	newPath := filepath.Join(directory.Path, file.Name)

	if err = f.Afs.Rename(filePath, newPath); err != nil {
		return err
	}
	return nil
}

func MoveFile(file *schema.File, newFullPath string) (err error) {
	var (
		db        *gorm.DB
		directory *schema.Directory
	)

	log.Infof("Received move request for: %s to %s", file.Name, newFullPath)

	db = database.DB
	newDirPath := filepath.Dir(newFullPath)
	newFileName := filepath.Base(newFullPath)

	log.Infof("%#v", newDirPath)
	log.Infof("%#v", newFileName)

	if len(newDirPath) < 2 {
		return errors.New("invalid directory in path")
	}

	// Find the directory or create it
	directory, err = GetDirByPath(newDirPath)
	if err != nil {
		directory, err = CreateDirectoryFromPath(newDirPath, db)
		if err != nil {
			return err
		}
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

	return tx.Error
}

func RenameFile(fileID uint, newName string) error {
	// TODO: This needs to update the file on the disk, perhaps just use the MoveFile method?
	tx := database.DB.Model(schema.File{}).Where("id = ?", fileID).Update("name", newName)
	return tx.Error
}

func UpdateFileByDirectory(file *schema.File, dir *schema.Directory) (*schema.File, error) {
	if err := MoveFileToNewDir(file, dir.ID); err != nil {
		return nil, err
	}
	file.DirectoryID = dir.ID
	tx := database.DB.Model(&schema.File{}).Where("id = ?", file.ID).Update("directory_id", dir.ID)
	if tx.Error != nil {
		return nil, tx.Error
	}
	return file, nil
}

func UpdateFilename(file *schema.File, newName string) (*schema.File, error) {
	if strings.ContainsAny(newName, "/\\") {
		return nil, errors.New("file names cannot contain slashes")
	}

	directory, err := GetDirectoryByID(file.DirectoryID)
	if err != nil {
		return nil, err
	}

	//filePath, err := file.Path(database.DB)
	//if err != nil {
	//	return nil, err
	//}

	oldPath := filepath.Join(directory.Path, file.Name)
	newPath := filepath.Join(directory.Path, newName)
	if err = f.Afs.Rename(oldPath, newPath); err != nil {
		log.WithError(err).Info("write err")
		return nil, err
	}

	file.Name = newName
	if err = RenameFile(file.ID, newName); err != nil {
		return nil, err
	}

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
