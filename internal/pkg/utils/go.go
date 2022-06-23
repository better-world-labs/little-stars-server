package utils

import (
	"aed-api-server/internal/pkg/cache"
	"github.com/jtolds/gls"
	log "github.com/sirupsen/logrus"
)

func Go(fn func()) {
	gls.Go(func() {
		defer func() {
			if err := recover(); err != nil {
				log.Errorf("handle panic: %v, %s", err, PanicTrace(2))
			}
		}()
		fn()
	})
}

func GoWrapWithNewTraceId(fn func()) func() {
	return func() {
		go func() {
			defer func() {
				if err := recover(); err != nil {
					log.Errorf("handle panic: %v, %s", err, PanicTrace(2))
				}
			}()
			SetTraceId("", fn)
		}()
	}
}

//JobLockDoWrap 将一个函数fn做包装：1.处理函数可能产生的panic；2.通过redis加锁，加锁成功才执行函数fn
//lockKey 加锁使用的key
//fn 被包装的函数
//expired 锁超时的时间
func JobLockDoWrap(lockKey string, fn func(), expired int64) func() {
	return GoWrapWithNewTraceId(func() {
		lock, err := cache.GetDistributeLock(lockKey, expired)
		if err != nil {
			log.Error("add lock err:", err)
		}

		defer lock.Release()

		if !lock.Locked() {
			log.Info("get " + lockKey + " not suc")
			return
		}

		log.Info("get " + lockKey + " suc")
		fn()
	})
}
