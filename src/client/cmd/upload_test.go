package cmd_test

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"gitlab.computing.dcu.ie/collint9/2021-ca400-collint9-coynemt2/src/client/apiclient"
	apiClientMocks "gitlab.computing.dcu.ie/collint9/2021-ca400-collint9-coynemt2/src/client/apiclient/mocks"
	"gitlab.computing.dcu.ie/collint9/2021-ca400-collint9-coynemt2/src/client/cmd"
	"gitlab.computing.dcu.ie/collint9/2021-ca400-collint9-coynemt2/src/client/cmd/mocks"
	"gitlab.computing.dcu.ie/collint9/2021-ca400-collint9-coynemt2/src/config"
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
	uploadCmd := cmd.NewUploadCmd(fileUploader.Execute)

	mockAuth := new(apiClientMocks.AuthenticatorFunc)
	mockAuth.On("Execute", mock.AnythingOfType("string")).Return(nil)
	apiclient.Authenticator = mockAuth.Execute

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
