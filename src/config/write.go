package config

import (
	"github.com/spf13/viper"
	"path/filepath"
)

func SetupAndWrite(name, path string, config interface{}) error {
	newConfig, err := Setup(config)
	if err != nil {
		return err
	}

	if path == "" {
		path = filepath.Join(SyncheDir, name+".yaml")
	}

	if err = Write(path, newConfig); err != nil {
		return err
	}
	return nil
}

func Write(path string, newConfig interface{}) error {
	viper.Set("config", newConfig)

	if err := viper.WriteConfigAs(path); err != nil {
		return err
	}

	viper.SetConfigFile(path)
	return nil
}
