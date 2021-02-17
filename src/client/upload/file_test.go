package upload_test

import (
	"github.com/spf13/afero"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"gitlab.computing.dcu.ie/collint9/2021-ca400-collint9-coynemt2/src/client/data"
	dataMocks "gitlab.computing.dcu.ie/collint9/2021-ca400-collint9-coynemt2/src/client/data/mocks"
	"gitlab.computing.dcu.ie/collint9/2021-ca400-collint9-coynemt2/src/client/upload"
	uploadMocks "gitlab.computing.dcu.ie/collint9/2021-ca400-collint9-coynemt2/src/client/upload/mocks"
	"testing"
)

func TestFileUpload_Upload(t *testing.T) {
	data.SetFileSystem(afero.NewMemMapFs())
	mockSplitter := new(dataMocks.Splitter)
	mockChunkUploader := new(uploadMocks.ChunkUploader)
	mockNewUploadRequester := new(uploadMocks.NewUploadRequester)
	fileUpload := upload.NewFileUpload(mockSplitter, mockChunkUploader, data.DefaultFileHashFunc, mockNewUploadRequester)
	file, err := data.Afs.Create("testFile")
	if err != nil {
		t.Fatalf("Could not create the test data file: %v", err)
	}

	testCases := []struct {
		Name          string
		FileData      []byte
		ExpectedError error
	}{
		{"empty file upload", []byte{}, nil},
		{"small file upload", []byte("small file content"), nil},
	}

	for _, tc := range testCases {
		t.Run(tc.Name, func(t *testing.T) {
			_, err = file.Write(tc.FileData)
			if err != nil {
				t.Fatalf("Could write to the test data file: %v", err)
			}
			mockSplitter.On("Split", mock.Anything).Return(&[]data.Chunk{}, nil)
			mockChunkUploader.On("UploadChunks", mock.Anything, mock.Anything).Return(nil)
			mockNewUploadRequester.On("CreateNewUpload", mock.Anything, mock.Anything, mock.Anything).Return("uploadRequestID", nil)

			err = fileUpload.Upload(file)
			require.Equal(t, tc.ExpectedError, err)
			err = file.Close()
			if err != nil {
				t.Errorf("Failed to close the test data file: %v", err)
			}
		})
	}
}
