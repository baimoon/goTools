package clients

import (
	"fmt"
	"github.com/garyburd/redigo/redis"
	"time"
)

type RedisClient struct {
	pool *redis.Pool
}

func NewRedisClient() (client *RedisClient) {
	client = &RedisClient{}
	return
}

func (client *RedisClient) Connect(conStr, pwd string) {
	fmt.Println(conStr)
	client.pool = client.newPool(conStr, pwd)
}

func (client *RedisClient) newPool(server, password string) *redis.Pool {
	return &redis.Pool{
		MaxIdle:     3,
		IdleTimeout: 60 * time.Second,
		Dial: func() (redis.Conn, error) {
			fmt.Println(server)
			c, err := redis.Dial("tcp", server)
			if err != nil {
				return nil, err
			}
			if password != "" {
				if _, err := c.Do("AUTH", password); err != nil {
					c.Close()
					return nil, err
				}
			}
			return c, err
		},
		TestOnBorrow: func(c redis.Conn, t time.Time) error {
			pong, err := c.Do("PING")
			fmt.Println(pong)
			return err
		},
	}
}

func (client RedisClient) Smembers(key string) ([]string, error) {
	result := make([]string, 0)
	c := client.pool.Get()
	smembersScript := redis.NewScript(1, `return redis.call('SMEMBERS', KEYS[1])`)
	reply, err := smembersScript.Do(c, key)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	res := reply.([]interface{})
	for _, v := range res {
		result = append(result, string(v.([]byte)))
	}
	return result, err
}

func (client RedisClient) Sadd(key string, members ...string) (int64, error) {
	c := client.pool.Get()
	var coun int64 = 0
	for _, v := range members {
		saddScript := redis.NewScript(1, `return redis.call('SADD', KEYS[1], ARGV[1])`)
		reply, err := saddScript.Do(c, key, v)
		if err != nil {
			return coun, err
		}
		coun = coun + reply.(int64)
	}
	return coun, nil
}

func (client RedisClient) Srem(key string, members ...string) (int64, error) {
	c := client.pool.Get()
	var count int64 = 0
	for _, v := range members {
		sremScript := redis.NewScript(1, `return redis.call('SREM', KEYS[1], ARGV[1])`)
		reply, err := sremScript.Do(c, key, v)
		if err != nil {
			return count, err
		}
		count = count + reply.(int64)
	}
	return count, nil
}

func (client RedisClient) Hgetall(key string) (map[string]string, error) {
	keys, err := client.Hkeys(key)
	if err != nil {
		return nil, err
	}
	result := make(map[string]string, 0)
	for _, k := range keys {
		v, e1 := client.Hget(key, k)
		if e1 != nil {
			continue
		}
		result[k] = v
	}
	return result, nil
}

func (client RedisClient) Hkeys(key string) ([]string, error) {
	c := client.pool.Get()
	set := make([]string, 0)
	hkeysScript := redis.NewScript(1, `return redis.call('HKEYS', KEYS[1])`)
	reply, err := hkeysScript.Do(c, key)
	if err != nil {
		return nil, err
	}
	list := reply.([]interface{})
	for _, v := range list {
		set = append(set, string(v.([]byte)))
	}
	return set, nil
}

func (client RedisClient) Hget(key, field string) (string, error) {
	c := client.pool.Get()
	hgetScript := redis.NewScript(1, `return redis.call('HGET', KEYS[1], ARGV[1])`)
	reply, err := hgetScript.Do(c, key, field)
	if err != nil {
		return "", err
	}
	return string(reply.([]byte)), nil
}
