package cmd

import (
	"context"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"gitlab.computing.dcu.ie/collint9/2021-ca400-collint9-coynemt2/src/client/apiclient"
	files2 "gitlab.computing.dcu.ie/collint9/2021-ca400-collint9-coynemt2/src/client/apiclient/files"
	"gitlab.computing.dcu.ie/collint9/2021-ca400-collint9-coynemt2/src/client/apiclient/transfer"
	"gitlab.computing.dcu.ie/collint9/2021-ca400-collint9-coynemt2/src/client/models"
	"gitlab.computing.dcu.ie/collint9/2021-ca400-collint9-coynemt2/src/files"
	"os"
)

var (
	getFileID     uint64
	getFilepath   string
	getFileOutput string
)

// writeFile Creates and writes the file being downloaded to dick
func writeFile(fileInfo *models.File) error {
	var outputFile = fileInfo.Name
	if getFileOutput != "" {
		outputFile = getFileOutput
	}

	fileWriter, err := files.Afs.OpenFile(outputFile, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		return err
	}

	defer fileWriter.Close()

	params := &transfer.DownloadFileParams{
		FileID:  int64(fileInfo.ID),
		Context: context.Background(),
	}

	_, err = apiclient.Client.Transfer.DownloadFile(params, apiclient.ClientAuth, fileWriter)
	if err != nil {
		return err
	}

	log.Infof("Finished downloading the file to: %s", outputFile)
	return nil
}

// getFileByPath Queries the server for the file by file path
func getFileByPath() error {
	fileInfo, err := apiclient.Client.Files.GetFilePathInfo(
		&files2.GetFilePathInfoParams{
			FilePath: getFilepath,
			Context:  context.Background(),
		},
		apiclient.ClientAuth)
	if err != nil {
		return err
	}

	if err := writeFile(fileInfo.Payload); err != nil {
		return err
	}
	return nil
}

// getFileByPath Queries the server for a file by ID
func getFileByID() error {
	fileInfo, err := apiclient.Client.Files.GetFileInfo(
		&files2.GetFileInfoParams{
			FileID:  getFileID,
			Context: context.Background(),
		},
		apiclient.ClientAuth)
	if err != nil {
		return err
	}

	if err := writeFile(fileInfo.Payload); err != nil {
		return err
	}
	return nil
}

// getFileJob Handles the user inputs from the command line and outputs the result of the get command
// retrieves a file from server and writes it to the local client
func getFileJob(cmd *cobra.Command, args []string) error {
	if len(args) > 0 && args[0] != "" {
		getFilepath = args[0]
		if err := getFileByPath(); err != nil {
			return err
		}
	} else if getFilepath != "" {
		if err := getFileByPath(); err != nil {
			return err
		}
	} else if getFileID != 0 {
		if err := getFileByID(); err != nil {
			return err
		}
	} else {
		return cmd.Help()
	}
	return nil
}

// NewGetCmd Handles the user inputs from the command line and requests to download a file from the server
func NewGetCmd() *cobra.Command {
	var getCmd = &cobra.Command{
		Use:     "get",
		Aliases: []string{"download"},
		Short:   "Download a file from a Synche server",
		Long: `Specify a file to download. Files can be specified by path or ID.
  The file will be downloaded to the current directory by default. 
  Examples:
  synche get downloaded_file.jpg
  synche get -i 2 -o downloaded_file.jpg
`,
		PreRun: authenticateUserPreRun,
		Run: func(cmd *cobra.Command, args []string) {
			log.Info("Downloading...")
			if err := getFileJob(cmd, args); err != nil {
				log.WithError(err).Fatal("download failed")
			}
		},
	}

	getCmd.Flags().Uint64VarP(&getFileID, "id", "i", 0, "ID of the file to download")
	getCmd.Flags().StringVarP(&getFilepath, "path", "p", "", "path to the file to download")
	getCmd.Flags().StringVarP(&getFileOutput, "output", "o", "", "download location. either a full file path or directory")
	return getCmd
}

func init() {
	rootCmd.AddCommand(NewGetCmd())
}
