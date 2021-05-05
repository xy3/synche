package cmd

import (
	"context"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"gitlab.computing.dcu.ie/collint9/2021-ca400-collint9-coynemt2/src/client/apiclient"
	"gitlab.computing.dcu.ie/collint9/2021-ca400-collint9-coynemt2/src/client/apiclient/files"
	"regexp"
)

var newDirParentID uint64

// mkdirCmd represents the mkdir command
var mkdirCmd = &cobra.Command{
	Use:    "mkdir",
	Short:  "Create a new directory",
	Long:   `Create a new directory on the server. The first argument should the name of the directory`,
	Args:   cobra.ExactArgs(1),
	PreRun: authenticateUserPreRun,
	Run: func(cmd *cobra.Command, args []string) {
		createDirJob(args[0], newDirParentID)
	},
}

func init() {
	rootCmd.AddCommand(mkdirCmd)
	mkdirCmd.Flags().Uint64VarP(
		&newDirParentID,
		"parent-dir-id",
		"p",
		0,
		"the id of the directory you want to create a new directory in. Default is the home directory.",
	)
}

// isValidDirName doesn't allow any special characters in the name
// note that " / " is not allowed
func isValidDirName(name string) bool {
	matched := regexp.MustCompile(`[\\/\?%\*|<>\(\)\[\]\{\}.,:;"]`)
	if matched.FindString(name) == "" {
		return true
	}
	return false
}

func createDirJob(name string, parentID uint64) {
	if len(name) < 3 {
		log.Error("directory name must be more than 3 characters long")
		return
	}

	var parentDirID *uint64
	parentDirID = nil

	if newDirParentID == 0 {
		log.Error("parent id not specified")
		return
	}

	if parentID != 0 {
		parentDirID = &parentID

		if !isValidDirName(name) {
			log.Error("directory name is invalid")
			return
		}

		directory, err := apiclient.Client.Files.CreateDirectory(&files.CreateDirectoryParams{
			DirectoryName:     name,
			ParentDirectoryID: parentDirID,
			Context:           context.Background(),
		}, apiclient.ClientAuth)

		if err != nil {
			log.WithError(err).Error("failed to create the directory")
			return
		}

		log.Info("Created the directory successfully.")
		log.Infof("Directory ID: %d", directory.Payload.ID)
	}
}
