package main

import (
	"github.com/spf13/cobra"
)

// newCmd Flags when something needs to be created. This can be a user, server, or directory
var newCmd = &cobra.Command{
	Use:   "new",
	Short: "Create something new",
	Long:  `The new command is used to register new users, add servers and create directories`,
	Run: func(cmd *cobra.Command, args []string) {
		_ = cmd.Help()
	},
}

func init() {
	rootCmd.AddCommand(newCmd)
}
