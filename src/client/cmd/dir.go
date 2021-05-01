package cmd

import (
	"github.com/spf13/cobra"
)

// dirCmd represents the dir command
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
