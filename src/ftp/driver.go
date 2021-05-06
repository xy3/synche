package ftp

import (
	"github.com/goftp/server"
	log "github.com/sirupsen/logrus"
	"gitlab.computing.dcu.ie/collint9/2021-ca400-collint9-coynemt2/src/client/data"
	"gitlab.computing.dcu.ie/collint9/2021-ca400-collint9-coynemt2/src/database/repo"
	"gitlab.computing.dcu.ie/collint9/2021-ca400-collint9-coynemt2/src/database/schema"
	"gorm.io/gorm"
	"io"
	"io/ioutil"
	"path/filepath"
)

type Driver struct {
	db      *gorm.DB
	conn    *server.Conn
	user    *schema.User
	homeDir *schema.Directory
	logger  *log.Logger
}

// Init is a hook which is called when a new connection is made
func (d *Driver) Init(conn *server.Conn) {
	d.conn = conn
}

func (d *Driver) buildPath(path string) (string, error) {
	if d.user == nil {
		email := d.conn.LoginUser()
		user, err := repo.GetUserByEmail(email, d.db)
		if err != nil {
			d.logger.Errorf("Could not find a user with the email: %s", email)
			return path, err
		}
		d.user = user
	}

	if d.homeDir == nil {
		dir, err := repo.GetOrCreateHomeDir(d.user, d.db)
		if err != nil {
			d.logger.WithError(err).Error("failed to get the user's home directory")
		}
		d.homeDir = dir
	}
	return filepath.Join(d.homeDir.Path, path), nil
}

func (d *Driver) Stat(path string) (fileInfo server.FileInfo, err error) {
	var (
		fullPath string
		file     *schema.File
		dir      *schema.Directory
	)

	if fullPath, err = d.buildPath(path); err != nil {
		return nil, err
	}

	d.logger.Debugf("Stat-ing path: %s", fullPath)

	if dir, err = repo.GetDirByPath(fullPath, d.db); err == nil {
		return &FileInfo{
			name:     dir.Name,
			size:     4 * data.KB,
			isDir:    true,
			modeTime: dir.UpdatedAt,
		}, nil
	}

	if file, err = repo.FindFileByFullPath(fullPath, d.db); err == nil {
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
func (d *Driver) ChangeDir(path string) error {
	var (
		err      error
		fullPath string
	)

	if fullPath, err = d.buildPath(path); err != nil {
		return err
	}

	_, err = repo.GetDirByPath(fullPath, d.db)
	return err
}

// ListDir is used to list files and subDir of current dir
func (d *Driver) ListDir(path string, callback func(server.FileInfo) error) (err error) {
	var (
		fullPath string
		dir      *schema.Directory
	)

	if fullPath, err = d.buildPath(path); err != nil {
		return err
	}

	dir, err = repo.GetDirWithContentsFromPath(fullPath, d.db)
	if err != nil {
		return err
	}

	d.logger.Debugf("directory has %d children and %d files", len(dir.Children), len(dir.Files))

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
	var (
		fullPath string
		dir      *schema.Directory
	)

	if fullPath, err = d.buildPath(path); err != nil {
		return err
	}

	if dir, err = repo.GetDirByPath(fullPath, d.db); err != nil {
		return
	}
	return dir.Delete(true, d.db)
}

func (d *Driver) DeleteFile(path string) error {
	var file *schema.File

	fullPath, err := d.buildPath(path)
	if err != nil {
		return err
	}

	d.logger.Debugf("Deleting: %s", fullPath)
	if file, err = repo.FindFileByFullPath(fullPath, d.db); err != nil {
		return err
	}
	return file.Delete(d.db)
}

func (d *Driver) Rename(fromPath string, toPath string) (err error) {
	var (
		fullFromPath string
		fullToPath   string
		file         *schema.File
	)
	if fullFromPath, err = d.buildPath(fromPath); err != nil {
		return err
	}
	if fullToPath, err = d.buildPath(toPath); err != nil {
		return err
	}

	if file, err = repo.FindFileByFullPath(fullFromPath, d.db); err != nil {
		return err
	}
	return repo.MoveFile(file, fullToPath, d.db)
}

// PutFile is used to upload file
func (d *Driver) PutFile(path string, connReader io.Reader, append bool) (bytes int64, err error) {
	var (
		fullPath string
		file     *schema.File
	)
	if fullPath, err = d.buildPath(path); err != nil {
		return 0, err
	}

	if append {
		if file, err = repo.FindFileByFullPath(fullPath, d.db); err != nil {
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
		fullPath       string
		fileReadSeeker io.ReadSeeker
		file           *schema.File
	)

	if fullPath, err = d.buildPath(path); err != nil {
		return 0, nil, err
	}

	if file, err = repo.FindFileByFullPath(fullPath, d.db); err != nil {
		return
	}

	if fileReadSeeker, err = file.Reader(d.db); err != nil {
		return
	}

	_, err = fileReadSeeker.Seek(offset, io.SeekStart)
	return file.Size, ioutil.NopCloser(fileReadSeeker), err
}

func (d *Driver) MakeDir(path string) (err error) {
	var fullPath string
	if fullPath, err = d.buildPath(path); err != nil {
		return err
	}

	if _, err = repo.GetDirByPath(fullPath, d.db); err != nil {
		_, err = repo.CreateDirectoryFromPath(fullPath, d.db)
		if err != nil {
			return err
		}
	}

	return err
}
