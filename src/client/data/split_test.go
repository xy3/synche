package data

import (
	"testing"
)

func TestSplitJob_Split(t *testing.T) {
	//chunkWriter := new(mocks.ChunkWriter)
	//chunkWriter.On("Execute", mock.AnythingOfType("string"), mock.Anything).Return(nil)
	//splitter := NewSplitJob(chunkWriter.Execute)
	//
	//testCases := []struct {
	//	Name           string
	//	Filepath       string
	//	ChunkMBS       int
	//	ExpectedChunks int
	//	ExpectedError  error
	//}{
	//	{"splitting into smaller chunks", "testdata/split_me", 1, 3, nil},
	//	{"splitting into 1 same size chunk", "testdata/split_me", 10, 1, nil},
	//}
	//for _, tc := range testCases {
	//	t.Run(fmt.Sprintf("%v ", tc.Name), func(t *testing.T) {
	//		splitJob := &SplitJob{
	//			chunkWriter: chunkWriter.Execute,
	//			hashFunc:    nil,
	//			chunkDir:    "",
	//			chunkSize:   0,
	//			file:        nil,
	//		}
	//
	//		chunks, err := splitter.Split(tc.Filepath, "", tc.ChunkMBS)
	//
	//		assert.Equal(t, tc.ExpectedError, err)
	//		assert.Equal(t, tc.ExpectedChunks, len(chunks))
	//	})
	//}
}
