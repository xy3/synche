package cmd

import (
	"github.com/fatih/color"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"gitlab.computing.dcu.ie/collint9/2021-ca400-collint9-coynemt2/src/client/apiclient"
	"gitlab.computing.dcu.ie/collint9/2021-ca400-collint9-coynemt2/src/client/apiclient/files"
	"gitlab.computing.dcu.ie/collint9/2021-ca400-collint9-coynemt2/src/client/models"
	"strconv"
)

func printFile(file *models.File) {
	log.Infof("--- File moved successfully. Updated file information: ---")
	log.WithFields(log.Fields{
		"Name":        file.Name,
		"ID":          file.ID,
		"DirectoryID": file.DirectoryID,
		"FileSize":    file.Size}).Infof(color.GreenString(file.Name))
}

func moveJob(params *files.UpdateFileParams) *files.UpdateFileOK {
	resp, err := apiclient.Client.Files.UpdateFile(params, apiclient.ClientAuth)
	if err != nil {
		log.WithError(err).Fatal("Failed to move file")
	}
	return resp
}

func moveToDirectoryID(fileID uint64) *models.File {
	resp := moveJob(files.NewUpdateFileParams().WithFileID(fileID).WithDirectoryID(&moveDirID))
	log.Infof("Directory name: %v\n", &moveDirName)
	return resp.GetPayload()
}

// TODO: replace with a path option (move to path) (this will also double as the rename function)
// func moveToDirectoryName(fileID uint64) *models.File {
// 	resp := moveJob(files.NewUpdateFileParams().WithFileID(fileID).WithDirectoryName(&moveDirName))
// 	log.Infof("Directory name: %v\n", &moveDirName)
// 	return resp.GetPayload()
// }

var moveDirName string
var moveDirID uint64
var moveCmd = &cobra.Command{
	Use:    "move",
	Short:  "move FileID DirectoryID",
	Long:   `Move a file from one specified location to another using the file ID and directory IDs`,
	PreRun: authenticateUserPreRun,
	Run: func(cmd *cobra.Command, args []string) {
		fileID, err := strconv.ParseUint(args[0], 10, 64)
		if err != nil {
			log.Error("Invalid file ID")
		}

		var updatedFile *models.File
		// if moveDirName != "" {
		// 	updatedFile = moveToDirectoryName(fileID)
		// } else {
		// }
		updatedFile = moveToDirectoryID(fileID)

		printFile(updatedFile)
	},
}

func init() {
	// moveCmd.Flags().StringVarP(&moveDirName, "dir-name", "n", "", "Specify the name of a directory to move a file to it")
	moveCmd.Flags().Uint64VarP(&moveDirID, "dir-id", "i", 0, "Specify an ID of a directory to move a file to it")
	rootCmd.AddCommand(moveCmd)
}
