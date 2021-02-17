package data_test

import (
	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
	"gitlab.computing.dcu.ie/collint9/2021-ca400-collint9-coynemt2/src/client/data"
	"os"
	"path"
	"testing"
)

const (
	chunkDir         = "testdata"
	defaultChunkName = "test_chunk_"
)

var (
	defaultPath = path.Join(chunkDir, defaultChunkName)
)

func TestDefaultChunkWriter(t *testing.T) {
	data.SetFileSystem(afero.NewMemMapFs())

	testCases := []struct {
		Name       string
		ChunkPath  string
		ChunkBytes []byte
		Expected   error
	}{
		{"writing content to a chunk file", defaultPath + "1", []byte("test content"), nil},
		{"writing empty byte array to chunk file", defaultPath + "2", []byte{}, nil},
	}
	for _, tc := range testCases {
		t.Run(tc.Name, func(t *testing.T) {
			chunk := data.NewChunk(tc.ChunkPath, data.DefaultChunkHashFunc(tc.ChunkBytes), 1)
			actual := data.DefaultChunkWriter(chunk, &tc.ChunkBytes)
			assert.Equal(t, tc.Expected, actual)

			if _, err := data.AppFS.Stat(tc.ChunkPath); os.IsNotExist(err) {
				t.Errorf("chunk file was not created: %v", err)
			}
		})
		// We don't need to remove the files because we are using the afero MemMapFs
	}
}
