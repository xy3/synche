package repo

import (
	"github.com/spf13/afero"
	"github.com/stretchr/testify/suite"
	"github.com/xy3/synche/src/files"
	"github.com/xy3/synche/src/server/models"
	"github.com/xy3/synche/src/server/schema"
	"gorm.io/gorm"
	"os"
	"testing"
)

func TestMain(m *testing.M) {
	files.SetFileSystem(afero.NewMemMapFs())
	os.Exit(m.Run())
}

type listTestSuite struct {
	suite.Suite
	user    *schema.User
	homeDir *schema.Directory
	down    func(t *testing.T)
	db      *gorm.DB
}

func Test_listTestSuite(t *testing.T) {
	s := new(listTestSuite)
	suite.Run(t, s)
}

func (s *listTestSuite) SetupTest() {
	user, homeDir, db, down, err := NewUserForTest(s.T())
	s.Assert().NoError(err)
	s.down = down
	s.homeDir = homeDir
	s.user = user
	s.db = db
}

func addTestFilesToDir(dirID, userID uint, db *gorm.DB) (testFiles []schema.File, err error) {
	file1 := schema.File{
		Name:           "file1",
		Size:           10,
		Hash:           "hash1",
		DirectoryID:    dirID,
		UserID:         userID,
		TotalChunks:    10,
		ChunksReceived: 9,
	}
	file2 := file1
	file2.Hash = "hash2"
	file2.Name = "file2"
	if err = db.Create(&file1).Error; err != nil {
		return nil, err
	}
	if err = db.Create(&file2).Error; err != nil {
		return nil, err
	}
	testFiles = []schema.File{file1, file2}
	return testFiles, err
}

func (s *listTestSuite) TestGetDirContentsByID() {
	defer s.down(s.T())
	// tx, down := database.NewTxForTest(s.T())
	// s.db = tx
	// defer down(s.T())

	testDir, err := CreateDirectory("testDir", s.homeDir.ID, s.db)
	s.Require().NoError(err)

	want := &models.DirectoryContents{
		CurrentDir:     testDir.ConvertToModelsDir(),
		Files:          []*models.File{},
		SubDirectories: []*models.Directory{},
	}

	s.Run("List empty dir", func() {
		got, err := GetDirContentsByID(testDir.ID, s.db)
		s.Assert().NoError(err)
		s.Assert().EqualValues(want, got)
	})

	s.Run("List dir with files", func() {
		testFiles, err := addTestFilesToDir(testDir.ID, s.user.ID, s.db)
		s.Assert().NoError(err)
		want.Files = []*models.File{
			testFiles[0].ConvertToFileModel(),
			testFiles[1].ConvertToFileModel(),
		}
		got, err := GetDirContentsByID(testDir.ID, s.db)
		s.Assert().NoError(err)
		s.Assert().EqualValues(want, got)
	})
}

func (s *listTestSuite) TestGetDirWithContentsFromPath() {
	defer s.down(s.T())

	testDir, err := CreateDirectory("testDir", s.homeDir.ID, s.db)
	s.Require().NoError(err)

	testDir.Files = []schema.File{}
	testDir.Children = []schema.Directory{}
	want := testDir

	s.Run("List empty dir", func() {
		got, err := GetDirWithContentsFromPath(testDir.Path, s.db)
		s.Require().NoError(err)
		got.CreatedAt = testDir.CreatedAt
		got.UpdatedAt = testDir.UpdatedAt
		s.Assert().EqualValues(want, got)
	})

	s.Run("List dir with files", func() {
		_, err = addTestFilesToDir(testDir.ID, s.user.ID, s.db)
		s.Assert().NoError(err)

		got, err := GetDirWithContentsFromPath(testDir.Path, s.db)
		s.Assert().NoError(err)
		s.Assert().Equal(len(got.Files), 2)
	})
}

func (s *listTestSuite) TestGetHomeDirContents() {
	defer s.down(s.T())

	want := &models.DirectoryContents{
		CurrentDir:     s.homeDir.ConvertToModelsDir(),
		Files:          []*models.File{},
		SubDirectories: []*models.Directory{},
	}

	s.Run("List empty home dir", func() {
		got, err := GetHomeDirContents(s.user, s.db)
		s.Assert().NoError(err)
		s.Assert().EqualValues(want, got)
	})

	s.Run("List home dir with files", func() {
		testFiles, err := addTestFilesToDir(s.homeDir.ID, s.user.ID, s.db)
		s.Assert().NoError(err)
		want.Files = []*models.File{
			testFiles[0].ConvertToFileModel(),
			testFiles[1].ConvertToFileModel(),
		}
		got, err := GetHomeDirContents(s.user, s.db)
		s.Assert().NoError(err)
		s.Assert().EqualValues(want, got)
	})
}
