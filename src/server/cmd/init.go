package cmd

import (
	"fmt"
	flag "github.com/spf13/pflag"

	"github.com/spf13/cobra"
)

// initCmd represents the init command
var initCmd = &cobra.Command{
	Use:   "init",
	Short: "A brief description of your command",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("init called")
	},
}

func init() {
	rootCmd.AddCommand(initCmd)

	ServerFlags = flag.CommandLine
	flag.CommandLine = initCmd.Flags()
}
