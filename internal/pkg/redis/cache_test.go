package redis

import (
	"aed-api-server/internal/pkg/cache"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func init() {
	cache.InitPool(cache.RedisConfig{
		Server:    "redis-star-dev.openviewtech.com:6379",
		Password:  "",
		MaxActive: 10,
		MaxIdle:   3,
	})
}

func TestNewCache(t *testing.T) {
	type Point struct {
		X int
		Y int
	}

	testStore := NewCache("test")
	point := Point{
		X: 100,
		Y: 200,
	}

	err := testStore.Put("point", &point, 0)

	assert.Nil(t, err)

	var newPoint Point
	_, err = testStore.Get("point", &newPoint)
	assert.Nil(t, err)
	assert.Equal(t, newPoint, point)

	err = testStore.Remove("point")
	assert.Nil(t, err)

	existed, err := testStore.Get("point", &newPoint)
	assert.Nil(t, err)
	assert.False(t, existed)

	err = testStore.Put("point", &point, 3*time.Second)
	assert.Nil(t, err)

	newPoint = Point{}
	existed, err = testStore.Get("point", &newPoint)
	assert.Nil(t, err)
	assert.True(t, existed)
	assert.Equal(t, newPoint, point)

	time.Sleep(3 * time.Second)
	newPoint = Point{}
	existed, err = testStore.Get("point", &newPoint)
	assert.Nil(t, err)
	assert.False(t, existed)
	assert.Equal(t, newPoint, Point{})
}
