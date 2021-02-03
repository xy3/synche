package jobs

import "os"

const (
	UploadDirectory = "../data/received/"
)

func SetupDirs() error {
	err := os.MkdirAll(UploadDirectory, os.ModePerm)
	if err != nil {
		return err
	}
	return nil
}