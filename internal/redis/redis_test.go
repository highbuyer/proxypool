//package redis
//
//import (
//	//"github.com/highbuyer/proxypool/redis" // 导入你自己项目路径下的 redis 包
//	"testing"
//)
//
//func TestRedisConnect(t *testing.T) {
//	conn := redis.Pool.Get()
//	defer conn.Close()
//
//	_, err := conn.Do("PING")
//
//	if err != nil {
//		t.Errorf("cannot connect to Redis: %s", err)
//		return
//	}
//}
// redis_test.go

package redis

import (
	"testing"
)

func TestConnectRedis(t *testing.T) {
	// 测试未授权访问
	conn, err := ConnectRedis("localhost:6379", "")
	if err != nil {
		t.Errorf("Expected nil error, but got %v instead", err)
	}

	defer conn.Close()

	_, setErr := conn.Do("SET", "test_key1", "test_value1")
	if setErr != nil {
		t.Errorf("Expected nil error, but got %v instead", setErr)
	}

	value, err := conn.Do("GET", "test_key1")
	if value != "test_value1" || err != nil {
		t.Errorf("Expected test_value1,nil ,but got %v,%v instead", value, err.Error())
	}

	// 测试授权访问
	conn2, authErr := ConnectRedis("localhost:6379", "password")
	if authErr != nil {
		t.Errorf("Expected nil error, but got %v instead", authErr)
	}

	defer conn2.Close()

	_, setAuthError := conn2.Do("SET", "test_key2", "test_value2")
	if setAuthError != nil {
		t.Errorf("Expected nil error,but got %v instead ", setAuthError.Error())
	}

	val, err := conn2.Do("GET ", " test_key2 ")
	if val != " test_value2 " || err != nil {
		t.Errorf("Expected test_value2,nil,but got %v,%v instead", val, err.Error())
	}
}
