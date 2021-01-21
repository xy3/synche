package cmd

import (
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(listCmd)
}

var moveCmd = &cobra.Command{
	Use:   "move",
	Short: "<Short description here>",
	Long:  `<Long description here>`,
	Run: func(cmd *cobra.Command, args []string) {
		moveJob()
	},
}

func moveJob() {
	// retrieve a list of files that are on the server and return them here
}