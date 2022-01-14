package server

import (
	"github.com/stretchr/testify/suite"
	"gorm.io/gorm"
	"testing"
)

type dataTestSuite struct {
	suite.Suite
	down func(t *testing.T)
	db   *gorm.DB
}

func Test_dataTestSuite(t *testing.T) {
	suite.Run(t, new(dataTestSuite))
}

func (s *dataTestSuite) SetupTest() {
	db, down := NewTxForTest(s.T())
	s.down = down
	s.db = db
	TestDB = s.db
}

func (s *dataTestSuite) TestInitSyncheData() {
	defer s.down(s.T())
	gotDb, err := InitSyncheData()
	s.Assert().NoError(err)
	s.Assert().NotNil(gotDb)
}

func (s *dataTestSuite) TestNewConnection() {
	defer s.down(s.T())
	gotDb, err := NewConnection()
	s.Assert().NoError(err)
	s.Assert().NotNil(gotDb)
}

func (s *dataTestSuite) TestRequireNewConnection() {
	defer s.down(s.T())
	gotDb := RequireNewConnection()
	s.Assert().NotNil(gotDb)
}

func (s *dataTestSuite) Test_configureConnection() {
	defer s.down(s.T())
	gotDb, err := configureConnection(s.db)
	s.Assert().NoError(err)
	s.Assert().NotNil(gotDb)
}

func (s *dataTestSuite) TestMigrateAll() {
	defer s.down(s.T())
	err := MigrateAll(s.db)
	s.Assert().NoError(err)
}
