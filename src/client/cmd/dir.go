package cmd

import (
	"github.com/spf13/cobra"
)

// dirCmd represents the dir command
var dirCmd = &cobra.Command{
	Use:   "dir",
	Short: "A brief description of your command",
	Long: ``,
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		createDirJob(args[0])
	},
}

func init() {
	newCmd.AddCommand(dirCmd)
}
