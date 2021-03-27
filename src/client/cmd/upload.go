package cmd

import (
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"gitlab.computing.dcu.ie/collint9/2021-ca400-collint9-coynemt2/src/client/data"
	"gitlab.computing.dcu.ie/collint9/2021-ca400-collint9-coynemt2/src/client/upload"
	"gitlab.computing.dcu.ie/collint9/2021-ca400-collint9-coynemt2/src/files"
	"path"
	"time"
)

func NewUploadCmd(uploader Uploader) *cobra.Command {
	uploadCmd := &cobra.Command{
		Use:   "upload [file path]",
		Short: "Uploads a specified file to the server",
		Long:  `Uploads a specified local file to the server using chunked uploading`,
		Args:  cobra.MinimumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			start := time.Now()
			filePath := args[0]
			err := uploader.Run(filePath)
			if err != nil {
				log.WithError(err).Fatal("Failed to upload the file")
			}
			elapsed := time.Since(start)
			log.Info(elapsed)
		},
	}
	uploadCmd.Flags().StringP("name", "n", "", "store the file on the server with this name instead")
	return uploadCmd
}

//go:generate mockery --name=Uploader --case underscore
type Uploader interface {
	Run(filePath string) error
}

type UploadJob struct {
	chunkUploader upload.ChunkUploader
	fileUploader  upload.FileUploader
	fileHashFunc  files.FileHashFunc
}

func (u UploadJob) Run(filePath string) error {
	file, err := files.AppFS.Open(filePath)
	if err != nil {
		return err
	}
	defer file.Close()
	stat, err := file.Stat()
	if err != nil {
		return err
	}
	hash, err := u.fileHashFunc(filePath)
	if err != nil {
		return err
	}
	splitFile := data.NewSplitFile(stat.Size(), viper.GetInt64("config.chunks.size"), filePath, path.Base(filePath), hash, file)
	return u.fileUploader.AsyncUpload(splitFile)
}

func NewUploadJob(newUploadRequester upload.NewUploadRequester, fileHashFunc files.FileHashFunc) *UploadJob {
	chunkUploader := new(upload.ChunkUpload)
	fileUploader := upload.NewFileUpload(chunkUploader, newUploadRequester)

	return &UploadJob{
		chunkUploader: chunkUploader,
		fileUploader:  fileUploader,
		fileHashFunc:  fileHashFunc,
	}
}

func NewDefaultUploadJob() *UploadJob {
	return NewUploadJob(upload.DefaultNewUploadRequester, files.HashFile)
}

func init() {
	uploadCmd := NewUploadCmd(NewDefaultUploadJob())
	rootCmd.AddCommand(uploadCmd)
	log.SetFormatter(&log.TextFormatter{
		FullTimestamp: true,
	})
	uploadCmd.Flags().Int64P("chunksize", "s", 1024, "size in KB for each chunk")
	uploadCmd.Flags().IntP("workers", "w", 10, "number of chunks to upload in parallel")

	_ = viper.BindPFlag("config.chunks.size", uploadCmd.Flags().Lookup("chunksize"))
	_ = viper.BindPFlag("config.chunks.workers", uploadCmd.Flags().Lookup("workers"))
}
