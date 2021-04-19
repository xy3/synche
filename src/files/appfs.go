package files

import "github.com/spf13/afero"

var AppFS = afero.NewOsFs()

// Afs provides utils on-top of the AppFS, e.g. Afs.ReadFile()
var Afs = &afero.Afero{Fs: AppFS}

func SetFileSystem(fs afero.Fs) {
	AppFS = fs
	Afs.Fs = fs
}
