package schema_test

import (
	"github.com/spf13/afero"
	"github.com/stretchr/testify/suite"
	"github.com/xy3/synche/src/files"
	schema2 "github.com/xy3/synche/src/schema"
	"github.com/xy3/synche/src/server/models"
	"github.com/xy3/synche/src/server/repo"
	"gorm.io/gorm"
	"strings"
	"testing"
)

type fileTestSuite struct {
	suite.Suite
	user     *schema2.User
	homeDir  *schema2.Directory
	file     *schema2.File
	filePath string
	down     func(t *testing.T)
	db       *gorm.DB
}

func Test_fileTestSuite(t *testing.T) {
	files.SetFileSystem(afero.NewMemMapFs())
	suite.Run(t, new(fileTestSuite))
}

func (s *fileTestSuite) SetupTest() {
	user, homeDir, db, down, err := repo.NewUserForTest(s.T())
	s.Assert().NoError(err)
	s.down = down
	s.homeDir = homeDir
	s.user = user
	s.db = db

	path, err := repo.BuildFullPath("testfile", s.user, s.db)
	s.Assert().NoError(err)
	file, err := repo.CreateFileFromReader(path, strings.NewReader("some content"), s.user.ID, s.db)
	s.Assert().NoError(err)

	s.filePath = path
	s.file = file
}

func (s *fileTestSuite) TestFile_AppendFromReader() {
	defer s.down(s.T())

	s.Run("Existing file", func() {
		err := s.file.AppendFromReader(strings.NewReader(" extra"), s.user.ID, s.db)
		s.Assert().NoError(err)
		readFile, err := files.Afs.ReadFile(s.filePath)
		s.Assert().NoError(err)
		s.Assert().Equal([]byte("some content extra"), readFile)
	})

	s.Run("Wrong user ID", func() {
		err := s.file.AppendFromReader(strings.NewReader("extra"), 1234, s.db)
		s.Assert().Error(err)
	})

	s.Run("Non-existing file", func() {
		err := s.file.Delete(s.db)
		s.Assert().NoError(err)
		err = s.file.AppendFromReader(strings.NewReader("extra"), s.user.ID, s.db)
		s.Assert().Error(err)
	})
}

func (s *fileTestSuite) TestFile_ConvertToFileModel() {
	defer s.down(s.T())
	got := s.file.ConvertToFileModel()
	want := &models.File{
		Available:      true,
		ChunksReceived: 0,
		DirectoryID:    1,
		Hash:           "0cbb30dba49e2bff48b93fea523ed308",
		ID:             1,
		Name:           "testfile",
		Size:           12,
		TotalChunks:    0,
	}
	s.Assert().EqualValues(want, got)
}

func (s *fileTestSuite) TestFile_Delete() {
	defer s.down(s.T())
	err := s.file.Delete(s.db)
	s.Assert().NoError(err)
}

func (s *fileTestSuite) TestFile_Move() {
	defer s.down(s.T())
	err := s.file.Move(s.filePath+"/dir/newfile", s.db)
	s.Assert().NoError(err)
}

func (s *fileTestSuite) TestFile_Path() {
	defer s.down(s.T())
	got, err := s.file.Path(s.db)
	s.Assert().NoError(err)
	s.Assert().Equal(s.filePath, got)
}

func (s *fileTestSuite) TestFile_Reader() {
	defer s.down(s.T())
	got, err := s.file.Reader(s.db)
	s.Assert().NoError(err)
	gotBytes, err := afero.ReadAll(got)
	s.Assert().NoError(err)
	s.Assert().Equal([]byte("some content"), gotBytes)
}

func (s *fileTestSuite) TestFile_Rename() {
	defer s.down(s.T())
	err := s.file.Rename("newName", s.db)
	s.Assert().NoError(err)
	s.Assert().Equal("newName", s.file.Name)
}

func (s *fileTestSuite) TestFile_SetAvailable() {
	defer s.down(s.T())
	s.file.Available = false
	err := s.file.SetAvailable(s.db)
	s.Assert().NoError(err)
	s.Assert().True(s.file.Available)
}

func (s *fileTestSuite) TestFile_SetUnavailable() {
	defer s.down(s.T())
	s.file.Available = true
	err := s.file.SetUnavailable(s.db)
	s.Assert().NoError(err)
	s.Assert().False(s.file.Available)
}

func (s *fileTestSuite) TestFile_ValidateHash() {
	defer s.down(s.T())

	s.Run("Valid file hash", func() {
		err := s.file.ValidateHash(s.db)
		s.Assert().NoError(err)
	})

	s.Run("Invalid file hash", func() {
		s.file.Hash = "wrong-hash"
		err := s.file.ValidateHash(s.db)
		s.Assert().Error(err)
	})
}
