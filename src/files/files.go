package files

import "github.com/spf13/afero"

var AppFS = afero.NewOsFs()

// Afs provides utils on-top of the AppFS, e.g. Afs.ReadFile()
var Afs = &afero.Afero{Fs: AppFS}

func SetFileSystem(fs afero.Fs) {
	AppFS = fs
	Afs.Fs = fs
}

func SetupDirs(fs afero.Fs, dirs []string) error {
	for _, dir := range dirs {
		if err := fs.MkdirAll(dir, 0755); err != nil {
			return err
		}
	}
	return nil
}
