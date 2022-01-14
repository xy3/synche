package client_test

import (
	"github.com/spf13/afero"
	"github.com/stretchr/testify/require"
	"github.com/xy3/synche/src/client"
	"github.com/xy3/synche/src/files"
	"testing"
)

func TestChunkUpload_NewParams(t *testing.T) {
	files.SetFileSystem(afero.NewMemMapFs())

	testCases := []struct {
		Name   string
		Chunk  files.Chunk
		FileID uint64
		ChunkBytes    []byte
		ExpectedError error
	}{
		{
			Chunk:         files.Chunk{Hash: "hash", Num: 1},
			FileID:        1,
			ChunkBytes:    []byte("test file content"),
			ExpectedError: nil,
		},
	}
	for _, tc := range testCases {
		t.Run(
			tc.Name, func(t *testing.T) {
				tc.Chunk.Bytes = &tc.ChunkBytes
				params := client.NewChunkUploadParams(tc.Chunk, tc.FileID)
				if tc.ExpectedError == nil {
					require.Equal(t, tc.FileID, params.FileID)
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
