package redis

import (
	"aed-api-server/internal/pkg/cache"
	"github.com/gomodule/redigo/redis"
)

func getConnection() redis.Conn {
	return cache.GetConn()
}
