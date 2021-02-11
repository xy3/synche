package cmd

import (
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func NewConfigCommand() *cobra.Command {
	// configCmd represents the config command
	configCmd := &cobra.Command{
		Use:   "config",
		Short: "Displays the current config path",
		Long:  `Displays the current config path that Synche client is currently using`,
		Run: func(cmd *cobra.Command, args []string) {
			log.Infof("Config file path: %s", viper.ConfigFileUsed())
		},
	}
	return configCmd
}

func init() {
	rootCmd.AddCommand(NewConfigCommand())
}
