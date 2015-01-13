package clients

import (
	"fmt"
	"github.com/garyburd/redigo/redis"
)

type RedisClient struct {
	conn redis.Conn
}

func NewClient(host, port string) *RedisClient {
	client := &RedisClient{}
	client.conn, _ = redis.Dial("tcp", host+":"+port)
	return client
}

func (this *RedisClient) Close() {
	fmt.Println(this.conn)
	if this.conn != nil {
		this.conn.Close()
	}
}

func (this *RedisClient) Get(key string) (string, error) {
	reply, err := this.conn.Do("get", key)
	if err != nil {
		return "", err
	}
	value := string(reply.([]byte))
	return value, nil
}

func (this *RedisClient) Smembers(key string) ([]string, error) {
	reply, err := this.conn.Do("SMEMBERS", key)
	if err != nil {
		return nil, err
	}

	result := make([]string, 0)
	res := reply.([]interface{})
	for _, v := range res {
		result = append(result, string(v.([]byte)))
	}
	return result, nil
}
