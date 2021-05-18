package repo

import (
	"github.com/spf13/afero"
	"github.com/stretchr/testify/suite"
	"gitlab.computing.dcu.ie/collint9/2021-ca400-collint9-coynemt2/src/database"
	"gitlab.computing.dcu.ie/collint9/2021-ca400-collint9-coynemt2/src/database/schema"
	"gitlab.computing.dcu.ie/collint9/2021-ca400-collint9-coynemt2/src/files"
	"gitlab.computing.dcu.ie/collint9/2021-ca400-collint9-coynemt2/src/files/hash"
	c "gitlab.computing.dcu.ie/collint9/2021-ca400-collint9-coynemt2/src/server/config"
	"gorm.io/gorm"
	"path/filepath"
	"strings"
	"testing"
)

type directoryTestSuite struct {
	suite.Suite
	user    *schema.User
	homeDir *schema.Directory
	down    func(t *testing.T)
	db      *gorm.DB
}

func Test_directoryTestSuite(t *testing.T) {
	files.SetFileSystem(afero.NewMemMapFs())
	s := new(directoryTestSuite)
	suite.Run(t, s)
}

func (s *directoryTestSuite) SetupTest() {
	user, homeDir, db, down, err := NewUserForTest(s.T())
	s.Assert().NoError(err)
	s.down = down
	s.homeDir = homeDir
	s.user = user
	s.db = db
}

func (s *directoryTestSuite) TestBuildFullPath() {
	defer s.down(s.T())

	s.Run("Correct path and user", func() {
		got, err := BuildFullPath("/somepath", s.user, s.db)
		s.Assert().NoError(err)
		s.Assert().NotNil(got)
		s.Assert().Equal(filepath.Join(s.homeDir.Path, "/somepath"), got)
	})

	s.Run("No user", func() {
		_, err := BuildFullPath("/somepath", &schema.User{}, s.db)
		s.Assert().Error(err)
	})

	s.Run("Empty path", func() {
		got, err := BuildFullPath("", s.user, s.db)
		s.Assert().NoError(err)
		s.Assert().Equal(s.homeDir.Path, got)
	})
}

func (s *directoryTestSuite) TestCreateDirectory() {
	defer s.down(s.T())

	s.Run("Dir in home dir", func() {
		gotDirectory, err := CreateDirectory("someDir", s.homeDir.ID, s.db)
		s.Assert().NoError(err)
		s.Assert().NotNil(gotDirectory)
		s.Assert().Equal(s.user.ID, gotDirectory.UserID)
		s.Assert().Equal(s.homeDir.ID, *gotDirectory.ParentID)
		s.Assert().Equal("someDir", gotDirectory.Name)
		wantPath := filepath.Join(s.homeDir.Path, "someDir")
		s.Assert().Equal(wantPath, gotDirectory.Path)
		wantPathHash := hash.PathHash(wantPath)
		s.Assert().Equal(wantPathHash, gotDirectory.PathHash)
	})

	s.Run("Empty name", func() {
		_, err := CreateDirectory("", s.homeDir.ID, s.db)
		s.Assert().Error(err)
	})

	s.Run("No parent dir", func() {
		_, err := CreateDirectory("someDir2", 1234, s.db)
		s.Assert().Error(err)
	})

	s.Run("Duplicate dir", func() {
		_, err := CreateDirectory("someDir", s.homeDir.ID, s.db)
		s.Assert().Error(err)
	})
}

