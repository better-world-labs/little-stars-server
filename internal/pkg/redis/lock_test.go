package redis

import (
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestNewLocker(t *testing.T) {
	testLocker := locker{
		prefix:      "test",
		lockTime:    4 * time.Second,
		checkPeriod: 2 * time.Second,
	}

	t.Run("LockFn", func(t *testing.T) {
		key := "until-finish"
		err := testLocker.LockFn(key, func() {
			time.Sleep(5 * time.Second)
			_, err := testLocker.Lock(key, 10*time.Second)
			assert.NotNil(t, err)
		})
		//log.Info("unlocked ok?")

		assert.Nil(t, err)
		unlock, err := testLocker.Lock(key, 10*time.Second)
		assert.Nil(t, err)
		defer unlock()
		//time.Sleep(20 * time.Second)
	})

	t.Run("Lock & unlock by expired", func(t *testing.T) {
		key := "lock-test"
		unlock, err := testLocker.Lock(key, 2*time.Second)
		assert.Nil(t, err)
		defer unlock()
		_, err = testLocker.Lock(key, 2*time.Second)
		assert.NotNil(t, err)
		time.Sleep(2 * time.Second)
		unlock2, err := testLocker.Lock(key, 2*time.Second)
		assert.Nil(t, err)
		unlock2()
	})

	t.Run("Lock & call unlock", func(t *testing.T) {
		key := "lock-test2"
		unlock, err := testLocker.Lock(key, 20*time.Second)
		assert.Nil(t, err)

		//未解锁前，不能再次加锁
		_, err = testLocker.Lock(key, 20*time.Second)
		assert.NotNil(t, err)

		//解锁
		unlock()

		//解锁后，允许再次加锁
		unlock, err = testLocker.Lock(key, 20*time.Second)
		assert.Nil(t, err)
		unlock()
	})
}
