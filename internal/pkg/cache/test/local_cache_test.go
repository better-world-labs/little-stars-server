package test

import (
	"aed-api-server/internal/pkg/cache"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestLocalCache(t *testing.T) {
	cache := cache.NewLocalLRUCache(2)
	err := cache.Put("a", "k1", 1)
	assert.Nil(t, err)

	v1, ok, err := cache.Get("a", "k1")
	assert.True(t, ok)
	assert.Equal(t, 1, v1)

	err = cache.Put("a", "k2", 2)
	assert.Nil(t, err)

	v2, ok, err := cache.Get("a", "k2")
	assert.True(t, ok)
	assert.Equal(t, 2, v2)

	err = cache.Put("a", "k3", 3)
	assert.Nil(t, err)

	v3, ok, err := cache.Get("a", "k3")
	assert.True(t, ok)
	assert.Equal(t, 3, v3)

	v1, ok, err = cache.Get("a", "k1")
	assert.False(t, ok)
}
