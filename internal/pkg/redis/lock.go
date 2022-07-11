package redis

import (
	"aed-api-server/internal/interfaces/facility"
	"aed-api-server/internal/pkg/cache"
	"context"
	"fmt"
	"github.com/gomodule/redigo/redis"
	log "github.com/sirupsen/logrus"
	"time"
)

type LockError struct{}

func (*LockError) Error() string {
	return "not lock success"
}

func NewLocker(prefix string) facility.Locker {
	return &locker{
		prefix:      prefix,
		lockTime:    10 * time.Second,
		checkPeriod: 5 * time.Second,
	}
}

type locker struct {
	prefix      string
	lockTime    time.Duration
	checkPeriod time.Duration
}

func (l *locker) buildKey(key string) string {
	return fmt.Sprintf("%s:%s", l.prefix, key)
}
func (l *locker) Lock(key string, ttl time.Duration) (unlock facility.Unlock, err error) {
	key = l.buildKey(key)
	lock, err := cache.GetDistributeLock(key, int64(ttl/time.Millisecond))
	if err != nil {
		return nil, err
	}
	if !lock.Locked() {
		return nil, &LockError{}
	}
	return func() {
		err2 := lock.Release()
		if err2 != nil {
			log.Error("lock.Release() err:", err2)
		}
	}, nil
}

func (l *locker) LockFn(key string, fn func()) (err error) {
	unlock, err := l.Lock(key, l.lockTime)
	if err != nil {
		return err
	}
	defer unlock()

	cancelCtx, stopWatch := context.WithCancel(context.Background())
	defer stopWatch()

	//监听任务完成，给锁续期
	go func() {
		for {
			//log.Info("---->ok")
			select {
			case <-cancelCtx.Done():
				log.Info("lock watch end")
				return

			default:
				time.Sleep(l.checkPeriod)
				err := l.Renewal(key, l.lockTime)
				if err != nil {
					log.Errorf("对 key=%s 续期失败", key)
				}
			}
		}
	}()

	fn()
	return nil
}

func (l *locker) Renewal(key string, ttl time.Duration) error {
	connection := getConnection()
	defer func(connection redis.Conn) {
		err := connection.Close()
		if err != nil {
			log.Error("redis connection.Close() err:", err)
		}
	}(connection)
	key = l.buildKey(key)
	return connection.Send("PEXPIRE", key, int64(ttl/time.Millisecond))
}
