package config

import (
	"github.com/mitchellh/go-homedir"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"path/filepath"
)

var (
	HomeDir   string
	SyncheDir string
)

func init() {
	home, err := homedir.Dir()
	if err != nil {
		log.WithError(err).Fatal("Could not retrieve the home directory location")
		return
	}
	HomeDir = home
	SyncheDir = filepath.Join(HomeDir, "synche")
}

// Read Reads a config in using Viper in order to facilitate different configuration file types
func Read(name string, path string) error {
	viper.SetConfigName(name)
	viper.SetConfigType("yaml")

	// Set config file locations
	if path != "" {
		// Use config file from the command line flag.
		viper.SetConfigFile(path)
	} else {
		// cfgFile not set, scan usual directories for existing config
		viper.AddConfigPath(SyncheDir)
		viper.AddConfigPath(HomeDir)
		viper.AddConfigPath(".")
	}

	// Enable reading environment variables
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		return err
	}

	log.Debugf("Using config file: %s", viper.ConfigFileUsed())
	return nil
}

// ReadOrCreate Reads in a config if it exists, if it does not exist it will call Setup() to prompt
// the user to set up their configuration
func ReadOrCreate(name, path string, defaultCfg, configStruct interface{}) (created bool, err error) {
	err = Read(name, path)
	if err != nil {
		log.Warn("No config file found")

		newConfig, err := Setup(defaultCfg)
		if err != nil {
			return false, err
		}

		if path == "" {
			path = filepath.Join(SyncheDir, name+".yaml")
		}

		viper.Set("config", newConfig)
		if err := viper.UnmarshalKey("config", &configStruct); err != nil {
			return false, err
		}

		if err = Write(path, newConfig); err != nil {
			return false, err
		}

		return true, nil
	}
	return false, err
}
