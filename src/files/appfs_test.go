package files

import (
	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestSetFileSystem(t *testing.T) {
	t.Run("MemMap filesystem", func(t *testing.T) {
		newFs := afero.NewMemMapFs()
		SetFileSystem(newFs)
		assert.Equal(t, AppFS, newFs)
		assert.Equal(t, Afs.Fs, newFs)
	})

	t.Run("OS filesystem", func(t *testing.T) {
		newFs := afero.NewOsFs()
		SetFileSystem(newFs)
		assert.Equal(t, AppFS, newFs)
		assert.Equal(t, Afs.Fs, newFs)
	})
}
