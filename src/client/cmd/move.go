package cmd

import (
	"github.com/fatih/color"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"gitlab.computing.dcu.ie/collint9/2021-ca400-collint9-coynemt2/src/client/apiclient"
	"gitlab.computing.dcu.ie/collint9/2021-ca400-collint9-coynemt2/src/client/apiclient/files"
	"gitlab.computing.dcu.ie/collint9/2021-ca400-collint9-coynemt2/src/client/models"
)

func printFile(file *models.File) {
	log.Infof("--- File moved successfully. Updated file information: ---")
	log.WithFields(log.Fields{
		"Name":        file.Name,
		"ID":          file.ID,
		"DirectoryID": file.DirectoryID,
		"FileSize":    file.Size}).Infof(color.GreenString(file.Name))
}

func moveIDJob(params *files.UpdateFileParams) *files.UpdateFileOK {
	resp, err := apiclient.Client.Files.UpdateFile(params, apiclient.ClientAuth)
	if err != nil {
		log.WithError(err).Fatal("Failed to move file")
	}
	return resp
}

func movePathJob(params *files.UpdateFilepathParams) *files.UpdateFilepathOK {
	resp, err := apiclient.Client.Files.UpdateFilepath(params, apiclient.ClientAuth)
	if err != nil {
		log.WithError(err).Fatal("Failed to move file")
	}
	return resp
}

func moveFileByIDAndDirID() *models.File {
	resp := moveIDJob(files.NewUpdateFileParams().WithFileID(moveFileID).WithDirectoryID(&moveDirID))
	return resp.GetPayload()
}

func moveFileByIDAndFilepath() *models.File {
	resp := moveIDJob(files.NewUpdateFileParams().WithFileID(moveFileID).WithFilepath(&moveNewFilepath))
	return resp.GetPayload()
}

func moveFileByPathAndDirectoryID() *models.File {
	resp := movePathJob(files.NewUpdateFilepathParams().WithFilepath(moveCurrentFilepath).WithDirectoryID(&moveDirID))
	return resp.GetPayload()
}

func moveFileByPaths() *models.File {
	resp := movePathJob(files.NewUpdateFilepathParams().WithFilepath(moveCurrentFilepath).WithNewFilepath(&moveNewFilepath))
	return resp.GetPayload()
}

var moveCurrentFilepath string
var moveNewFilepath string
var moveFileID uint64
var moveDirID uint64
var moveCmd = &cobra.Command{
	Use:    "move",
	Short:  "Move a file",
	Long:   `Move a file from one specified location to another using the full path to the current location or file ID, and the full path to the new location or the directory ID`,
	PreRun: authenticateUserPreRun,
	Run: func(cmd *cobra.Command, args []string) {
		var file *models.File
		if moveCurrentFilepath != "" && moveNewFilepath != "" {
			file = moveFileByPaths()
		} else if moveCurrentFilepath != "" && moveDirID != 0 {
			file = moveFileByPathAndDirectoryID()
		} else if moveFileID != 0 && moveNewFilepath != "" {
			file = moveFileByIDAndFilepath()
		} else if moveFileID != 0 && moveDirID != 0 {
			file = moveFileByIDAndDirID()
		} else {
			log.Error("Invalid argument supplied")
		}
		printFile(file)
	},
}

func init() {
	moveCmd.Flags().StringVarP(&moveCurrentFilepath, "current-file-path", "c", "", "Specify the name of a directory to move a file to it")
	moveCmd.Flags().StringVarP(&moveNewFilepath, "new-file-path", "p", "", "Specify the name of a directory to move a file to it")
	moveCmd.Flags().Uint64VarP(&moveFileID, "file-id", "f", 0, "Specify the ID of a file to move it")
	moveCmd.Flags().Uint64VarP(&moveDirID, "dir-id", "d", 0, "Specify an ID of a directory to move a file to it")
	rootCmd.AddCommand(moveCmd)
}
