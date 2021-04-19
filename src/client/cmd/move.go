package cmd

import (
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(moveCmd)
}

var moveCmd = &cobra.Command{
	Use:   "move",
	Short: "Move a file",
	Long:  `<Move a file from one specified location to another using the file ID and directory IDs>`,
	Run: func(cmd *cobra.Command, args []string) {
		moveJob()
	},
}

func moveJob() {
	// retrieve a list of files that are on the server and return them here
}