func (s *directoryTestSuite) TestCreateDirectoryFromPath() {
	defer s.down(s.T())

	path, err := BuildFullPath("/someDir", s.user, s.db)
	s.Assert().NoError(err)

	var someDirID uint

	s.Run("New dir in home dir", func() {
		gotDir, err := CreateDirectoryFromPath(path, s.db)
		s.Assert().NoError(err)
		s.Assert().NotNil(gotDir)
		s.Assert().Equal("someDir", gotDir.Name)
		someDirID = gotDir.ID
	})

	s.Run("New dir in nested dir", func() {
		gotDir, err := CreateDirectoryFromPath(path+"/dir2", s.db)
		s.Assert().NoError(err)
		s.Assert().NotNil(gotDir)
		s.Assert().Equal("dir2", gotDir.Name)
		s.Assert().Equal(someDirID, *gotDir.ParentID)
		s.Assert().Equal(path+"/dir2", gotDir.Path)
		s.Assert().Equal(s.user.ID, gotDir.UserID)
	})

	s.Run("New dir in non-existing nested dir", func() {
		gotDir, err := CreateDirectoryFromPath(path+"/dir-that-doesnt-exist/dir3", s.db)
		s.Assert().NoError(err)
		s.Assert().NotNil(gotDir)
		s.Assert().Equal("dir3", gotDir.Name)

		parentDir, err := GetDirectoryByID(*gotDir.ParentID, s.db)
		s.Assert().NoError(err)

		s.Assert().Equal("dir-that-doesnt-exist", parentDir.Name)
		s.Assert().Equal(path+"/dir-that-doesnt-exist/dir3", gotDir.Path)
		s.Assert().Equal(s.user.ID, gotDir.UserID)
		s.Assert().Equal(s.user.ID, parentDir.UserID)
	})
}

func (s *directoryTestSuite) TestGenerateUserDirName() {
	defer s.down(s.T())

	got := GenerateUserDirName(s.user)
	s.Assert().Equal("synche_a26ee852bc", got)
}

func (s *directoryTestSuite) TestGetDirByPath() {
	defer s.down(s.T())

	path, err := BuildFullPath("/someDir", s.user, s.db)
	s.Assert().NoError(err)

	s.Run("Existing path", func() {
		newDir, err := CreateDirectoryFromPath(path, s.db)
		s.Assert().NoError(err)

		gotDir, err := GetDirByPath(path, s.db)
		s.Assert().NoError(err)
		gotDir.UpdatedAt = newDir.UpdatedAt
		gotDir.CreatedAt = newDir.CreatedAt
		s.Assert().EqualValues(newDir, gotDir)
	})

	s.Run("Non-existing path", func() {
		_, err := GetDirByPath("fakepath", s.db)
		s.Assert().Error(err)
	})
}

func (s *directoryTestSuite) TestGetDirectoryByID() {
	defer s.down(s.T())

	path, err := BuildFullPath("/someDir", s.user, s.db)
	s.Assert().NoError(err)

	s.Run("Existing dir ID", func() {
		newDir, err := CreateDirectoryFromPath(path, s.db)
		s.Assert().NoError(err)

		gotDir, err := GetDirectoryByID(newDir.ID, s.db)
		s.Assert().NoError(err)
		gotDir.UpdatedAt = newDir.UpdatedAt
		gotDir.CreatedAt = newDir.CreatedAt
		s.Assert().EqualValues(newDir, gotDir)
	})

	s.Run("Non-existing dir ID", func() {
		_, err := GetDirectoryByID(1234, s.db)
		s.Assert().Error(err)
	})
}

func (s *directoryTestSuite) TestGetDirectoryForFileID() {
	defer s.down(s.T())

	path, err := BuildFullPath("/someDir", s.user, s.db)
	s.Assert().NoError(err)
	newFile, err := CreateFileFromReader(path+"/testfile.txt", strings.NewReader("content"), s.user.ID, s.db)
	s.Assert().NoError(err)

	s.Run("Existing file ID", func() {
		gotDir, err := GetDirectoryForFileID(newFile.ID, s.db)
		s.Assert().NoError(err)
		s.Assert().Equal(newFile.DirectoryID, gotDir.ID)
	})

	s.Run("Non-existing file ID", func() {
		_, err := GetDirectoryForFileID(1234, s.db)
		s.Assert().Error(err)
	})
}

func (s *directoryTestSuite) TestGetHomeDir() {
	defer s.down(s.T())

	s.Run("Existing user home dir", func() {
		gotHomeDir, err := GetHomeDir(s.user.ID, s.db)
		s.Assert().NoError(err)
		gotHomeDir.UpdatedAt = s.homeDir.UpdatedAt
		gotHomeDir.CreatedAt = s.homeDir.CreatedAt
		s.Assert().EqualValues(s.homeDir, gotHomeDir)
	})

	s.Run("Non-existing user home dir", func() {
		_, err := GetHomeDir(1234, s.db)
		s.Assert().Error(err)
	})
}

