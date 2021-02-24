package data

import (
	c "gitlab.computing.dcu.ie/collint9/2021-ca400-collint9-coynemt2/src/client/config"
	"os"
)

func SetupDirs() error {
	// Create "chunks" directory if it doesn't exist
	dirs := []string{
		c.Config.Synche.Dir,
		c.Config.Synche.DataDir,
	}
	for _, dir := range dirs {
		if err := os.MkdirAll(dir, 0755); err != nil {
			return err
		}
	}
	return nil
}
