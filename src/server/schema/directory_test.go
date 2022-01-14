package schema_test

import (
	"github.com/spf13/afero"
	"github.com/stretchr/testify/suite"
	"github.com/xy3/synche/src/client/models"
	"github.com/xy3/synche/src/files"
	"github.com/xy3/synche/src/hash"
	"github.com/xy3/synche/src/server/repo"
	"github.com/xy3/synche/src/server/schema"
	"gorm.io/gorm"
	"testing"
)

type directoryTestSuite struct {
	suite.Suite
	user    *schema.User
	homeDir *schema.Directory
	dir     *schema.Directory
	down    func(t *testing.T)
	db      *gorm.DB
}

func Test_directoryTestSuite(t *testing.T) {
	files.SetFileSystem(afero.NewMemMapFs())
	suite.Run(t, new(directoryTestSuite))
}

func (s *directoryTestSuite) SetupTest() {
	user, homeDir, db, down, err := repo.NewUserForTest(s.T())
	s.Assert().NoError(err)
	s.down = down
	s.homeDir = homeDir
	s.user = user
	s.db = db

	dir, err := repo.CreateDirectory("newDir", s.homeDir.ID, s.db)
	s.Assert().NoError(err)

	s.dir = dir
}

func (s *directoryTestSuite) TestDirectory_ConvertToModelsDir() {
	defer s.down(s.T())
	want := &models.Directory{
		FileCount:         0,
		ID:                uint64(s.dir.ID),
		Name:              "newDir",
		ParentDirectoryID: uint64(s.homeDir.ID),
		Path:              s.dir.Path,
		PathHash:          hash.PathHash(s.dir.Path),
	}
	got := s.dir.ConvertToModelsDir()
	s.Assert().EqualValues(want, got)
}

func (s *directoryTestSuite) TestDirectory_Delete() {
	defer s.down(s.T())
	err := s.dir.Delete(true, s.db)
	s.Assert().NoError(err)
}

func (s *directoryTestSuite) TestDirectory_UpdateFileCount() {
	defer s.down(s.T())
	gotNum, err := s.dir.UpdateFileCount(s.db)
	s.Assert().NoError(err)
	s.Assert().Zero(gotNum)
}
