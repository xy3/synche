package config

import (
	"github.com/mitchellh/go-homedir"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"path/filepath"
)

// type SyncheConfigManager interface {
// 	Read() error
// 	ReadOrCreate(config interface{}) error
// 	Create(config interface{}) error
// 	Write(config interface{}) error
// }

var config SyncheConfig

func Config() SyncheConfig {
	return config
}

type SyncheConfig struct {
	Home  string
	Dir   string
	Path  string
	Name  string
	IsNew bool
}

func New(path string, name string) (*SyncheConfig, error) {
	home, err := homedir.Dir()
	if err != nil {
		return nil, err
	}
	syncheDir := filepath.Join(home, ".synche")
	config = SyncheConfig{Home: home, Dir: syncheDir, Path: path, Name: name, IsNew: false}
	return &config, nil
}

func (cfg SyncheConfig) Read() error {
	viper.SetConfigName(cfg.Name)
	viper.SetConfigType("yaml")

	// Set config file locations
	if cfg.Path != "" {
		// Use config file from the command line flag.
		viper.SetConfigFile(cfg.Path)
	} else {
		// cfgFile not set, scan usual directories for existing config
		viper.AddConfigPath(cfg.Dir)
		viper.AddConfigPath(cfg.Home)
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

func (cfg *SyncheConfig) ReadOrCreate(defaultCfg interface{}) error {
	err := cfg.Read()
	if _, ok := err.(viper.ConfigFileNotFoundError); ok {
		log.Warn("No config file found")
		err = cfg.Create(defaultCfg)
		if err != nil {
			return err
		}
		cfg.IsNew = true
	} else {
		return err
	}
	return nil
}
