package config

import (
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"os"
	"path"
)

// TODO Unmarshal the config into a struct

func SetDefaults() error {
	viper.SetDefault("ChunkDir", "../data/chunks")
	viper.SetDefault("ChunkSize", 1) // 1MB
	viper.SetDefault("verbose", false)
	return nil
}

// initConfig reads in config file and ENV variables if set.
func InitConfig(cfgFile string) error {
	err := SetDefaults()
	if err != nil {
		log.Errorf("Failed to set config defaults: %v", err)
		return err
	}

	viper.SetConfigType("yaml")
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// config flag not set, search home dir for config
		home, err := os.UserHomeDir()
		if err != nil {
			log.Fatalf("Could not get $HOME directory: %v", err)
		}

		// Search config in home directory with name ".synche" (without extension).
		viper.AddConfigPath(home)
		viper.SetConfigName(".synche")
		cfgFile = path.Join(home, ".synche.yaml")
		viper.SetConfigFile(cfgFile)
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		log.Debug("Using config file:", viper.ConfigFileUsed())
	} else {
		// the config file does not exist, so create a new one
		err = viper.WriteConfigAs(cfgFile)
		if err != nil {
			log.Fatalf("Unable to create new config file, %v", err)
		}
	}
	return nil
}
