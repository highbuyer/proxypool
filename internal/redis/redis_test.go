package redis

import (
	"fmt"
	"github.com/gomodule/redigo/redis"
	"testing"
)

func TestConnectRedis(t *testing.T) {
	conn, err := ConnectRedis("localhost:6379", "")
	if err != nil {
		t.Errorf("Expected nil error, but got %v instead", err)
	}

	defer conn.Close()

	value1 := "test_value1"
	key1 := "test_key1"

	setErr := setValue(conn, key1, value1, t)
	if setErr != nil {
		t.Errorf("Expected no errors when setting value to key. Got %v instead.", setErr.Error())
		return
	}

	resultVal, err := getValue(conn, key1, t)
	if resultVal != value1 || err != nil {
		t.Errorf("Expected the same value as was set in Redis for key=%s. Got %v,%v respectively.", key1,
			resultVal, err.Error())
		return
	}

	conn2, authErr := ConnectRedis("localhost:6379", "password")
	if authErr != nil {
		t.Errorf("Expected no errors when connecting with authorization. Got %v instead.", authErr.Error())
		return
	}

	defer conn2.Close()

	value2 := "test_value2"
	key2 := "test_key2"

	setAuthError := setValue(conn2,key2,value2,t )
	if setAuthError!=nil{
		t.Errorf("Expected no errors when setting a value with authorization. Got %v instead.",setAuthError.Error() )
		return
	}

	authVal ,autherr:=getValue(conn,key,t );
	if(authVal!=value||autherr!=nil){
		t.Errorf(fmt.Sprintf(“expect :%s actual:%s”,value,result,err))
		return
	}

}