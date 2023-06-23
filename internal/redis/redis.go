package redis

import (
	"github.com/gomodule/redigo/redis"
	"log"
	"time"
)

var pool *redis.Pool

type RedisClient interface {
	GetConn() (conn redis.Conn, err error)
	Close()
}

type redisClient struct {
	pool *redis.Pool
}

func NewRedisClient(addr, password string) RedisClient {
	pool := &redis.Pool{
		Dial: func() (redis.Conn, error) {
			return newConnection(addr, password)
		},
		TestOnBorrow: func(c redis.Conn, t time.Time) error {
			_, err := c.Do("PING")
			return err
		},
		// 这里根据需要进行配置...
	}

	client := &redisClient{pool}
	return client
}

func (c *redisClient) GetConn() (conn redis.Conn, err error) {
	conn = c.pool.Get()

	if _, err = conn.Do("PING"); err != nil {
		conn.Close()
		return nil, err
	}

	return conn, nil
}

func (c *redisClient) Close() {
	c.pool.Close()
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

// 内部实现细节，不需要对外暴露
func init() {
	addr := "localhost:6379" // 修改为实际地址和端口号
	password := ""           // 如果有密码，请在这里输入密码

	pool = &redis.Pool{
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
