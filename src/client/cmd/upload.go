package cmd

import (
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	c "gitlab.computing.dcu.ie/collint9/2021-ca400-collint9-coynemt2/src/client/config"
	"gitlab.computing.dcu.ie/collint9/2021-ca400-collint9-coynemt2/src/client/data"
	"gitlab.computing.dcu.ie/collint9/2021-ca400-collint9-coynemt2/src/client/upload"
	"gitlab.computing.dcu.ie/collint9/2021-ca400-collint9-coynemt2/src/files"
	"gitlab.computing.dcu.ie/collint9/2021-ca400-collint9-coynemt2/src/files/hash"
	"path"
	"time"
)

var uploadDirID uint

func NewUploadCmd(fileUploadFunc FileUploadFunc) *cobra.Command {
	uploadCmd := &cobra.Command{
		Use:    "upload [file path]",
		Short:  "Uploads a specified file to the server",
		Long:   `Uploads a specified local file to the server using chunked uploading`,
		Args:   cobra.ExactArgs(1),
		PreRun: authenticateUserPreRun,
		Run: func(cmd *cobra.Command, args []string) {
			start := time.Now()
			filePath := args[0]

			if err := fileUploadFunc(filePath); err != nil {
				log.WithError(err).Fatal("Failed to upload the file")
			}
			elapsed := time.Since(start)
			log.Info(elapsed)
		},
	}

	uploadCmd.Flags().StringP("name", "n", "", "store the file on the server with this name instead")
	uploadCmd.Flags().Int64VarP(&c.Config.Chunks.SizeKB, "chunk-size", "s", 1024, "size in KB for each chunk")
	uploadCmd.Flags().IntVarP(&c.Config.Chunks.Workers, "workers", "w", 10, "number of chunks to upload in parallel")
	uploadCmd.Flags().UintVarP(&uploadDirID, "dir-id", "d", 0, "the ID of the directory to store the file in. default is your home directory on the server")

	return uploadCmd
}

//go:generate mockery --name=FileUploadFunc --case=underscore
type FileUploadFunc func(filePath string) error

func FileUpload(filePath string) error {
	file, err := files.AppFS.Open(filePath)

	if err != nil {
		return err
	}
	defer file.Close()

	stat, err := file.Stat()

	if err != nil {
		return err
	}
	fileHash, err := hash.File(filePath)

	if err != nil {
		return err
	}
	splitFile := data.NewSplitFile(stat.Size(), c.Config.Chunks.SizeKB, filePath, path.Base(filePath), fileHash, file)
	return upload.AsyncUpload(splitFile, uploadDirID, upload.NewUpload, upload.AsyncChunkUpload)
}

func init() {
	uploadCmd := NewUploadCmd(FileUpload)
	rootCmd.AddCommand(uploadCmd)
}
