package ftp

import (
	"github.com/goftp/server"
	log "github.com/sirupsen/logrus"
	"gitlab.computing.dcu.ie/collint9/2021-ca400-collint9-coynemt2/src/client/data"
	"gitlab.computing.dcu.ie/collint9/2021-ca400-collint9-coynemt2/src/database/repo"
	"gitlab.computing.dcu.ie/collint9/2021-ca400-collint9-coynemt2/src/database/schema"
	"gitlab.computing.dcu.ie/collint9/2021-ca400-collint9-coynemt2/src/files"
	"gorm.io/gorm"
	"io"
	"io/ioutil"
	"path/filepath"
)

type Driver struct {
	db      *gorm.DB
	conn    *server.Conn
	user    *schema.User
	rootDir *schema.Directory
	logger  *log.Logger
}

// Init is a hook which is called when a new connection is made
func (d *Driver) Init(conn *server.Conn) {
	d.conn = conn
}

func (d *Driver) buildPath(path string) string {
	if d.user == nil {
		email := d.conn.LoginUser()
		user, err := repo.GetUserByEmail(email)
		if err != nil {
			d.logger.Errorf("Could not find a user with the email: %s", email)
			d.conn.Close()
			return ""
		}
		d.user = user
	}

	if d.rootDir == nil {
		dir, err := repo.GetHomeDir(d.user.ID)
		if err != nil {
			if _, err = repo.SetupUserHomeDir(d.user); err != nil {
				d.logger.WithError(err).Error("failed to setup the user's home directory")
				return path
			}
		}
		if exists, _ := files.Afs.IsDir(dir.Path); !exists {
			if _, err = repo.CreateUserHomeDir(d.user); err != nil {
				d.logger.WithError(err).Error("failed to create the user's home directory")
				return path
			}
		}
		d.rootDir = dir
	}
	return filepath.Join(d.rootDir.Path, path)
}

func (d *Driver) Stat(path string) (fileInfo server.FileInfo, err error) {
	fullPath := d.buildPath(path)
	d.logger.Infof("statting path: %s", fullPath)

	var (
		file *schema.File
		dir  *schema.Directory
	)

	if dir, err = repo.GetDirByPath(fullPath); err == nil {
		return &FileInfo{
			name:     dir.Name,
			size:     4 * data.KB,
			isDir:    true,
			modeTime: dir.UpdatedAt,
		}, nil
	}

	if file, err = repo.FindFileByFullPath(fullPath); err == nil {
		return &FileInfo{
			name:     file.Name,
			size:     file.Size,
			isDir:    false,
			modeTime: file.UpdatedAt,
		}, nil
	}

	return nil, err
}

// ChangeDir is used to change the current directory, if the directory doesn't exist, it will be created.
func (d *Driver) ChangeDir(path string) (err error) {
	_, err = repo.GetDirByPath(d.buildPath(path))
	return err
}

// ListDir is used to list files and subDir of current dir
func (d *Driver) ListDir(path string, callback func(server.FileInfo) error) (err error) {
	fullPath := d.buildPath(path)
	d.logger.Infof("listing dir: %s", fullPath)

	var dir *schema.Directory
	dir, err = repo.GetDirWithContentsFromPath(fullPath, d.db)
	if err != nil {
		return err
	}

	d.logger.Infof("directory has %d children and %d files", len(dir.Children), len(dir.Files))

	for _, child := range dir.Children {
		if err = callback(&FileInfo{
			name:     child.Name,
			size:     4 * data.KB,
			isDir:    true,
			modeTime: child.UpdatedAt,
		}); err != nil {
			return
		}
	}

	for _, file := range dir.Files {
		if !file.Available {
			continue
		}
		if err = callback(&FileInfo{
			name:     file.Name,
			size:     file.Size,
			isDir:    false,
			modeTime: file.UpdatedAt,
		}); err != nil {
			return
		}
	}
	return
}

func (d *Driver) DeleteDir(path string) (err error) {
	var dir *schema.Directory
	if dir, err = repo.GetDirByPath(d.buildPath(path)); err != nil {
		return
	}
	return dir.Delete(true, d.db)
}

func (d *Driver) DeleteFile(path string) (err error) {
	fullPath := d.buildPath(path)
	var file *schema.File
	d.logger.Infof("Delete request received for: %s", fullPath)
	if file, err = repo.FindFileByFullPath(fullPath); err != nil {
		return
	}
	return file.Delete(d.db)
}

func (d *Driver) Rename(fromPath string, toPath string) error {
	file, err := repo.FindFileByFullPath(d.buildPath(fromPath))
	if err != nil {
		return err
	}
	return repo.MoveFile(file, d.buildPath(toPath))
}

// PutFile is used to upload file
func (d *Driver) PutFile(path string, connReader io.Reader, append bool) (bytes int64, err error) {
	fullPath := d.buildPath(path)

	var file *schema.File
	if append {
		if file, err = repo.FindFileByFullPath(fullPath); err != nil {
			return
		}
		originSize := file.Size

		if err = file.AppendFromReader(connReader, d.user.ID, d.db); err != nil {
			return
		}
		return file.Size - originSize, nil
	}

	if file, err = repo.CreateFileFromReader(fullPath, connReader, d.user.ID, d.db); err != nil {
		return
	}
	return file.Size, nil
}

// GetFile is used to download a file
func (d *Driver) GetFile(path string, offset int64) (size int64, rc io.ReadCloser, err error) {
	var (
		fileReadSeeker io.ReadSeeker
		file           *schema.File
	)
	if file, err = repo.FindFileByFullPath(d.buildPath(path)); err != nil {
		return
	}

	if fileReadSeeker, err = file.Reader(d.db); err != nil {
		return
	}

	_, err = fileReadSeeker.Seek(offset, io.SeekStart)
	return file.Size, ioutil.NopCloser(fileReadSeeker), err
}

func (d *Driver) MakeDir(path string) error {
	path = d.buildPath(path)
	dir, err := repo.GetDirByPath(path)
	if err != nil {
		dir, err = repo.CreateDirectoryFromPath(path, d.db)
		if err != nil {
			return err
		}
	}
	d.logger.Infof("created: %v", dir)
	return err
}
