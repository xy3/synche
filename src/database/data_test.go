package database_test

import (
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"gitlab.computing.dcu.ie/collint9/2021-ca400-collint9-coynemt2/src/database"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"testing"
)

type dataSuite struct {
	suite.Suite
	DB     *gorm.DB
	mockDB sqlmock.Sqlmock
}

func (s *dataSuite) SetupSuite() {
	db, mockDB, err := sqlmock.New()
	require.NoError(s.T(), err)

	s.mockDB = mockDB
	s.mockDB.ExpectQuery("SELECT VERSION()").WillReturnRows(sqlmock.NewRows([]string{"VERSION()"}).AddRow(""))
	s.DB, err = gorm.Open(mysql.New(mysql.Config{Conn: db}))

	require.NoError(s.T(), err)

	database.DB = s.DB
}

// Todo: Add migration test cases and configure the sqlmock to expect the queries
func TestDataSuite(t *testing.T) {
	suite.Run(t, new(dataSuite))
}
