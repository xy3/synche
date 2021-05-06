package cmd

import (
	"github.com/spf13/cobra"
)

// mkdirCmd Handles the user inputs from the command line and outputs the result of the dir command
// creates a directory on the server
var dirCmd = &cobra.Command{
	Use:    "dir",
	Short:  "A brief description of your command",
	Long:   ``,
	Args:   cobra.ExactArgs(1),
	PreRun: authenticateUserPreRun,
	Run: func(cmd *cobra.Command, args []string) {
		createDirJob(args[0], newDirParentID)
	},
}

// TODO: fix bug that doesn't allow file paths to be the default args
func init() {
	newCmd.AddCommand(dirCmd)
	dirCmd.Flags().Uint64VarP(
		&newDirParentID,
		"parent-dir-id",
		"p",
		0,
		"the id of the directory you want to create a new directory in. Default is the home directory.",
	)
}
