package upload

import (
	"github.com/go-openapi/runtime"
	log "github.com/sirupsen/logrus"
	"gitlab.computing.dcu.ie/collint9/2021-ca400-collint9-coynemt2/src/client/apiclient"
	"gitlab.computing.dcu.ie/collint9/2021-ca400-collint9-coynemt2/src/client/apiclient/files"
	"gitlab.computing.dcu.ie/collint9/2021-ca400-collint9-coynemt2/src/client/data"
	"sync"
)

//go:generate mockery --name=ChunkUploader --case underscore
type ChunkUploader interface {
	UploadChunks(chunks []data.Chunk, uploadRequestID string) error
	Upload(wg *sync.WaitGroup, params *files.UploadChunkParams, uploadErrors chan error)
	NewParams(chunk data.Chunk, requestID string) (*files.UploadChunkParams, error)
}

type ChunkUpload struct {
	hashFunc data.ChunkHashFunc
}

func NewChunkUpload(hashFunc data.ChunkHashFunc) *ChunkUpload {
	return &ChunkUpload{hashFunc: hashFunc}
}

func (cu *ChunkUpload) UploadChunks(chunks []data.Chunk, uploadRequestID string) error {
	var wg sync.WaitGroup
	uploadErrors := make(chan error, len(chunks))

	for _, chunk := range chunks {
		wg.Add(1)
		params, err := cu.NewParams(chunk, uploadRequestID)
		if err != nil {
			panic(err)
		}
		go cu.Upload(&wg, params, uploadErrors)
	}

	// TODO: return when an error is encountered (here we could wait for N errors and then 'call it a day'?)
	err := <-uploadErrors
	if err != nil {
		wg.Wait()
		close(uploadErrors)
		return err
	}

	//// Attempt to cast the error as an UploadChunkBadRequest
	//if err, ok := err.(*files.UploadChunkBadRequest); ok {
	//	log.Errorf("%v", *err.Payload.Message)
	//} else {
	//	// otherwise, just print the Error() as supplied by go-swagger
	//	log.Errorf("%s", err.Error())
	//}

	log.Infof("Waiting for upload workers to finish.")
	wg.Wait()
	log.Infof("All upload workers completed.")
	close(uploadErrors)

	//for err := range uploadErrors {
	//	if err != nil {
	//		log.Errorf("%v\n", err)
	//		return err
	//	}
	//}

	return nil
}

func (cu *ChunkUpload) Upload(wg *sync.WaitGroup, params *files.UploadChunkParams, uploadErrors chan error) {
	defer wg.Done()

	// TODO: Have a limit of errors before we consider it "not working"?
	resp, err := apiclient.Client.Files.UploadChunk(params)
	if err != nil {
		// TODO: Bug - the channel seems to be closing prematurely when an upload fails... see [ch213]
		uploadErrors <- err
		return
	}

	log.Infof("Response received: %#v", resp.Payload)

	// TODO: Do something here with the response payload to check if the chunk was uploaded correctly
}

func (cu *ChunkUpload) NewParams(chunk data.Chunk, requestID string) (*files.UploadChunkParams, error) {
	file, err := data.AppFS.Open(chunk.Path)
	if err != nil {
		return nil, err
	}

	readCloser := runtime.NamedReader(chunk.Hash, file)
	params := files.NewUploadChunkParams().
		WithChunkNumber(int64(chunk.Num)). // TODO: Change the API endpoint to use uint64 instead of int64
		WithChunkHash(chunk.Hash).
		WithUploadRequestID(requestID).
		WithChunkData(readCloser)

	return params, nil
}
