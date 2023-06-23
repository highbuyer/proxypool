package redis

import (
	"github.com/highbuyer/proxypool/pool"
)

type RedisClient interface {
	GetConn() (conn redis.Conn, err error)
	Close()
}

type redisClient struct {
	pool pool.RedisClient // 修改此处为公共模块中的接口类型。
}

func NewRedisClient(addr, password string) RedisClient {
	client := &redisClient{pool.NewRedisPool(addr, password)}
	return client // 返回新构建的客户端对象即可。
}
func (c *redisClient) GetConn() (conn redis.Conn, err error) {
	conn = c.pool.Get()
	return conn, nil
}

func (c *redisClient) Close() {
	c.pool.Close()
}
