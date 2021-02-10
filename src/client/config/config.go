package config

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"os"
	"path"
)

// TODO Unmarshal the config into a struct

func SetDefaults() {
	viper.SetDefault("ChunkDir", "../data/chunks")
	viper.SetDefault("ChunkSize", 1) // 1MB
	viper.SetDefault("verbose", false)
}

// initConfig reads in config file and ENV variables if set.
func InitConfig(cfgFile string) {
	SetDefaults()

	viper.SetConfigType("yaml")
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// config flag not set, search home dir for config
		home, err := os.UserHomeDir()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		// Search config in home directory with name ".synche" (without extension).
		viper.AddConfigPath(home)
		viper.SetConfigName(".synche")
		cfgFile = path.Join(home, ".synche.yaml")
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		log.Trace("Using config file:", viper.ConfigFileUsed())
	} else {
		// the config file does not exist, so create a new one
		err = viper.WriteConfigAs(cfgFile)
		if err != nil {
			fmt.Printf("Unable to create new config file, %v", err)
			os.Exit(1)
		}
	}
}
