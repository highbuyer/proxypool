// pool/pool.go

package pool

import (
	"github.com/gomodule/redigo/redis"
	"log"
	"time"
)

type RedisClient interface {
	GetConn() (conn redis.Conn, err error)
	Close()
}

func NewRedisPool(addr, password string) *redis.Pool {
	return &redis.Pool{
		Dial: func() (redis.Conn, error) {
			return newConnection(addr, password)
		},
		TestOnBorrow: func(c redis.Conn, t time.Time) error {
			_, err := c.Do("PING")
			return err
		},
		MaxIdle:     10,
		MaxActive:   100,
		IdleTimeout: (60 * time.Second),
	}
}

// 内部实现细节，不需要对外暴露
func newConnection(addr, password string) (conn redis.Conn, err error) {
	c, err := redigo.Dial("tcp", addr)
	if err != nil {
		log.Fatal(err)
		return nil, err
	}

	if len(password) > 0 {
		_, authErr := c.Do("AUTH", password)
		if authErr != nil {
			log.Fatal(authErr)
			return nil, authErr
		}
	}

	return c, nil
}
