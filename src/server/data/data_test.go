package data_test

import (
	"errors"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"gitlab.computing.dcu.ie/collint9/2021-ca400-collint9-coynemt2/src/server/data"
	"gitlab.computing.dcu.ie/collint9/2021-ca400-collint9-coynemt2/src/server/data/mocks"
	"gitlab.computing.dcu.ie/collint9/2021-ca400-collint9-coynemt2/src/server/data/schema"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"regexp"
	"testing"
)

type dataSuite struct {
	suite.Suite
	DB         *gorm.DB
	dbMock     sqlmock.Sqlmock
	syncheData data.SyncheData
	cacheMock  *mocks.UploadCache
}

func (s *dataSuite) SetupSuite() {
	db, dbMock, err := sqlmock.New()
	require.NoError(s.T(), err)
	s.dbMock = dbMock
	s.dbMock.ExpectQuery("SELECT VERSION()").WillReturnRows(sqlmock.NewRows([]string{"VERSION()"}).AddRow(""))
	s.DB, err = gorm.Open(mysql.New(mysql.Config{Conn: db}))
	require.NoError(s.T(), err)

	s.syncheData = data.SyncheData{
		Cache: nil,
		DB:    s.DB,
	}
}

func (s *dataSuite) TestSyncheData_GetNumChunks() {
	tests := []struct {
		name            string
		uploadId        uint
		fromCache       bool
		chunksFromCache int64
		fromDB          bool
		chunksFromDB    int64
		wantChunks      int64
		wantErr         bool
	}{
		{name: "from cache", uploadId: 1, fromCache: true, chunksFromCache: 10, wantChunks: 10},
		{name: "from SQL db", uploadId: 1, fromDB: true, chunksFromDB: 5, wantChunks: 5},
		{name: "when uploadID does not exist", uploadId: 1, wantChunks: 0, wantErr: true},
		{
			name: "cache num is returned first if its in cache",
			uploadId: 500,
			fromCache: true,
			chunksFromCache: 5,
			chunksFromDB: 10,
			wantChunks: 5,
			wantErr: false,
		},
		{
			name: "DB result is returned if num is not in cache",
			uploadId: 100,
			fromDB: true,
			chunksFromDB: 10,
			wantChunks: 10,
			wantErr: false,
		},
	}

	for _, tc := range tests {
		s.cacheMock = new(mocks.UploadCache)

		s.Run(tc.name, func() {

			get := s.cacheMock.On("GetUpload", s.syncheData.Cache, tc.uploadId)
			set := s.cacheMock.On("SetUpload", s.syncheData.Cache, tc.uploadId, mock.AnythingOfType("models.Upload"))
			upload := schema.Upload{}

			if tc.fromCache {
				upload.NumChunks = tc.chunksFromCache
				get.Return(upload, nil)
			} else {
				get.Return(upload, errors.New(""))
				set.Return(nil)

				query := "SELECT * FROM `uploads` WHERE `uploads`.`id` = ? AND `uploads`.`deleted_at` IS NULL ORDER BY `uploads`.`id` LIMIT 1"
				expectedQuery := s.dbMock.ExpectQuery(regexp.QuoteMeta(query)).WithArgs(tc.uploadId)

				if tc.fromDB {
					expectedQuery.WillReturnRows(sqlmock.NewRows([]string{"num_chunks"}).AddRow(tc.chunksFromDB))
				} else {
					expectedQuery.WillReturnError(gorm.ErrRecordNotFound)
				}
			}

			gotChunks, err := s.syncheData.GetNumChunks(s.cacheMock, tc.uploadId)

			require.Equal(s.T(), err != nil, tc.wantErr)
			require.Equal(s.T(), tc.wantChunks, gotChunks)
		})
	}
}

// Todo: Add migration test cases and configure the sqlmock to expect the queries

func TestDataSuite(t *testing.T) {
	suite.Run(t, new(dataSuite))
}
