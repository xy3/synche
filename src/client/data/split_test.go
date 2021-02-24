package data_test

import (
	"bytes"
	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/require"
	"gitlab.computing.dcu.ie/collint9/2021-ca400-collint9-coynemt2/src/client/data"
	"io"
	"io/ioutil"
	"testing"
)

func init() {
	log.SetOutput(ioutil.Discard)
}

func TestSplitFile_NumChunks(t *testing.T) {
	testCases := []struct {
		Name           string
		ChunkSize      int64
		FileSize       int64
		ExpectedChunks int64
	}{
		{"1MB chunkSize and .5MB fileSize", data.MB, 0.5 * data.MB, 1},
		{"1MB chunkSize and 2MB fileSize", data.MB, 2 * data.MB, 2},
		{"10MB chunkSize and 3MB fileSize", 10 * data.MB, 3 * data.MB, 1},
		{"10MB chunkSize and 100MB fileSize", 10 * data.MB, 100 * data.MB, 10},
	}
	for _, tc := range testCases {
		t.Run(
			tc.Name, func(t *testing.T) {
				splitFile := new(data.SplitFile)
				splitFile.FileSize = tc.FileSize
				splitFile.ChunkSize = tc.ChunkSize
				require.Equal(t, tc.ExpectedChunks, splitFile.NumChunks())
			},
		)
	}
}

func TestSplitFile_NextChunk(t *testing.T) {
	testCases := []struct {
		Name          string
		FileSize      int64
		ChunkSize     int64
		CurrentIndex  int64
		FileBytes     []byte
		ExpectedBytes []byte
		ExpectedError error
	}{
		{"empty file", 0, data.MB, 0, []byte{}, nil, nil},
		{"1MB file, 1MB chunkSize", 4 * data.BYTE, data.MB, 0, []byte("test"), []byte("test"), nil},
		{"current index at end", data.MB, data.MB, 1, []byte("test"), nil, nil},
		{"broken io reader", data.MB, data.MB, 0, nil, nil, io.EOF},
	}
	for _, tc := range testCases {
		t.Run(
			tc.Name, func(t *testing.T) {
				splitFile := new(data.SplitFile)
				splitFile.FileSize = tc.FileSize
				splitFile.ChunkSize = tc.ChunkSize
				splitFile.CurrentIndex = tc.CurrentIndex
				splitFile.Reader = bytes.NewReader(tc.FileBytes)
				chunkBytes, err := splitFile.NextChunk()
				require.Equal(t, tc.ExpectedError, err)
				if err == nil {
					require.Equal(t, tc.ExpectedBytes, chunkBytes)
				}
			},
		)
	}
}

func TestSplitFile_Split(t *testing.T) {
	testCases := []struct {
		Name          string
		FileSize      int64
		ChunkSize     int64
		CurrentIndex  int64
		FileBytes     []byte
		ExpectedNumChunks int
		ExpectedError error
	}{
		{"empty file", 0, data.MB, 0, []byte{}, 0, nil},
		{"1MB file, 1MB chunkSize", 4 * data.BYTE, data.MB, 0, []byte("test"), 1, nil},
		{"current index at end", data.MB, data.MB, 1, []byte("test"), 0, nil},
		{"broken io reader", data.MB, data.MB, 0, nil, 0, io.EOF},
		{"5MB file, 1MB chunkSize", 5*data.MB, 1*data.MB, 0, make([]byte, (data.MB/data.BYTE) * 5), 5, nil},
	}
	for _, tc := range testCases {
		t.Run(
			tc.Name, func(t *testing.T) {
				splitFile := new(data.SplitFile)
				splitFile.FileSize = tc.FileSize
				splitFile.ChunkSize = tc.ChunkSize
				splitFile.CurrentIndex = tc.CurrentIndex
				splitFile.Reader = bytes.NewReader(tc.FileBytes)
				count := 0
				err := splitFile.Split(
					func(chunk *data.Chunk) error {
						count++
						return nil
					},
				)
				require.Equal(t, tc.ExpectedNumChunks, count)
				require.Equal(t, tc.ExpectedError, err)
			},
		)
	}
}
