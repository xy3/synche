package main_test

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/xy3/synche/src/client"
	"github.com/xy3/synche/src/client/mocks"
	"github.com/xy3/synche/src/cmd/client"
	"github.com/xy3/synche/src/config"
	"os"
	"testing"
)

var (
	textFile = "testdata/test_upload_file.txt"
)

func TestMain(m *testing.M) {
	config.TestMode = true
	os.Exit(m.Run())
}

func TestFileUpload(t *testing.T) {
	// TODO: Add test cases.
}

func TestUploadCommand(t *testing.T) {
	// Create a new mock FileUploadFunc and return no error for any string input
	fileUploader := new(mocks.FileUploadFunc)
	fileUploader.On("Execute", mock.AnythingOfType("string"), mock.AnythingOfType("uint")).Return(nil)
	uploadCmd := main.NewUploadCmd(fileUploader.Execute)

	mockAuth := new(mocks.AuthenticatorFunc)
	mockAuth.On("Execute", mock.AnythingOfType("string")).Return(nil)
	client.Authenticator = mockAuth.Execute

	testCases := []struct {
		Name     string
		Args     []string
		Expected error
	}{
		{"uploading a text file", []string{textFile}, nil},
		{"uploading with a specified name", []string{textFile, "-n", "test_name.pdf"}, nil},
	}
	for _, tc := range testCases {
		t.Run(tc.Name, func(t *testing.T) {
			uploadCmd.SetArgs(tc.Args)
			actual := uploadCmd.Execute()
			assert.Equal(t, tc.Expected, actual)
			mockAuth.AssertExpectations(t)
		})
	}
}
