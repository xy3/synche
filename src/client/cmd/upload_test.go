package cmd

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"testing"
)

var (
	textFile = "testdata/test_upload_file.txt"
)

type MockFileUploader struct {
	mock.Mock
}

func (m *MockFileUploader) Upload(filePath string) error {
	m.Called(filePath)
	return nil
}

func TestUploadCommand(t *testing.T) {
	// Create a new mock FileUploader and return no error for any string input
	fileUploader := new(MockFileUploader)
	fileUploader.On("Upload", mock.AnythingOfType("string")).Return(nil)
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
