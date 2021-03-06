package cache

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func init() {
	InitPool(RedisConfig{Server: "116.62.220.222:6379", Password: "caepheip9Phu9dae", MaxActive: 10, MaxIdle: 3})
}
func TestDistributeLock(t *testing.T) {
	lock, err := GetDistributeLock("lockname", 30000)
	assert.Nil(t, err)
	assert.True(t, lock.Locked())

	lock2, err2 := GetDistributeLock("lockname", 30000)
	assert.Nil(t, err2)
	assert.False(t, lock2.Locked())

	err = lock.Release()
	assert.Nil(t, err)

	lock, err = GetDistributeLock("lockname", 30000)
	assert.Nil(t, err)
	assert.True(t, lock.Locked())

	err = lock.Release()
	assert.Nil(t, err)

}
