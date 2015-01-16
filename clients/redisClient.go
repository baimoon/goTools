package clients

import (
	"errors"
	"github.com/garyburd/redigo/redis"
	"strconv"
	"time"
)

type RedisClient struct {
	pool *redis.Pool
}

func NewRedisClient(newFn func() (redis.Conn, error), max int) *RedisClient {
	client := &RedisClient{}

	client.pool = redis.NewPool(newFn, max)
	return client
}

func NewRedisClientForWeb(server, password string) *RedisClient {
	client := new(RedisClient)
	client.pool = newPool(server, password)
	return client
}

//采用作者建议的web application的创建线程池的方法
func newPool(server, password string) *redis.Pool {
	return &redis.Pool{
		MaxIdle:     5,
		IdleTimeout: 240 * time.Second,
		Dial: func() (redis.Conn, error) {
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
			_, err := c.Do("PING")
			return err
		},
	}
}

//---------------------------key类的操作-----------------------------
//判断指定的key是否存在
func (this *RedisClient) Exists(key string) (bool, error) {
	conn := this.pool.Get()
	defer conn.Close()
	reply, err := conn.Do("EXISTS", key)
	if err != nil {
		return false, err
	}
	return reply.(int64) == 1, nil
}

//删除指定的key
func (this *RedisClient) Del(key string) (bool, error) {
	conn := this.pool.Get()
	defer conn.Close()
	reply, err := conn.Do("DEL", key)
	if err != nil {
		return false, err
	}
	return reply.(int64) == 1, nil
}

//---------------------------字符串类的操作---------------------------------
//设置字符串类型的值
func (this *RedisClient) Get(key string) (string, error) {
	conn := this.pool.Get()
	defer conn.Close()
	reply, err := conn.Do("get", key)
	if err != nil || reply == nil {
		return "", err
	}
	value := string(reply.([]byte))
	return value, nil
}

//获取字符串类型的值
func (this *RedisClient) Set(key, value string) error {
	conn := this.pool.Get()
	defer conn.Close()
	_, err := conn.Do("SET", key, value)
	if err != nil {
		return err
	}
	return nil
}

//-------------------------SET类的操作--------------------------
//获取set的所有成员
func (this *RedisClient) Smembers(key string) ([]string, error) {
	conn := this.pool.Get()
	defer conn.Close()
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
	return result, nil
}

//向set中添加信息
func (this *RedisClient) Sadd(key string, values ...string) (int64, error) {
	conn := this.pool.Get()
	defer conn.Close()
	var count int64 = 0
	for _, v := range values {
		reply, err := conn.Do("SADD", key, v)
		if err != nil {
			continue
		}
		count = count + reply.(int64)
	}
	return count, nil
}

//从set中删除数据
func (this *RedisClient) Srem(key string, values ...string) (int64, error) {
	conn := this.pool.Get()
	defer conn.Close()
	var count int64 = 0
	for _, v := range values {
		reply, err := conn.Do("SREM", key, v)
		if err != nil {
			continue
		}
		count = count + reply.(int64)
	}
	return count, nil
}

//set中是否存在
func (this *RedisClient) Sismembers(key string, field string) (bool, error) {
	conn := this.pool.Get()
	defer conn.Close()

	reply, err := conn.Do("SISMEMBER", key, field)
	if err != nil {
		return true, err
	}
	return reply.(int64) == 1, nil
}

//set中元素的个数
func (this *RedisClient) Scard(key string) (int64, error) {
	conn := this.pool.Get()
	defer conn.Close()

	reply, err := conn.Do("SCARD", key)
	if err != nil {
		return -1, err
	}
	return reply.(int64), nil
}

//------------------------HASH类的操作-------------------------
func (this *RedisClient) Hgetall(key string) (map[string]string, error) {
	conn := this.pool.Get()
	defer conn.Close()
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
	return result, nil
}

//从hash中获取指定field的数据
func (this *RedisClient) Hget(key string, field string) (string, error) {
	conn := this.pool.Get()
	defer conn.Close()
	reply, err := this._hget(conn, key, field)
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
	defer conn.Close()
	reply, err := this._hkeys(conn, key)
	if err != nil {
		return nil, err
	}
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
	defer conn.Close()
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
	return result, nil
}

func (this *RedisClient) _hset(conn redis.Conn, key string, field string, value string) (int64, error) {
	reply, err := conn.Do("HSET", key, field, value)
	if err != nil {
		return 0, err
	}
	return reply.(int64), nil
}

//向hash中添加数据
func (this *RedisClient) Hset(key string, field string, value string) (int64, error) {
	conn := this.pool.Get()
	defer conn.Close()
	return this._hset(conn, key, field, value)
}

//向hash中添加指定的数据
func (this *RedisClient) Hmset(key string, values map[string]string) (int64, error) {
	conn := this.pool.Get()
	defer conn.Close()

	var count int64 = 0
	for field, value := range values {
		reply, err := this._hset(conn, key, field, value)
		if err != nil {
			continue
		}
		count = count + reply
	}
	return count, nil
}

//---------------------zset---------------------
//获取zset中的分数
func (this *RedisClient) Zscore(key string, member string) (float64, error) {
	conn := this.pool.Get()
	defer conn.Close()

	reply, err := conn.Do("ZSCORE", key, member)
	if err != nil {
		return -1, err
	} else if reply == nil {
		return -1, errors.New("没有对应的数据")
	}

	var score float64
	scores := reply.([]uint8)

	if len(scores) > 0 {
		scoreStr := string(scores)
		score, err = strconv.ParseFloat(scoreStr, 64)
		if err != nil {
			return -1, nil
		}
	} else {
		return -1, errors.New("没有对应的数据")
	}

	return float64(score), nil
}
