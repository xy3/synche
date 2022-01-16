package files_test

import (
	"bytes"
	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/require"
	"github.com/xy3/synche/src/files"
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
		{"1MB chunkSize and .5MB fileSize", files.MB, 0.5 * files.MB, 1},
		{"1MB chunkSize and 2MB fileSize", files.MB, 2 * files.MB, 2},
		{"10MB chunkSize and 3MB fileSize", 10 * files.MB, 3 * files.MB, 1},
		{"10MB chunkSize and 100MB fileSize", 10 * files.MB, 100 * files.MB, 10},
	}
	for _, tc := range testCases {
		t.Run(
			tc.Name, func(t *testing.T) {
				splitFile := new(files.SplitFile)
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
		{"empty file", 0, files.MB, 0, []byte{}, nil, nil},
		{"1MB file, 1MB chunkSize", 4 * files.BYTE, files.MB, 0, []byte("test"), []byte("test"), nil},
		{"current index at end", files.MB, files.MB, 1, []byte("test"), nil, nil},
		{"broken io reader", files.MB, files.MB, 0, nil, nil, io.EOF},
	}
	for _, tc := range testCases {
		t.Run(
			tc.Name, func(t *testing.T) {
				splitFile := new(files.SplitFile)
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
		Name              string
		FileSize          int64
		ChunkSize         int64
		CurrentIndex      int64
		FileBytes         []byte
		ExpectedNumChunks int
		ExpectedError     error
	}{
		{"empty file", 0, files.MB, 0, []byte{}, 0, nil},
		{"1MB file, 1MB chunkSize", 4 * files.BYTE, files.MB, 0, []byte("test"), 1, nil},
		{"current index at end", files.MB, files.MB, 1, []byte("test"), 0, nil},
		{"broken io reader", files.MB, files.MB, 0, nil, 0, io.EOF},
		{"5MB file, 1MB chunkSize", 5 * files.MB, 1 * files.MB, 0, make([]byte, (files.MB/files.BYTE)*5), 5, nil},
	}
	for _, tc := range testCases {
		t.Run(
			tc.Name, func(t *testing.T) {
				splitFile := new(files.SplitFile)
				splitFile.FileSize = tc.FileSize
				splitFile.ChunkSize = tc.ChunkSize
				splitFile.CurrentIndex = tc.CurrentIndex
				splitFile.Reader = bytes.NewReader(tc.FileBytes)
				count := 0
				err := splitFile.Split(
					func(chunk *files.Chunk, index int64) error {
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
