package main

import (
	"github.com/spf13/cobra"
)

// registerCmd Registers a new user
var registerCmd = &cobra.Command{
	Use:   "register",
	Short: "Create a new user on the server",
	Run:   userCmd.Run,
}

func init() {
	rootCmd.AddCommand(registerCmd)
	registerCmd.Flags().StringVarP(&email, "email", "e", "", "User email address")
	registerCmd.Flags().StringVarP(&name, "name", "n", "", "Your name")
}
