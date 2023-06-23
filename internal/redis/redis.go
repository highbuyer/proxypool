// Package redis.go
package redis

import (
	"github.com/gomodule/redigo/redis"
	"time"
)

var Pool *redis.Pool // Redis 连接池对象

// 创建新 Redis 连接
func newConnection(addr, password string) (redis.Conn, error) {
	c, err := redis.Dial("tcp", addr)
	if err != nil {
		return nil, err
	}

	if len(password) > 0 {
		_, authErr := c.Do("AUTH", password)
		if authErr != nil {
			c.Close()
			return nil, authErr
		}
	}

	return c, nil
}

// 初始化 Redis 连接池（仅需执行一次）
func init() {
	addr := "localhost:6379" // 修改为实际地址和端口号
	password := ""           // 如果有密码，请在这里输入密码

	Pool = &redis.Pool{
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

// ConnectRedis 返回已连接到指定 Redis 服务器的连接对象。
func ConnectRedis(addr, password string) (conn redis.Conn, err error) {
	conn = Pool.Get()

	_, err = conn.Do("PING")
	if err != nil {
		conn.Close()
		return nil, err
	}
	return conn, nil
}
