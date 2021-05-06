package cmd

import (
	"github.com/spf13/cobra"
	flag "github.com/spf13/pflag"
)

// initCmd Represents the init command
var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize or update the server config",
	Long:  `Set up a new Synche server config or update the values in the current config if it exists.`,
	Run: func(cmd *cobra.Command, args []string) {
		// err := config.Update(serverConfig.Config)
		// if err != nil {
		// 	log.WithError(err).Fatal("Failed to initialize the config")
		// }
	},
}

func init() {
	rootCmd.AddCommand(initCmd)
	ServerFlags = flag.CommandLine
	flag.CommandLine = initCmd.Flags()
}
