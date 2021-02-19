package upload_test

import (
	"github.com/spf13/afero"
	"github.com/stretchr/testify/require"
	"gitlab.computing.dcu.ie/collint9/2021-ca400-collint9-coynemt2/src/client/data"
	"gitlab.computing.dcu.ie/collint9/2021-ca400-collint9-coynemt2/src/client/upload"
	"reflect"
	"testing"
)

func TestNewChunkUpload(t *testing.T) {
	t.Run("new chunk upload", func(t *testing.T) {
		actual := upload.NewChunkUpload(data.DefaultChunkHashFunc)
		value := reflect.ValueOf(*actual)
		require.Equal(t, 1, value.NumField())
		require.True(t, value.Field(0).IsValid())
	})
}

func TestChunkUpload_NewParams(t *testing.T) {
	data.SetFileSystem(afero.NewMemMapFs())

	testCases := []struct {
		Name             string
		Chunk            data.Chunk
		UploadRequestID  string
		ChunkFileContent []byte
		ExpectedError    error
	}{
		{
			Chunk: data.Chunk{Path: "testdata/chunk", Hash: "hash", Num: 1},
			UploadRequestID: "reqId",
			ChunkFileContent: []byte("test file content"),
			ExpectedError: nil,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.Name, func(t *testing.T) {
			err := data.Afs.WriteFile(tc.Chunk.Path, tc.ChunkFileContent, 0644)
			if err != nil {
				t.Errorf("Failed to write the test file content: %v", err)
			}

			chunkUpload := new(upload.ChunkUpload)
			params, err := chunkUpload.NewParams(tc.Chunk, tc.UploadRequestID)
			require.Equal(t, tc.ExpectedError, err)
			if tc.ExpectedError == nil {
				require.Equal(t, tc.UploadRequestID, params.UploadRequestID)
				require.Equal(t, int64(tc.Chunk.Num), params.ChunkNumber)
				require.Equal(t, tc.Chunk.Hash, params.ChunkHash)
				paramData := make([]byte, len(tc.ChunkFileContent))
				_, err = params.ChunkData.Read(paramData)
				if err != nil {
					t.Errorf("Could not read the param chunk data: %v", paramData)
				}
				require.Equal(t, tc.ChunkFileContent, paramData)
			}
		})
	}
}

func TestChunkUpload_Upload(t *testing.T) {
	// TODO: Test Chunk Upload function
}