func (s *directoryTestSuite) TestGetOrCreateDirectory() {
	defer s.down(s.T())

	path, err := BuildFullPath("/some/dir", s.user, s.db)
	s.Assert().NoError(err)

	s.Run("Correct path format", func() {
		gotDir, err := GetOrCreateDirectory(path, s.db)
		s.Assert().NoError(err)
		s.Assert().Equal("dir", gotDir.Name)
		s.Assert().Equal(s.user.ID, gotDir.UserID)
		s.Assert().Equal(path, gotDir.Path)
	})

	s.Run("Bad path format", func() {
		_, err := GetOrCreateDirectory("dir/path", s.db)
		s.Assert().Error(err)
	})
}

func (s *directoryTestSuite) TestGetOrCreateHomeDir() {
	defer s.down(s.T())

	var (
		gotHomeDir *schema.Directory
		err        error
	)

	s.Run("Existing home dir", func() {
		gotHomeDir, err = GetOrCreateHomeDir(s.user, s.db)
		s.Assert().NoError(err)
		gotHomeDir.UpdatedAt = s.homeDir.UpdatedAt
		gotHomeDir.CreatedAt = s.homeDir.CreatedAt
		s.Assert().EqualValues(s.homeDir, gotHomeDir)
		isDir, _ := files.Afs.IsDir(gotHomeDir.Path)
		s.Assert().True(isDir)
	})

	s.Run("Non-existing home dir", func() {
		_, err := GetOrCreateHomeDir(&schema.User{}, s.db)
		s.Assert().Error(err)
	})

	s.Run("Creates new dir", func() {
		isDir, _ := files.Afs.IsDir(gotHomeDir.Path)
		s.Assert().True(isDir)
		err = files.Afs.RemoveAll(gotHomeDir.Path)
		s.Assert().NoError(err)

		_, err = GetOrCreateHomeDir(s.user, s.db)
		s.Assert().NoError(err)

		isDir, _ = files.Afs.IsDir(gotHomeDir.Path)
		s.Assert().True(isDir)
	})
}

func (s *directoryTestSuite) TestMakeUserHomeDir() {
	defer s.down(s.T())

	shouldExist := filepath.Join(c.Config.Server.StorageDir, GenerateUserDirName(s.user))
	isDir, _ := files.Afs.Exists(shouldExist)
	if isDir {
		err := files.Afs.RemoveAll(shouldExist)
		s.Assert().NoError(err)
	}

	gotHomeDir, err := MakeUserHomeDir(s.user)
	s.Assert().NoError(err)
	isDir, _ = files.Afs.Exists(shouldExist)
	s.Assert().True(isDir)
	s.Assert().Equal(shouldExist, gotHomeDir)
}

func (s *directoryTestSuite) TestSetupUserHomeDir() {
	s.down(s.T())

	db, down := database.NewTxForTest(s.T())
	s.db = db
	s.down = down

	defer s.down(s.T())

	shouldExist := filepath.Join(c.Config.Server.StorageDir, GenerateUserDirName(s.user))
	pathHash := hash.PathHash(shouldExist)

	s.Run("Create home dir", func() {
		got, err := SetupUserHomeDir(s.user, s.db)
		s.Assert().NoError(err)
		isDir, _ := files.Afs.Exists(got.Path)
		s.Assert().True(isDir)
		s.Assert().Equal(s.user.ID, got.UserID)
		s.Assert().Equal(shouldExist, got.Path)
		s.Assert().Equal(pathHash, got.PathHash)
		s.Assert().Equal("home", got.Name)
		s.Assert().Nil(got.ParentID)
	})

	s.Run("Create duplicate home dir", func() {
		_, err := SetupUserHomeDir(s.user, s.db)
		s.Assert().Error(err)
	})
}

func (s *directoryTestSuite) TestUpdateDirFileCount() {
	defer s.down(s.T())

	s.Run("Existing dir", func() {
		err := UpdateDirFileCount(s.homeDir.ID, s.db)
		s.Assert().NoError(err)
	})

	s.Run("Non existing dir", func() {
		err := UpdateDirFileCount(1234, s.db)
		s.Assert().Error(err)
	})
}
