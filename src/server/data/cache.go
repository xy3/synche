package data

import (
	"github.com/gomodule/redigo/redis"
	log "github.com/sirupsen/logrus"
	"gitlab.computing.dcu.ie/collint9/2021-ca400-collint9-coynemt2/src/server/config"
)

type Cache interface {
	SetNumberOfChunks(uploadRequestId string, numberOfChunks int64) error
	GetNumberOfChunks(uploadRequestId string) (numberOfChunks int64, err error)
	DeleteKeyNumberOfChunks(uploadRequestId string) error
}

type CacheData struct {
	// wrap redis connection and any other driver that may be needed
	redis *redis.Pool
}

func ping(c redis.Conn) error {
	s, err := redis.String(c.Do("PING"))
	if err != nil {
		return err
	}

	// Log PONG
	log.Infof("CacheData PING response = %s\n", s)
	return nil
}

func newPool(network string, address string, port string, password string, db int) *redis.Pool {
	return &redis.Pool{
		MaxIdle: 80,
		MaxActive: 12000,
		// Create and configure the connection
		Dial: func() (redis.Conn, error) {
			c, err := redis.Dial(network, address + ":" + port, redis.DialPassword(password), redis.DialDatabase(db))
			return c, err
		},
	}
}

func BuildRedisClient(rConfig config.RedisConfig) *CacheData {
	// Pointer to redis.Pool
	pool := newPool(rConfig.Network, rConfig.Address, rConfig.Port, rConfig.Password, rConfig.DB)

	// Connection from the pool
	conn := pool.Get()
	defer conn.Close()

	// Test connectivity
	err := ping(conn)
	if err != nil {
		log.Errorf("Issue connecting to Redis: %v", err)
	}

	return &CacheData{redis: pool}
}
