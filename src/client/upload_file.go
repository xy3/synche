package client

import (
	log "github.com/sirupsen/logrus"
	"github.com/xy3/synche/src/client/transfer"
	"github.com/xy3/synche/src/files"
	"sync"
)

//go:generate mockery --name=AsyncFileUploader --case underscore
type AsyncFileUploader func(files.Splitter, NewUploadFunc, AsyncChunkUploader) error

func AsyncUpload(splitter files.Splitter, uploadDirID uint, newUploadFunc NewUploadFunc, asyncChunkUploader AsyncChunkUploader) error {
	// Try to create a new upload file on the server
	uploadFile, err := newUploadFunc(NewUploadParamsFromSplitter(splitter, uploadDirID))
	if err != nil {
		return err
	}

	skipChunks, toSkip, err := getChunksToSkip(uploadFile.ID)
	if err != nil {
		log.WithError(err).Error("Failed to get the number of remaining chunks")
	}

	if int64(toSkip) == splitter.NumChunks() {
		log.Infof("All %d chunks are already on the server", toSkip)
		return nil
	}

	log.WithFields(log.Fields{"workers": Config.Chunks.Workers, "chunksize": Config.Chunks.SizeKB}).Info("Chunk config")
	log.Infof("%#v", uploadFile)

	var wg sync.WaitGroup
	uploadErrors := make(chan error)
	// The closure func here is called everytime a new chunk is read from the file
	err = splitter.Split(
		func(chunk *files.Chunk, index int64) error {
			if skipChunks[chunk.Num] {
				log.Infof("Skipping chunk %d", chunk.Num)
				return nil
			}

			log.WithFields(log.Fields{"chunk": *chunk, "index": index}).Info("")
			if index%int64(Config.Chunks.Workers) == 0 {
				log.Infof("%d - Waiting for %d workers...", index, Config.Chunks.Workers)
				wg.Wait()
			}
			params := NewChunkUploadParams(*chunk, uploadFile.ID)
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


func getChunksToSkip(fileID uint64) (skipChunks map[int64]bool, toSkip int, err error) {
	skipChunks = make(map[int64]bool)
	resp, err := Client.Transfer.CheckUploadedChunks(transfer.NewCheckUploadedChunksParams().WithFileID(fileID), ClientAuth)
	if err != nil {
		return skipChunks, 0, err
	}
	existingChunks := resp.GetPayload()
	for _, n := range existingChunks.ChunkNumbers {
		skipChunks[n] = true
	}
	return skipChunks, len(existingChunks.ChunkNumbers), nil
}