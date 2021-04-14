package upload_test

import (
	"github.com/spf13/afero"
	"github.com/stretchr/testify/require"
	"gitlab.computing.dcu.ie/collint9/2021-ca400-collint9-coynemt2/src/client/data"
	"gitlab.computing.dcu.ie/collint9/2021-ca400-collint9-coynemt2/src/client/upload"
	"gitlab.computing.dcu.ie/collint9/2021-ca400-collint9-coynemt2/src/files"
	"testing"
)

func TestChunkUpload_NewParams(t *testing.T) {
	files.SetFileSystem(afero.NewMemMapFs())

	testCases := []struct {
		Name            string
		Chunk           data.Chunk
		UploadRequestID uint64
		ChunkBytes      []byte
		ExpectedError   error
	}{
		{
			Chunk:           data.Chunk{Hash: "hash", Num: 1},
			UploadRequestID: 1,
			ChunkBytes:      []byte("test file content"),
			ExpectedError:   nil,
		},
	}
	for _, tc := range testCases {
		t.Run(
			tc.Name, func(t *testing.T) {
				tc.Chunk.Bytes = &tc.ChunkBytes
				params := upload.NewChunkUploadParams(tc.Chunk, tc.UploadRequestID)
				if tc.ExpectedError == nil {
					require.Equal(t, tc.UploadRequestID, params.UploadID)
					require.Equal(t, tc.Chunk.Num, params.ChunkNumber)
					require.Equal(t, tc.Chunk.Hash, params.ChunkHash)
					paramData := make([]byte, len(tc.ChunkBytes))
					_, err := params.ChunkData.Read(paramData)
					if err != nil {
						t.Errorf("Could not read the param chunk data: %v", paramData)
					}
					require.Equal(t, tc.ChunkBytes, paramData)
				}
			},
		)
	}
}

func TestChunkUpload_Upload(t *testing.T) {
	// TODO: Test Chunk AsyncUpload function
}
