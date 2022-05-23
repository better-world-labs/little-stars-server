package test

import (
	"aed-api-server/internal/pkg/cache"
	"github.com/stretchr/testify/assert"
	"sync"
	"testing"
)

func TestCacheManager(t *testing.T) {

	c := cache.NewLocalLRUCache(2)
	c2 := cache.NewLocalLRUCache(5)
	manager := cache.NewManager(c, c2)

	err := manager.Put("a", "k1", 1)
	err = manager.Put("a", "k2", 2)
	err = manager.Put("a", "k3", 3)
	err = manager.Put("a", "k4", 4)
	err = manager.Put("a", "k5", 5)
	assert.Nil(t, err)

	group := sync.WaitGroup{}
	group.Add(3)
	go func() {
		o, exists, err := manager.Get("a", "k1")
		assert.Nil(t, err)
		assert.True(t, exists)
		assert.Equal(t, 1, o)
		group.Done()
	}()

	go func() {
		o, exists, err := manager.Get("a", "k2")
		assert.Nil(t, err)
		assert.True(t, exists)
		assert.Equal(t, 2, o)

		o, exists, err = manager.Get("a", "k2")
		assert.Nil(t, err)
		assert.True(t, exists)
		assert.Equal(t, 2, o)
		group.Done()
		group.Done()
	}()

	go func() {

	}()

	group.Wait()
}
