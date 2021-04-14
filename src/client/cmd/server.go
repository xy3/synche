package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// TODO: Could implement this as a way to manage multiple servers from one client

// serverCmd represents the server command
var serverCmd = &cobra.Command{
	Use:   "server",
	Short: "Add a new server connection to this client",
	Long: `You can have multiple servers connected to this client. 
They can be added using this command and providing the server host
and login credentials.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("server called")
	},
}

func init() {
	newCmd.AddCommand(serverCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// serverCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// serverCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
