package jobs

import (
	"github.com/spf13/viper"
	"os"
)

func SetupDirs() {
	// Create "chunks" directory if it doesn't exist
	_ = os.MkdirAll(viper.GetString("ChunkDir"), os.ModePerm)
}