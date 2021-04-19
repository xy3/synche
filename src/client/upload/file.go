package upload

import (
	log "github.com/sirupsen/logrus"
	c "gitlab.computing.dcu.ie/collint9/2021-ca400-collint9-coynemt2/src/client/config"
	"gitlab.computing.dcu.ie/collint9/2021-ca400-collint9-coynemt2/src/client/data"
	"sync"
)

//go:generate mockery --name=AsyncFileUploader --case underscore
type AsyncFileUploader func(data.Splitter, NewUploadFunc, AsyncChunkUploader) error

func AsyncUpload(splitter data.Splitter, newUploadFunc NewUploadFunc, asyncChunkUploader AsyncChunkUploader) error {
	upload, err := newUploadFunc(NewUploadParamsFromSplitter(splitter))
	if err != nil {
		return err
	}

	log.WithFields(log.Fields{"workers": c.Config.Chunks.Workers, "chunksize": c.Config.Chunks.SizeKB}).Info("Chunk config")
	log.Infof("%#v", upload)

	var wg sync.WaitGroup
	uploadErrors := make(chan error)
	// The closure func here is called everytime a new chunk is read from the file
	err = splitter.Split(
		func(chunk *data.Chunk, index int64) error {
			log.WithFields(log.Fields{"chunk": *chunk, "index": index}).Info("")
			if index%int64(c.Config.Chunks.Workers) == 0 {
				log.Infof("%d - Waiting for %d workers...", index, c.Config.Chunks.Workers)
				wg.Wait()
			}
			params := NewChunkUploadParams(*chunk, upload.ID)
			wg.Add(1)
			go asyncChunkUploader(&wg, params, uploadErrors)
			return nil
		},
	)
	if err != nil {
		return err
	}

	log.Info("Waiting for upload workers to finish.")
	wg.Wait()
	log.Infof("Finished uploading all %d chunks to the server", splitter.NumChunks())
	close(uploadErrors)

	// Here we could attempt to cast the error as an UploadChunkBadRequest or other relevant error types
	for err = range uploadErrors {
		if err != nil {
			log.Error(err)
			return err
		}
	}

	return nil
}
