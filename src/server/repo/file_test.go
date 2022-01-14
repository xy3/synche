package repo

import (
	"github.com/stretchr/testify/suite"
	"github.com/xy3/synche/src/server/schema"
	"gorm.io/gorm"
	"path/filepath"
	"strings"
	"testing"
)

type fileTestSuite struct {
	suite.Suite
	user    *schema.User
	homeDir *schema.Directory
	down    func(t *testing.T)
	dbc     *gorm.DB
	db      *gorm.DB
}

func Test_fileTestSuite(t *testing.T) {
	s := new(fileTestSuite)

	// user, homeDir, db, down, err := NewUserForTest(t)
	// // s.Assert().NoError(err)
	// assert.NoError(t, err)
	// s.down = down
	// s.homeDir = homeDir
	// s.user = user
	// s.dbc = db

	suite.Run(t, s)
}

func (s *fileTestSuite) SetupTest() {
	user, homeDir, db, down, err := NewUserForTest(s.T())
	s.Assert().NoError(err)
	s.down = down
	s.homeDir = homeDir
	s.user = user
	s.db = db
}

func (s *fileTestSuite) TestCreateFileFromReader() {
	defer s.down(s.T())

	path, err := BuildFullPath("testfile.txt", s.user, s.db)
	s.Assert().NoError(err)
	reader := strings.NewReader("file content")

	s.Run("Create file from string reader", func() {
		gotFile, err := CreateFileFromReader(path, reader, s.user.ID, s.db)
		s.Assert().NoError(err)
		s.Assert().Equal("testfile.txt", gotFile.Name)
		s.Assert().Equal(int64(12), gotFile.Size)
		s.Assert().Equal("0ceff6c19f5e2d0520ea3c145a250510", gotFile.Hash)
		s.Assert().Equal(s.homeDir.ID, gotFile.DirectoryID)
		s.Assert().Equal(s.user.ID, gotFile.UserID)
	})

	s.Run("Try to create duplicate file", func() {
		_, err = CreateFileFromReader(path, reader, s.user.ID, s.db)
		s.Assert().Error(err)
	})

	s.Run("Create file in nested dir", func() {
		path, err = BuildFullPath("dirname/testfile.txt", s.user, s.db)
		s.Assert().NoError(err)
		reader = strings.NewReader("new file content")

		gotFile, err := CreateFileFromReader(path, reader, s.user.ID, s.db)
		s.Assert().NoError(err)
		s.Assert().Equal(int64(16), gotFile.Size)
		s.Assert().Equal("10121ba5df6f1e38b19bebc9c3102600", gotFile.Hash)
		s.Assert().NotEqual(s.homeDir.ID, gotFile.DirectoryID)
		s.Assert().Equal(s.user.ID, gotFile.UserID)
	})
}

func (s *fileTestSuite) TestFindFileByFullPath() {
	defer s.down(s.T())

	path, err := BuildFullPath("testfile.txt", s.user, s.db)
	s.Assert().NoError(err)
	reader := strings.NewReader("file content")

	s.Run("Non-existing file", func() {
		_, err := FindFileByFullPath(path, s.db)
		s.Assert().Error(err)
	})

	s.Run("Existing file", func() {
		newFile, err := CreateFileFromReader(path, reader, s.user.ID, s.db)
		s.Assert().NoError(err)

		gotFile, err := FindFileByFullPath(path, s.db)
		s.Assert().NoError(err)

		// These values need to be set because timestamps change when they are accessed
		gotFile.UpdatedAt = newFile.UpdatedAt
		gotFile.CreatedAt = newFile.CreatedAt
		gotFile.Directory = newFile.Directory

		s.Assert().EqualValues(newFile, gotFile)
	})
}

func (s *fileTestSuite) TestGetFileByID() {
	defer s.down(s.T())

	path, err := BuildFullPath("testfile.txt", s.user, s.db)
	s.Assert().NoError(err)
	reader := strings.NewReader("file content")

	s.Run("Non-existing file", func() {
		_, err := GetFileByID(2, s.db)
		s.Assert().Error(err)
	})

	s.Run("Existing file", func() {
		newFile, err := CreateFileFromReader(path, reader, s.user.ID, s.db)
		s.Assert().NoError(err)

		gotFile, err := GetFileByID(newFile.ID, s.db)
		s.Assert().NoError(err)

		// These values need to be set because timestamps change when they are accessed
		gotFile.UpdatedAt = newFile.UpdatedAt
		gotFile.CreatedAt = newFile.CreatedAt

		s.Assert().EqualValues(newFile, gotFile)
	})
}

func (s *fileTestSuite) TestMoveFile() {
	defer s.down(s.T())

	path, err := BuildFullPath("testfile.txt", s.user, s.db)
	s.Assert().NoError(err)
	reader := strings.NewReader("file content")

	s.Run("Non-existing file", func() {
		err = MoveFile(&schema.File{}, "some/path", s.db)
		s.Assert().Error(err)
	})

	s.Run("Existing file", func() {
		newFile, err := CreateFileFromReader(path, reader, s.user.ID, s.db)
		s.Assert().NoError(err)

		oldDirID := newFile.DirectoryID

		err = MoveFile(newFile, filepath.Join(path, "newdir/newfile.txt"), s.db)
		s.Assert().NoError(err)
		s.Assert().NotEqual(newFile.DirectoryID, oldDirID)
		s.Assert().Equal("newfile.txt", newFile.Name)
	})
}

func (s *fileTestSuite) TestRenameFile() {
	defer s.down(s.T())

	path, err := BuildFullPath("testfile.txt", s.user, s.db)
	s.Assert().NoError(err)
	reader := strings.NewReader("file content")

	s.Run("Non-existing file", func() {
		_, err = RenameFile(10, "new-name.txt", s.db)
		s.Assert().Error(err)
	})

	s.Run("Existing file", func() {
		newFile, err := CreateFileFromReader(path, reader, s.user.ID, s.db)
		s.Assert().NoError(err)

		gotFile, err := RenameFile(newFile.ID, "new-name.txt", s.db)
		s.Assert().NoError(err)
		s.Assert().Equal("new-name.txt", gotFile.Name)
		s.Assert().Equal(newFile.DirectoryID, gotFile.DirectoryID)
		s.Assert().Equal(newFile.Size, gotFile.Size)
		s.Assert().Equal(newFile.ID, gotFile.ID)
		s.Assert().Equal(newFile.Hash, gotFile.Hash)
	})
}

func (s *fileTestSuite) Test_writeFileData() {
	defer s.down(s.T())

	path, err := BuildFullPath("testfile.txt", s.user, s.db)
	s.Assert().NoError(err)
	reader := strings.NewReader("file content")

	got, err := writeFileData(path, reader)
	s.Assert().NoError(err)
	s.Assert().Equal(int64(12), got)
}
