package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// TODO: Could implement this as a way to manage multiple servers from one client

// serverCmd Represents the server command
var serverCmd = &cobra.Command{
	Use:   "server",
	Short: "Add a new server connection to this client",
	Long: `You can have multiple servers connected to this client. 
They can be added using this command and providing the server host
and login credentials.`,
	PreRun: authenticateUserPreRun,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("server called")
	},
}

func init() {
	newCmd.AddCommand(serverCmd)
}
