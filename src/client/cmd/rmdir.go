package cmd

import (
	"context"
	"fmt"
	log "github.com/sirupsen/logrus"
	"gitlab.computing.dcu.ie/collint9/2021-ca400-collint9-coynemt2/src/client/apiclient"
	"gitlab.computing.dcu.ie/collint9/2021-ca400-collint9-coynemt2/src/client/apiclient/files"
	"strings"

	"github.com/spf13/cobra"
)

var rmDirID uint64
var rmDirForce bool

func removeDirByID(id uint64) (deleted bool, err error) {
	if !rmDirForce {
		var confirm string
		fmt.Printf("Are you sure? [Y/n]: ")
		_, _ = fmt.Scanln(&confirm)
		if !strings.HasPrefix(strings.ToLower(confirm), "y") {
			return false, nil
		}
	}

	_, err = apiclient.Client.Files.DeleteDirectory(&files.DeleteDirectoryParams{
		ID:      id,
		Context: context.Background(),
	}, apiclient.ClientAuth)

	if err != nil {
		return false, err
	}

	return true, nil
}

var rmdirCmd = &cobra.Command{
	Use:    "rmdir",
	Short:  "Removes a directory on the server",
	PreRun: authenticateUserPreRun,
	Run: func(cmd *cobra.Command, args []string) {
		if rmDirID == 0 {
			_ = cmd.Help()
			return
		}

		deleted, err := removeDirByID(rmDirID)
		if err != nil {
			log.WithError(err).Error("Failed to remove the directory")
		}

		if deleted {
			log.Info("Directory deleted successfully")
		}
	},
}

func init() {
	rootCmd.AddCommand(rmdirCmd)
	rmdirCmd.Flags().Uint64VarP(&rmDirID, "dir-id", "i", 0, "Specify an ID of a directory to delete. ")
	rmdirCmd.Flags().BoolVarP(&rmDirForce, "force", "f", false, "Skip confirmation and delete the directory immediately")
}
