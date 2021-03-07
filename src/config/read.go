package config

import (
	"github.com/mitchellh/go-homedir"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"path"
	"path/filepath"
)

func Read(cfgFile, cfgName string, setDefaults func(home string) interface{}) error {
	home, err := homedir.Dir()
	if err != nil {
		return err
	}

	syncheDir := filepath.Join(home, ".synche")
	defaultCfg := setDefaults(home)

	viper.SetConfigName(cfgName)
	viper.SetConfigType("yaml")

	// Set config file locations
	if cfgFile != "" {
		// Use config file from the command line flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// cfgFile not set, scan usual directories for existing config
		viper.AddConfigPath(syncheDir)
		viper.AddConfigPath(home)
		viper.AddConfigPath(".")
	}

	// Enable reading environment variables
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			// Ask the user if they want to customize the default values before proceeding
			newConfig, err := Setup(defaultCfg)
			if err != nil {
				return err
			}
			viper.SetDefault("config", newConfig)

			// Create the config file
			cfgPath := path.Join(syncheDir, cfgName+".yaml")
			err = viper.WriteConfigAs(cfgPath)
			if err != nil {
				return err
			}

			viper.SetConfigFile(cfgPath)
		} else {
			log.Fatalf("Could not read in the config file: %v", err)
		}
	}

	log.Infof("Using config file: %s", viper.ConfigFileUsed())
	return nil
}
