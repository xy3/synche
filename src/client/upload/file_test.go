package upload_test

//func TestAsyncUpload(t *testing.T) {
//	tests := []struct {
//		name    string
//		wantErr bool
//	}{
//		{"basic test", false},
//	}
//	for _, tt := range tests {
//		t.Run(tt.name, func(t *testing.T) {
//			mockSplitter := new(dataMocks.Splitter)
//			mockSplitter.On("Split", mock.Anything).Return(nil)
//			mockSplitter.On("File").Return(new(data.SplitFile))
//			mockSplitter.On("NumChunks").Return(int64(1))
//
//			mockNewUploadFunc := new(mocks.NewUploadFunc)
//			mockNewUploadFunc.On("Execute", mock.Anything).Return(nil, nil)
//			mockChunkUploader := new(mocks.AsyncChunkUploader)
//
//			err := upload.AsyncUpload(mockSplitter, 0, mockNewUploadFunc.Execute, mockChunkUploader.Execute)
//			if (err != nil) != tt.wantErr {
//				t.Errorf("AsyncUpload() error = %v, wantErr %v", err, tt.wantErr)
//			}
//		})
//	}
//}
