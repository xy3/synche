package data_test

import (
	"github.com/gomodule/redigo/redis"
	"github.com/rafaeljusto/redigomock"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"gitlab.computing.dcu.ie/collint9/2021-ca400-collint9-coynemt2/src/server/data"
	"testing"
)

type cacheSuite struct {
	suite.Suite
	pool  *redis.Pool
	mock  *redigomock.Conn
	cache *data.RedisCache
}

func (s *cacheSuite) SetupSuite() {
	conn := redigomock.NewConn()
	s.mock = conn
	s.pool = &redis.Pool{
		MaxIdle: 10,
		Dial:    func() (redis.Conn, error) { return s.mock, nil },
	}
	s.cache = &data.RedisCache{Pool: s.pool}
}

func (s *cacheSuite) TestCache_Delete() {
	tests := []struct {
		name    string
		key     string
		want    interface{}
		wantErr bool
	}{
		{name: "delete a key", key: "key", want: nil, wantErr: false},
		{name: "delete a key that doesn't exist", key: "key", want: nil, wantErr: true},
	}
	for _, tc := range tests {
		s.Run(tc.name, func() {
			require.NoError(s.T(), s.mock.Flush())
			cmd := s.mock.Command("DEL", tc.key).Expect(tc.want).ExpectError(nil)
			if tc.wantErr {
				cmd = s.mock.Command("DEL", tc.key).ExpectError(redis.ErrNil)
			}

			got, err := s.cache.Delete(tc.key)
			require.True(s.T(), cmd.Called)
			require.Equal(s.T(), tc.wantErr, err != nil)
			require.Equal(s.T(), tc.want, got)
			require.NoError(s.T(), s.mock.ExpectationsWereMet())
		})
	}
}

func (s *cacheSuite) TestCache_GetAll() {
	tests := []struct {
		name    string
		key     string
		want    interface{}
		wantErr bool
	}{
		{name: "get a key", key: "key", want: nil, wantErr: false},
		{name: "get a key that doesn't exist", key: "key", want: nil, wantErr: true},
	}
	for _, tc := range tests {
		s.Run(tc.name, func() {
			require.NoError(s.T(), s.mock.Flush())
			cmd := s.mock.Command("HGETALL", tc.key).Expect(tc.want).ExpectError(nil)
			if tc.wantErr {
				cmd = s.mock.Command("HGETALL", tc.key).ExpectError(redis.ErrNil)
			}

			got, err := s.cache.GetAll(tc.key)
			require.True(s.T(), cmd.Called)
			require.Equal(s.T(), tc.wantErr, err != nil)
			require.Equal(s.T(), tc.want, got)
			require.NoError(s.T(), s.mock.ExpectationsWereMet())
		})
	}
}

func (s *cacheSuite) TestCache_Ping() {
	tests := []struct {
		name    string
		want    interface{}
		wantErr bool
	}{
		{name: "redis ping", want: nil, wantErr: false},
	}
	for _, tc := range tests {
		s.Run(tc.name, func() {
			require.NoError(s.T(), s.mock.Flush())
			cmd := s.mock.Command("PING").Expect(tc.want)
			if tc.wantErr {
				cmd.ExpectError(redis.ErrNil)
			} else {
				cmd.ExpectError(nil)
			}

			got, err := s.cache.Ping()
			require.True(s.T(), cmd.Called)
			require.Equal(s.T(), tc.wantErr, err != nil)
			require.Equal(s.T(), tc.want, got)
			require.NoError(s.T(), s.mock.ExpectationsWereMet())
		})
	}
}

func (s *cacheSuite) TestCache_Set() {
	tests := []struct {
		name    string
		key     string
		value   interface{}
		want    interface{}
		wantErr bool
	}{
		{name: "set a value to a key", key: "key", value: "value", want: nil, wantErr: false},
	}
	for _, tc := range tests {
		s.Run(tc.name, func() {
			require.NoError(s.T(), s.mock.Flush())
			cmd := s.mock.GenericCommand("HSET").Expect(nil).ExpectError(nil)
			got, err := s.cache.SetAll(tc.key, tc.value)

			require.True(s.T(), cmd.Called)
			require.Equal(s.T(), tc.wantErr, err != nil)
			require.Equal(s.T(), tc.want, got)
			require.NoError(s.T(), s.mock.ExpectationsWereMet())
		})
	}
}

func TestCacheSuite(t *testing.T) {
	suite.Run(t, new(cacheSuite))
}
