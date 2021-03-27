package data

import (
	"github.com/gomodule/redigo/redis"
	"gitlab.computing.dcu.ie/collint9/2021-ca400-collint9-coynemt2/src/server/config"
	"os"
	"os/signal"
	"syscall"
)

//go:generate mockery --name=Cache --case underscore
type Cache interface {
	GetAll(key string) (interface{}, error)
	SetAll(key string, value interface{}) (interface{}, error)
	Delete(key string) (interface{}, error)
	Ping() (interface{}, error)
	Execute(commandName string, args ...interface{}) (interface{}, error)
}

type RedisCache struct {
	Pool *redis.Pool
	UploadCache UploadCache
}

func NewRedisCache(redisCfg config.RedisConfig, uploadCache UploadCache) *RedisCache {
	return &RedisCache{
		Pool: NewRedisPool(redisCfg),
		UploadCache: uploadCache,
	}
}

func (c *RedisCache) GetAll(key string) (interface{}, error) {
	return c.Execute("HGETALL", key)
}

func (c *RedisCache) SetAll(key string, value interface{}) (interface{}, error) {
	return c.Execute("HSET", redis.Args{}.Add(key).AddFlat(&value)...)
}

func (c *RedisCache) Delete(key string) (interface{}, error) {
	return c.Execute("DEL", key)
}

func (c *RedisCache) Ping() (interface{}, error) {
	return c.Execute("PING")
}

func (c *RedisCache) Execute(commandName string, args ...interface{}) (interface{}, error) {
	conn := c.Pool.Get()
	res, err := conn.Do(commandName, args...)
	if errClose := conn.Close(); errClose != nil {
		return res, errClose
	}
	return res, err
}

func NewRedisPool(redisCfg config.RedisConfig) *redis.Pool {
	pool := &redis.Pool{
		MaxIdle: 80,
		MaxActive: 12000,
		// Create and configure the connection
		Dial: func() (redis.Conn, error) {
			c, err := redis.Dial(redisCfg.Protocol, redisCfg.Address,
				redis.DialPassword(redisCfg.Password),
				redis.DialDatabase(redisCfg.DB),
			)
			return c, err
		},
	}
	cleanupHook(pool)
	return pool
}

func cleanupHook(pool *redis.Pool) {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	signal.Notify(c, syscall.SIGTERM)
	signal.Notify(c, syscall.SIGKILL)
	go func() {
		<-c
		_ = pool.Close()
		os.Exit(0)
	}()
}