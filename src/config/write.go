package config

import (
	"github.com/spf13/viper"
	"path/filepath"
)

func (cfg SyncheConfig) Create(config interface{}) error {
	newConfig, err := Setup(config)
	if err != nil {
		return err
	}

	if cfg.Path == "" {
		cfg.Path = filepath.Join(cfg.Dir, cfg.Name+".yaml")
	}

	err = cfg.Write(newConfig)
	if err != nil {
		return err
	}
	return nil
}

func (cfg SyncheConfig) Write(newConfig interface{}) error {
	viper.Set("config", newConfig)

	err := viper.WriteConfigAs(cfg.Path)
	if err != nil {
		return err
	}

	viper.SetConfigFile(cfg.Path)
	return nil
}