// internal/redis/redis_test.go

package redis

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSetAndGet(t *testing.T) {
	client := NewRedisClient("127.0.0.1:6379", "")
	conn, err := client.GetConn()
	assert.Nil(t, err)

	key := "test_key"
	value := "test_value"

	_, err = conn.Do("SET", key, value)
	assert.Nil(t, err)

	reply, _ := redis.String(conn.Do("GET", key))
	assert.Equal(t, value, reply)

	client.Close()
}
