package repo

import (
	"github.com/patrickmn/go-cache"
	"github.com/spf13/afero"
	"github.com/stretchr/testify/suite"
	"gitlab.computing.dcu.ie/collint9/2021-ca400-collint9-coynemt2/src/database/schema"
	"gitlab.computing.dcu.ie/collint9/2021-ca400-collint9-coynemt2/src/files"
	"gorm.io/gorm"
	"strconv"
	"testing"
)

type chunkTestSuite struct {
	suite.Suite
	user    *schema.User
	homeDir *schema.Directory
	down    func(t *testing.T)
	db      *gorm.DB
}

func Test_chunkTestSuite(t *testing.T) {
	files.SetFileSystem(afero.NewMemMapFs())
	suite.Run(t, new(chunkTestSuite))
}

func (s *chunkTestSuite) SetupTest() {
	user, homeDir, db, down, err := NewUserForTest(s.T())
	s.Assert().NoError(err)
	s.down = down
	s.homeDir = homeDir
	s.user = user
	s.db = db
}

func (s *chunkTestSuite) TestGetFileChunksInOrder() {
	defer s.down(s.T())

	testFile := &schema.File{
		Name:           "testfile.txt",
		Size:           2000,
		Hash:           "fakehash",
		DirectoryID:    s.homeDir.ID,
		UserID:         s.user.ID,
		Available:      true,
		TotalChunks:    10,
		ChunksReceived: 10,
	}

	tx := s.db.Create(testFile)
	s.Assert().NoError(tx.Error)

	s.Run("Existing file with no chunks", func() {
		gotChunks, err := GetFileChunksInOrder(testFile.ID, s.db)
		s.Assert().NoError(err)
		s.Assert().Empty(gotChunks)
	})

	s.Run("Existing file with chunks", func() {
		// Add some fake chunks
		for i := 1; i <= 10; i++ {
			tx = s.db.Create(&schema.FileChunk{
				Number: int64(i),
				Chunk: schema.Chunk{
					Hash: "hash_" + strconv.Itoa(i),
					Size: 20,
				},
				FileID: testFile.ID,
			})
			s.Assert().NoError(tx.Error)
		}

		gotChunks, err := GetFileChunksInOrder(testFile.ID, s.db)
		s.Assert().NoError(err)
		s.Assert().Len(gotChunks, 10)
		for i, gotChunk := range gotChunks {
			s.Assert().Equal(int64(i+1), gotChunk.Number)
		}
	})

	s.Run("Non-existing file", func() {
		gotChunks, err := GetFileChunksInOrder(1234, s.db)
		s.Assert().NoError(err)
		s.Assert().Len(gotChunks, 0)
	})
}

func (s *chunkTestSuite) TestGetCachedChunksReceived() {
	defer s.down(s.T())

	s.Run("Non-existing file", func() {
		_, err := GetCachedChunksReceived(1234)
		s.Assert().Error(err)
	})

	s.Run("Existing file", func() {
		FileIDChunkCountCache.Set(strconv.Itoa(1000), uint64(10), cache.DefaultExpiration)
		got, err := GetCachedChunksReceived(1000)
		s.Assert().NoError(err)
		s.Assert().Equal(uint64(10), got)
		FileIDChunkCountCache.Delete(strconv.Itoa(1000))
	})
}
