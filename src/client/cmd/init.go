package cmd

import (
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	clientConfig "gitlab.computing.dcu.ie/collint9/2021-ca400-collint9-coynemt2/src/client/config"
	"gitlab.computing.dcu.ie/collint9/2021-ca400-collint9-coynemt2/src/config"
)

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize or update the client config",
	Long:  `Set up a new Synche client config or update the values in the current config if it exists.`,
	Run: func(cmd *cobra.Command, args []string) {
		cfg := config.Config()
		err := cfg.Update(clientConfig.Config)
		if err != nil {
			log.WithError(err).Fatal("Failed to initialize the config")
		}
	},
}

func init() {
	rootCmd.AddCommand(initCmd)
}
