package cache

import (
	lru "github.com/hashicorp/golang-lru"
	"go.uber.org/zap/buffer"
)

type LocalLRUCache struct {
	lru *lru.Cache
}

func NewLocalLRUCache(size int) Cache {
	cache, err := lru.New(size)
	if err != nil {
		panic(err)
	}
	return &LocalLRUCache{
		lru: cache,
	}
}

func (l *LocalLRUCache) Put(cacheName string, key string, value interface{}) error {
	l.lru.Add(wrapKey(cacheName, key), value)
	return nil
}

func (l *LocalLRUCache) Get(cacheName string, key string) (interface{}, bool, error) {
	o, ok := l.lru.Get(wrapKey(cacheName, key))
	return o, ok, nil
}

func (l *LocalLRUCache) Evict(cacheName string, key string) error {
	l.lru.Remove(wrapKey(cacheName, key))
	return nil
}

func wrapKey(cacheName string, key string) string {
	b := buffer.Buffer{}
	b.AppendString(cacheName)
	b.AppendString("-")
	b.AppendString(key)
	return b.String()
}
