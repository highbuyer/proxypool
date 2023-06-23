// redis/client.go

package redis

import (
	"github.com/gomodule/redigo/redis"
	"time"
)

var Pool *redis.Pool

func init() {
	// 可以从环境变量或配置文件读取这些参数
	addr := "localhost:6379"
	password := ""

	// 定义 Dial 函数以创建新连接
	dialFunc := func() (redis.Conn, error) {
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

		return c, err
	}

	// 初始化连接池对象 Pool
	Pool = &redis.Pool{
		Dial: dialFunc,
		TestOnBorrow: func(c redis.Conn, t time.Time) error {
			_, err := c.Do("PING")
			return err
		},
	}
}

// ConnectRedis 使用指定的参数连接 Redis 服务器，并返回与之对应的 Connection 实例。
func ConnectRedis(addr string, password string) (redis.Conn, error) {
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

// 封装一些常用操作函数，如 Set、Get 等等。
