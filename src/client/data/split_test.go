package data_test

import (
	log "github.com/sirupsen/logrus"
	"github.com/spf13/afero"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"gitlab.computing.dcu.ie/collint9/2021-ca400-collint9-coynemt2/src/client/data"
	"gitlab.computing.dcu.ie/collint9/2021-ca400-collint9-coynemt2/src/client/data/mocks"
	"io/ioutil"
	"reflect"
	"testing"
)

var (
	mockChunkWriter = new(mocks.ChunkWriter)
)

func TestSplitJob_Split(t *testing.T) {
	mockChunkWriter.On("Execute", mock.Anything, mock.Anything).Return(nil)
	log.SetOutput(ioutil.Discard)
	data.SetFileSystem(afero.NewOsFs())

	err := data.AppFS.MkdirAll("testdata/chunks", 0755)
	if err != nil {
		t.Errorf("Failed to create the chunks directory: %v", err)
	}

	testCases := []struct {
		Name           string
		Filepath       string
		ChunkMBS       uint64
		ExpectedChunks int
		ExpectedError  error
	}{
		{"splitting into smaller chunks", "testdata/split_me", 1, 3, nil},
		{"splitting into 1 same size chunk", "testdata/split_me", 10, 1, nil},
	}
	for _, tc := range testCases {
		t.Run(tc.Name, func(t *testing.T) {
			splitter := data.NewSplitJob(mockChunkWriter.Execute, data.DefaultChunkHashFunc, "testdata/chunks", tc.ChunkMBS)

			file, err := data.Afs.Open(tc.Filepath)
			if err != nil {
				t.Fatalf("Could not open the file for testing:, %v", err)
			}

			chunks, err := splitter.Split(file)

			require.Equal(t, tc.ExpectedError, err)
			require.Equal(t, tc.ExpectedChunks, len(*chunks))
		})
	}
	err = data.Afs.RemoveAll("testdata/chunks/")
	if err != nil {
		t.Errorf("Could not remove the test chunk directory: %v", err)
	}
}

func TestNewSplitJob(t *testing.T) {
	testCases := []struct {
		Name              string
		ChunkWriter       data.ChunkWriter
		ChunkMBS          uint64
		ChunkDir          string
		ChunkHashFunc     data.ChunkHashFunc
		ExpectedChunkSize uint64
		ExpectedChunkDir  string
	}{
		{
			"new split job with 0 chunkMBs as the input",
			mockChunkWriter.Execute,
			0,
			"testdata",
			data.DefaultChunkHashFunc,
			0x100000,
			"testdata",
		},
		{
			"new split job with 10 chunkMBs as the input",
			mockChunkWriter.Execute,
			10,
			"testdata",
			data.DefaultChunkHashFunc,
			0x100000 * 10,
			"testdata",
		},
	}
	for _, tc := range testCases {
		t.Run(tc.Name, func(t *testing.T) {
			actual := data.NewSplitJob(tc.ChunkWriter, tc.ChunkHashFunc, tc.ChunkDir, tc.ChunkMBS)
			require.True(t, reflect.ValueOf(actual.HashFunc).IsValid(), "HashFunc on Splitter was not set correctly")
			require.True(t, reflect.ValueOf(actual.ChunkWriter).IsValid(), "ChunkWriter on Splitter was not set correctly")
			require.Equal(t, tc.ExpectedChunkSize, actual.ChunkSize)
			require.Equal(t, tc.ExpectedChunkDir, actual.ChunkDir)
		})
	}
}
