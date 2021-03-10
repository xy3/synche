package upload_test

import (
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	dataMocks "gitlab.computing.dcu.ie/collint9/2021-ca400-collint9-coynemt2/src/client/data/mocks"
	"gitlab.computing.dcu.ie/collint9/2021-ca400-collint9-coynemt2/src/client/upload"
	uploadMocks "gitlab.computing.dcu.ie/collint9/2021-ca400-collint9-coynemt2/src/client/upload/mocks"
	"testing"
)

// TODO: This test seems a bit pointless since almost everything is being mocked - is there a way to improve it?
func TestFileUpload_Upload(t *testing.T) {
	mockSplitter := new(dataMocks.Splitter)
	mockNewUploadRequester := new(uploadMocks.NewUploadRequester)
	fileUpload := upload.NewFileUpload(new(uploadMocks.ChunkUploader), mockNewUploadRequester)
	testCases := []struct {
		Name          string
		ExpectedError error
	}{
		{"empty file upload", nil},
		{"small file upload", nil},
	}
	for _, tc := range testCases {
		t.Run(tc.Name, func(t *testing.T) {
			mockSplitter.On("NumChunks").Return(int64(1))
			mockSplitter.On("Split", mock.Anything).Return(nil)
			mockNewUploadRequester.On("CreateNewUpload", mock.Anything).Return("uploadRequestID", nil)

			err := fileUpload.AsyncUpload(mockSplitter)
			require.Equal(t, tc.ExpectedError, err)
			if err != nil {
				t.Errorf("Failed to close the test data file: %v", err)
			}
		})
	}
}
