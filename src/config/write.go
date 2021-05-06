package config

import (
	"github.com/spf13/viper"
)

// Write Updates the config file
func Write(path string, newConfig interface{}) error {
	viper.Set("config", newConfig)

	if err := viper.WriteConfigAs(path); err != nil {
		return err
	}

	viper.SetConfigFile(path)
	return nil
}
