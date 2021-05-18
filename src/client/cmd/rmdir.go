package cmd

import (
	"context"
	"fmt"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"gitlab.computing.dcu.ie/collint9/2021-ca400-collint9-coynemt2/src/client/apiclient"
	"gitlab.computing.dcu.ie/collint9/2021-ca400-collint9-coynemt2/src/client/apiclient/files"
	"strings"
)

var (
	rmDirID    uint64
	rmDirPath  string
	rmDirForce bool
)

// logResponse Logs whether the directory was deleted or not based on the server's response
func logResponse(deleted bool, err error) {
	if err != nil {
		log.WithError(err).Error("Failed to remove the directory")
	}

	if deleted {
		log.Info("Directory deleted successfully")
	}
}

// getConfirmation ensures the user wishes to delete a directory if they did
// not specify the -f flag
func getConfirmation() (delete bool) {
	var confirm string
	fmt.Printf("Are you sure? [Y/n]: ")
	_, _ = fmt.Scanln(&confirm)
	return strings.HasPrefix(strings.ToLower(confirm), "y")
}

// removeDirByID Removes a directory that is specified by the path to the directory
func removeDirByPath() (deleted bool, err error) {
	if !rmDirForce {
		if !getConfirmation() {
			return false, err
		}
	}

	params := files.NewDeleteDirPathParams().WithDirPath(rmDirPath).WithContext(context.Background())
	_, err = apiclient.Client.Files.DeleteDirPath(params, apiclient.ClientAuth)

	if err != nil {
		return false, err
	}

	return true, nil
}

// removeDirByID Removes a directory that is specified by an ID
func removeDirByID() (deleted bool, err error) {
	if !rmDirForce {
		if !getConfirmation() {
			return false, err
		}
	}

	_, err = apiclient.Client.Files.DeleteDirectory(&files.DeleteDirectoryParams{
		ID:      rmDirID,
		Context: context.Background(),
	}, apiclient.ClientAuth)

	if err != nil {
		return false, err
	}

	return true, nil
}

// rmdirCmd Handles the user inputs from the command line and outputs the result of the remove directory command
var rmdirCmd = &cobra.Command{
	Use:    "rmdir",
	Short:  "Removes a directory on the server",
	PreRun: authenticateUserPreRun,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) > 0 && args[0] != "" {
			rmDirPath = args[0]
			deleted, err := removeDirByPath()
			logResponse(deleted, err)
		} else if rmDirPath != "" {
			deleted, err := removeDirByPath()
			logResponse(deleted, err)
		} else if rmDirID != 0 {
			deleted, err := removeDirByID()
			logResponse(deleted, err)
		}
	},
}

func init() {
	rootCmd.AddCommand(rmdirCmd)
	rmdirCmd.Flags().Uint64VarP(&rmDirID, "dir-id", "d", 0, "Specify an ID of a directory to delete")
	rmdirCmd.Flags().StringVarP(&rmDirPath, "dir-path", "p", "", "Specify the path to a directory to delete")
	rmdirCmd.Flags().BoolVarP(&rmDirForce, "force", "f", false, "Skip confirmation and delete the directory immediately")
}
