package cmd

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"gitlab.computing.dcu.ie/collint9/2021-ca400-collint9-coynemt2/src/client/apiclient"
	"gitlab.computing.dcu.ie/collint9/2021-ca400-collint9-coynemt2/src/client/apiclient/files"
	"gitlab.computing.dcu.ie/collint9/2021-ca400-collint9-coynemt2/src/client/apiclient/users"
	"golang.org/x/net/context"
	"strings"
)

var (
	deleteFilepath  string
	deleteFileID    uint64
	deleteUserEmail string
)

// deleteUserJob Deletes an authenticated user specified by the user's email address
func deleteUserJob() error {
	var input string
	log.Infof("Are you sure you wish to delete the user with the following email address? Y/N\n %v", deleteUserEmail)
	if _, err := fmt.Scanf("%s\n", &input); err != nil {
		log.Error(err)
	}

	input = strings.ToLower(input)
	if input == "y" || input == "yes" {
		_, err := apiclient.Client.Users.DeleteUser(users.NewDeleteUserParams().WithEmail(deleteUserEmail), apiclient.ClientAuth)
		if err != nil {
			return err
		}
		log.Info("User account deleted")
	} else if input == "n" || input == "no" {
		log.Info("The user account will not be deleted")
	} else {
		log.Infof("Invalid response: %v", input)
	}
	return nil
}

// deleteJobByPath Sends a request to the server to delete a file specified by an ID
func deleteJobByPath() error {
	_, err := apiclient.Client.Files.DeleteFilepath(
		files.NewDeleteFilepathParams().WithFilePath(
			deleteFilepath).WithContext(
			context.Background()), apiclient.ClientAuth)
	if err != nil {
		return err
	}
	log.Infof("Deleted file with path: %v", deleteFilepath)
	return nil
}

// deleteJobByPath Sends a request to the server to delete a file specified by the path to a file
func deleteJobByID() error {
	_, err := apiclient.Client.Files.DeleteFile(
		files.NewDeleteFileParams().WithFileID(
			deleteFileID).WithContext(
			context.Background()), apiclient.ClientAuth)
	if err != nil {
		return err
	}
	log.Infof("Deleted file with file ID: %v", deleteFileID)
	return nil
}

// deleteCmd Handles the user inputs from the command line and outputs the result of the delete command
var deleteCmd = &cobra.Command{
	Use:     "delete",
	Aliases: []string{"rm", "del"},
	Short:   "Delete a file on the server",
	Long:    `Sends a request to the server to delete file by specified file id or file path.`,
	PreRun:  authenticateUserPreRun,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) > 0 && args[0] != "" {
			deleteFilepath = args[0]
			if err := deleteJobByPath(); err != nil {
				log.WithError(err).Fatal("Failed to delete file by specified path")
			}
		} else if deleteFilepath != "" {
			if err := deleteJobByPath(); err != nil {
				log.WithError(err).Fatal("Failed to delete file by specified path")
			}
		} else if deleteFileID != 0 {
			if err := deleteJobByID(); err != nil {
				log.WithError(err).Fatal("Failed to delete file by specified ID")
			}
		} else if deleteUserEmail != "" {
			if err := deleteUserJob(); err != nil {
				log.WithError(err).Fatal("Failed to delete user")
			}
		} else {
			log.Info(cmd.Help())
		}
	},
}

func init() {
	deleteCmd.Flags().StringVarP(&deleteFilepath, "file-path",
		"p",
		"",
		"Specify the path to a file to delete it")
	deleteCmd.Flags().Uint64VarP(&deleteFileID,
		"file-id",
		"i",
		0,
		"Specify the ID of a file to delete it")
	deleteCmd.Flags().StringVarP(&deleteUserEmail,
		"user-email",
		"u",
		"",
		"Specify the email address of a user account to delete it")
	rootCmd.AddCommand(deleteCmd)
}
