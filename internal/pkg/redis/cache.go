package redis

import (
	"aed-api-server/internal/interfaces/facility"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gomodule/redigo/redis"
	log "github.com/sirupsen/logrus"
	"time"
)

func NewCache(prefix string) facility.Cache {
	return &store{prefix: prefix}
}

type store struct {
	prefix string
}

func (s *store) buildKey(key string) string {
	return fmt.Sprintf("%s:%s", s.prefix, key)
}

func (s *store) Get(key string, v interface{}) (existed bool, err error) {
	connection := getConnection()
	defer func(connection redis.Conn) {
		err := connection.Close()
		if err != nil {
			log.Error("redis connection.Close() err:", err)
		}
	}(connection)

	key = s.buildKey(key)

	reply, err := connection.Do("GET", key)
	if err != nil {
		return false, err
	}
	if reply == nil {
		return false, nil
	}
	//log.Infof("reply:%v", reply)
	s2, ok := reply.([]byte)

	if ok {
		err = json.Unmarshal(s2, v)
		if err != nil {
			return false, err
		}
		return true, nil
	}
	return false, errors.New("unknown err")
}

func (s *store) Put(key string, v interface{}, ttl time.Duration) (err error) {
	connection := getConnection()
	defer func(connection redis.Conn) {
		err := connection.Close()
		if err != nil {
			log.Error("redis connection.Close() err:", err)
		}
	}(connection)
	key = s.buildKey(key)

	bt, err := json.Marshal(v)
	if err != nil {
		return err
	}

	args := make([]interface{}, 0)
	args = append(args, key, bt)
	if ttl > 0 {
		args = append(args, "PX", int64(ttl/time.Millisecond))
	}
	reply, err := connection.Do("SET", args...)
	if err != nil {
		return nil
	}
	if reply != "OK" {
		return errors.New(fmt.Sprintf("err:%v", reply))
	}
	return nil
}

func (s *store) Remove(key string) (err error) {
	connection := getConnection()
	defer func(connection redis.Conn) {
		err := connection.Close()
		if err != nil {
			log.Error("redis connection.Close() err:", err)
		}
	}(connection)
	key = s.buildKey(key)

	return connection.Send("DEL", key)
}
