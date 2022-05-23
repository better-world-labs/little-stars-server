package cache

import "github.com/gomodule/redigo/redis"

type RedisConfig struct {
	Server    string `yaml:"server"`
	Password  string `yaml:"password"`
	MaxIdle   int    `yaml:"max-idle"`
	MaxActive int    `yaml:"max-active"`
}

var pool *redis.Pool

func InitPool(c RedisConfig) {
	pool = &redis.Pool{
		MaxIdle:   c.MaxIdle,   /*最大的空闲连接数*/
		MaxActive: c.MaxActive, /*最大的激活连接数*/
		Dial: func() (redis.Conn, error) {
			c, err := redis.Dial("tcp", c.Server, redis.DialPassword(c.Password))
			if err != nil {
				return nil, err
			}
			return c, nil
		},
	}
	_, err := pool.Get().Do("ping")
	if err != nil {
		panic(err)
	}
}

func GetConn() redis.Conn {
	return pool.Get()
}
