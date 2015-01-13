package clients

import (
	"github.com/garyburd/redigo/redis"
)

type RedisClient struct {
	pool *redis.Pool
}

func NewRedisClient(newFn func() (redis.Conn, error), max int) *RedisClient {
	client := &RedisClient{}

	client.pool = redis.NewPool(newFn, max)
	return client
}

//设置字符串类型的值
func (this *RedisClient) Get(key string) (string, error) {
	conn := this.pool.Get()
	reply, err := conn.Do("get", key)
	if err != nil || reply == nil {
		return "", err
	}
	value := string(reply.([]byte))
	conn.Flush()
	conn.Close()
	return value, nil
}

//获取字符串类型的值
func (this *RedisClient) Set(key, value string) error {
	conn := this.pool.Get()
	_, err := conn.Do("SET", key, value)
	if err != nil {
		return err
	}
	conn.Flush()
	conn.Close()
	return nil
}

//获取set的所有成员
func (this *RedisClient) Smembers(key string) ([]string, error) {
	conn := this.pool.Get()
	reply, err := conn.Do("SMEMBERS", key)
	if err != nil {
		return nil, err
	}

	result := make([]string, 0)
	if reply != nil {
		res := reply.([]interface{})
		for _, v := range res {
			result = append(result, string(v.([]byte)))
		}
	}
	conn.Flush()
	conn.Close()
	return result, nil
}

//向set中添加信息
func (this *RedisClient) Sadd(key string, values ...string) (int64, error) {
	conn := this.pool.Get()
	var count int64 = 0
	for _, v := range values {
		reply, err := conn.Do("SADD", key, v)
		if err != nil {
			continue
		}
		count = count + reply.(int64)
	}

	conn.Flush()
	conn.Close()
	return count, nil
}

//从set中删除数据
func (this *RedisClient) Srem(key string, values ...string) (int64, error) {
	conn := this.pool.Get()
	var count int64 = 0
	for _, v := range values {
		reply, err := conn.Do("SREM", key, v)
		if err != nil {
			continue
		}
		count = count + reply.(int64)
	}

	conn.Flush()
	conn.Close()
	return count, nil
}

func (this *RedisClient) Hgetall(key string) (map[string]string, error) {
	conn := this.pool.Get()
	reply, err := this._hkeys(conn, key)
	if err != nil {
		return nil, err
	}

	result := make(map[string]string, 0)
	for _, f := range reply {
		v, e1 := this._hget(conn, key, f)
		if e1 != nil {
			continue
		}
		result[f] = v
	}
	conn.Flush()
	conn.Close()
	return result, nil
}

func (this *RedisClient) Hget(key string, field string) (string, error) {
	conn := this.pool.Get()
	reply, err := this._hget(conn, key, field)
	conn.Flush()
	conn.Close()
	return reply, err
}

func (this *RedisClient) _hget(conn redis.Conn, key string, field string) (string, error) {
	reply, err := conn.Do("HGET", key, field)
	if err != nil || reply == nil {
		return "", err
	}
	return string(reply.([]byte)), err
}

//获取hash的keys
func (this *RedisClient) Hkeys(key string) ([]string, error) {
	conn := this.pool.Get()
	reply, err := this._hkeys(conn, key)
	if err != nil {
		return nil, err
	}
	conn.Flush()
	conn.Close()
	return reply, nil
}

func (this *RedisClient) _hkeys(conn redis.Conn, key string) ([]string, error) {
	reply, err := conn.Do("HKEYS", key)
	if err != nil {
		return nil, err
	}
	result := make([]string, 0)
	if reply != nil {
		res := reply.([]interface{})
		for _, v := range res {
			result = append(result, string(v.([]byte)))
		}
	}
	return result, nil
}

//获取hash的values
func (this *RedisClient) Hvals(key string) ([]string, error) {
	conn := this.pool.Get()
	reply, err := conn.Do("HVALS", key)
	if err != nil {
		return nil, err
	}
	result := make([]string, 0)
	if reply != nil {
		res := reply.([]interface{})
		for _, v := range res {
			result = append(result, string(v.([]byte)))
		}
	}
	conn.Flush()
	conn.Close()
	return result, nil
}
