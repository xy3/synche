package files

import (
	"github.com/spf13/afero"
)

func SetupDirs(fs afero.Fs, dirs []string) error {
	for _, dir := range dirs {
		if err := fs.MkdirAll(dir, 0755); err != nil {
			return err
		}
	}
	return nil
}
