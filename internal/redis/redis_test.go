package redis

import (
	"github.com/highbuyer/proxypool/redis" // 导入你自己项目路径下的 redis 包
	"testing"
)

func TestRedisConnect(t *testing.T) {
	conn := redis.Pool.Get()
	defer conn.Close()

	_, err := conn.Do("PING")

	if err != nil {
		t.Errorf("cannot connect to Redis: %s", err)
		return
	}
}
