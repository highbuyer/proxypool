package redistest

import (
	"fmt"
	"github.com/gomodule/redigo/redis"
	"github.com/highbuyer/proxypool/redis" // 替换为您的实际项目导入路径
	"testing"
)

func TestConnectRedis(t *testing.T) {
	client := NewRedisClient("localhost:6379", "")
	conn, err := client.GetConn()
	if err != nil {
		t.Errorf("Expected nil error, but got %v instead", err)
	}

	defer conn.Close()

	value1 := "test_value1"
	key1 := "test_key1"

	setErr := setValue(conn, key1, value1)
	if setErr != nil {
		t.Errorf("Expected no errors when setting value to key. Got %v instead.", setErr.Error())
		return
	}

	resultVal, err := getValue(conn, key1)
	if resultVal != value1 || err != nil {
		t.Errorf("Expected the same value as was set in Redis for key=%s. Got %v,%v respectively.", key1,
			resultVal, err.Error())
		return
	}

	// 使用带密码的客户端连接到 Redis 服务器进行授权操作和数据读写。
	authClient := NewRedisClient("localhost:6379", "password")
	defer authClient.Close()

	conn2, authErr := authClient.GetConn()
	if authErr != nil {
		t.Errorf("Expected no errors when connecting with authorization. Got %v instead.", authErr.Error())
		return
	}

	defer conn2.Close()

	value2 := "test_value2"
	key2 := "test_key2"

	setAuthError := setValue(conn2, key2, value2)
	if setAuthError != nil {
		t.Errorf("Expected no errors when setting a value with authorization. Got %v instead.", setAuthError.Error())
		return
	}

	authVal, autherr := getValue(conn2, key2)
	if authVal != value2 || autherr != nil {
		//t.Errorf(fmt.Sprintf("expected: %v but got: %v (error: %v)", expectedValue, resultValue, err))
		t.Errorf(fmt.Sprintf("expected: %v but got: %v (error: %v)", value2, authVal, err))
		return
	}
}
func setValue(conn redis.Conn, key string, value string) error {
	_, err := conn.Do("SET", key, value)
	return err
}

func getValue(conn redis.Conn, key string) (string, error) {
	valueBytes, err := redis.Bytes(conn.Do("GET", key))
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%s", valueBytes), nil

}
