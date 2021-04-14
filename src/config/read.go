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
	SyncheDir = filepath.Join(HomeDir, ".synche")
}

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

	log.Infof("Using config file: %s", viper.ConfigFileUsed())
	return nil
}

func ReadOrCreate(name, path string, defaultCfg interface{}) (created bool, err error) {
	err = Read(name, path)
	if _, ok := err.(viper.ConfigFileNotFoundError); ok {
		log.Warn("No config file found")
		return true, SetupAndWrite(name, path, defaultCfg)
	}
	return false, err
}
