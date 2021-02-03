package upload

import (
	"encoding/hex"
	"fmt"
	"github.com/go-openapi/runtime"
	"github.com/kalafut/imohash"
	"gitlab.computing.dcu.ie/collint9/2021-ca400-collint9-coynemt2/src/client/apiclient"
	"gitlab.computing.dcu.ie/collint9/2021-ca400-collint9-coynemt2/src/client/apiclient/files"
	"os"
	"sync"
)

func Chunk(wg *sync.WaitGroup, params *files.UploadChunkParams) {
	defer wg.Done()

	resp, err := apiclient.Default.Files.UploadChunk(params)
	if err != nil {
		panic(err) // TODO
	}

	fmt.Printf("%#v\n", resp.Payload)
}

func NewChunkUploadParams(chunkPath, chunkName, uploadRequestId string, chunkNum int64) (*files.UploadChunkParams, error) {
	chunkHash, err := imohash.SumFile(chunkPath)
	if err != nil {
		return nil, err
	}

	file, err := os.Open(chunkPath)
	if err != nil {
		return nil, err
	}

	hash := hex.EncodeToString(chunkHash[:])
	readCloser := runtime.NamedReader(chunkName, file)
	params := files.NewUploadChunkParams().
		WithChunkNumber(chunkNum).
		WithChunkHash(hash).
		WithUploadRequestID(uploadRequestId).
		WithChunkData(readCloser)

	return params, nil
}
