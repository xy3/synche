package cmd

import (
	"context"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"gitlab.computing.dcu.ie/collint9/2021-ca400-collint9-coynemt2/src/client/apiclient"
	"gitlab.computing.dcu.ie/collint9/2021-ca400-collint9-coynemt2/src/client/apiclient/files"
	"regexp"
)

var (
	newDirPath     string
	newDirParentID uint64
)

// isValidDirName Ensures that the directory name is valid before sending a request to make the directory on the server
// note that " / " is not allowed
func isValidDirName(name string) bool {
	if len(name) < 3 {
		return false
	}
	matched := regexp.MustCompile(`[\\/\?%\*|<>\(\)\[\]\{\}.,:;"]`)
	if matched.FindString(name) == "" {
		return true
	}
	return false
}

// isValidDirParentID Ensures that the parent ID is 0 so that the directory is created in the home folder
func isValidDirParentID(dirID uint64) bool { return dirID != 0 }

// createDirByPath ends a request to create a directory on the server, specified by the path to the directory
func createDirByPath() {
	params := files.NewCreateDirPathParams().WithDirPath(newDirPath).WithContext(context.Background())
	directory, err := apiclient.Client.Files.CreateDirPath(params, apiclient.ClientAuth)

	if err != nil {
		log.WithError(err).Error("failed to create the directory")
		return
	}

	log.Info("Created the directory successfully.")
	log.Infof("Directory ID: %d", directory.Payload.ID)
}

// createDirByID Sends a request to create a directory on the server, the parent directory is specified by ID
// and the directory is specified by a name
func createDirByID(name string, parentID uint64) {
	var parentDirID *uint64
	parentDirID = nil

	if !isValidDirParentID(newDirParentID) {
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

// mkdirCmd Handles the user inputs from the command line and outputs the result of the mkdir command
// creates a directory on the server
var mkdirCmd = &cobra.Command{
	Use:     "mkdir",
	Short:   "Create a new directory",
	Aliases: []string{"md"},
	Long:    `Create a new directory on the server.`,
	PreRun:  authenticateUserPreRun,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) > 0 && newDirParentID == 0 {
			newDirPath = args[0]
			createDirByPath()
		} else if newDirPath != "" {
			createDirByPath()
		} else if len(args) > 0 && newDirParentID != 0 {
			createDirByID(args[0], newDirParentID)
		} else {
			log.Info(cmd.Help())
		}
	},
}

func init() {
	rootCmd.AddCommand(mkdirCmd)
	mkdirCmd.Flags().Uint64VarP(&newDirParentID,
		"parent-dir-id",
		"d",
		0,
		"the id of the directory you want to create a new directory in. Default is the home directory.",
	)
	mkdirCmd.Flags().StringVarP(&newDirPath,
		"dir-path",
		"p",
		"",
		"the path to the directory you want to create",
	)
}
