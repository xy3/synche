package data_test

import (
	"errors"
	"github.com/stretchr/testify/require"
	"gitlab.computing.dcu.ie/collint9/2021-ca400-collint9-coynemt2/src/server/data"
	"gitlab.computing.dcu.ie/collint9/2021-ca400-collint9-coynemt2/src/server/data/mocks"
	"gitlab.computing.dcu.ie/collint9/2021-ca400-collint9-coynemt2/src/server/data/schema"
	"strconv"
	"testing"
)

func Test_uploadCache_DelUpload(t *testing.T) {
	tests := []struct {
		name     string
		uploadId uint
		wantErr  bool
	}{
		{name: "delete upload from cache", uploadId: 1, wantErr: false},
		{name: "handles failed delete", uploadId: 1, wantErr: true},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			c := data.NewUploadCache()
			mockCache := new(mocks.Cache)
			cmd := mockCache.On("Delete", c.FormatKey(tc.uploadId))
			if tc.wantErr {
				cmd.Return(nil, errors.New(""))
			} else {
				cmd.Return(nil, nil)
			}
			err := c.DelUpload(mockCache, tc.uploadId)
			mockCache.AssertExpectations(t)
			if tc.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func Test_uploadCache_GetUpload(t *testing.T) {
	tests := []struct {
		name       string
		uploadId   uint
		wantUpload schema.Upload
		wantErr    bool
	}{
		{name: "get upload by id", wantUpload: schema.Upload{NumChunks: 10, FileID: 99, DirectoryID: 23}},
		{name: "handles failed GetAll call", wantUpload: schema.Upload{}, wantErr: true},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			c := data.NewUploadCache()
			mockCache := new(mocks.Cache)
			var returnErr error
			if tc.wantErr {
				returnErr = errors.New("")
			}

			mockCache.On("GetAll", c.FormatKey(tc.uploadId)).Return([]interface{}{
				[]byte("NumChunks"), strconv.FormatInt(tc.wantUpload.NumChunks, 10),
				[]byte("DirectoryID"), strconv.Itoa(int(tc.wantUpload.DirectoryID)),
				[]byte("FileID"), strconv.Itoa(int(tc.wantUpload.FileID)),
			}, returnErr)

			gotUpload, err := c.GetUpload(mockCache, tc.uploadId)
			if tc.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
			mockCache.AssertExpectations(t)
			require.Equal(t, tc.wantUpload.FileID, gotUpload.FileID)
			require.Equal(t, tc.wantUpload.DirectoryID, gotUpload.DirectoryID)
			require.Equal(t, tc.wantUpload.NumChunks, gotUpload.NumChunks)
			require.Equal(t, tc.wantUpload, gotUpload)
		})
	}
}

func Test_uploadCache_SetUpload(t *testing.T) {
	tests := []struct {
		name     string
		uploadId uint
		upload   schema.Upload
		wantErr  bool
	}{
		{
			name:     "set upload to upload id",
			uploadId: 1,
			upload:   schema.Upload{},
			wantErr:  false,
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			c := data.NewUploadCache()
			mockCache := new(mocks.Cache)
			mockCache.On("SetAll", c.FormatKey(tc.uploadId), tc.upload).Return(nil, nil)
			err := c.SetUpload(mockCache, tc.uploadId, tc.upload)
			mockCache.AssertExpectations(t)
			if tc.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}
