package cmd

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"gitlab.computing.dcu.ie/collint9/2021-ca400-collint9-coynemt2/src/client/cmd/mocks"
	"testing"
)

var (
	textFile = "testdata/test_upload_file.txt"
)

func TestUploadCommand(t *testing.T) {
	// Create a new mock FileUploader and return no error for any string input
	fileUploader := new(mocks.Uploader)
	fileUploader.On("Run", mock.AnythingOfType("string")).Return(nil)
	uploadCmd := NewUploadCmd(fileUploader)

	testCases := []struct {
		Name     string
		Args     []string
		Expected error
	}{
		{"uploading a text file", []string{textFile}, nil},
		{"uploading with a specified name", []string{textFile, "-n", "test_name.pdf"}, nil},
	}
	for _, tc := range testCases {
		t.Run(fmt.Sprintf("%v ", tc.Name), func(t *testing.T) {
			uploadCmd.SetArgs(tc.Args)
			actual := uploadCmd.Execute()
			assert.Equal(t, tc.Expected, actual)
		})
	}
}
