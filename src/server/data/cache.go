package data

import (
	"github.com/gomodule/redigo/redis"
	log "github.com/sirupsen/logrus"
	"gitlab.computing.dcu.ie/collint9/2021-ca400-collint9-coynemt2/src/server/config"
)

type Cache struct {
	// wrap redis connection and any other driver that may be needed
	redis *redis.Pool
}

func ping(c redis.Conn) error {
	s, err := redis.String(c.Do("PING"))
	if err != nil {
		return err
	}

	log.Debugf("Cache PING response: %s", s)
	return nil
}

func newPool(redisCfg config.RedisConfig) *redis.Pool {
	return &redis.Pool{
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
}

func NewRedisCache(redisCfg config.RedisConfig) *Cache {
	// Pointer to redis.Pool
	pool := newPool(redisCfg)

	// Connection from the pool
	conn := pool.Get()
	defer conn.Close()

	// Test connectivity
	err := ping(conn)
	if err != nil {
		log.Errorf("Issue connecting to Redis: %v", err)
	}

	return &Cache{redis: pool}
}
