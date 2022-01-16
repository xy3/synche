package files_test

import (
	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
	"github.com/xy3/synche/src/files"
	"testing"
)

func TestSetFileSystem(t *testing.T) {
	t.Run("MemMap filesystem", func(t *testing.T) {
		newFs := afero.NewMemMapFs()
		files.SetFileSystem(newFs)
		assert.Equal(t, files.AppFS, newFs)
		assert.Equal(t, files.Afs.Fs, newFs)
	})

	t.Run("OS filesystem", func(t *testing.T) {
		newFs := afero.NewOsFs()
		files.SetFileSystem(newFs)
		assert.Equal(t, files.AppFS, newFs)
		assert.Equal(t, files.Afs.Fs, newFs)
	})
}

func TestDirs(t *testing.T) {
	files.SetFileSystem(afero.NewMemMapFs())
	dirs := []string{"dir1", "dir2", "/nested/dir", "relative/dir"}
	err := files.SetupDirs(files.AppFS, dirs)
	assert.NoError(t, err)
	for _, dir := range dirs {
		isDir, _ := files.Afs.IsDir(dir)
		assert.True(t, isDir)
	}
}

