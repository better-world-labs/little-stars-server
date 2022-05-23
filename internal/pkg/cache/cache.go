package cache

// Cache 缓存，包含基本缓存操作
type Cache interface {
	Put(cacheName string, key string, value interface{}) error
	Get(cacheName string, key string) (interface{}, bool, error)
	Evict(cacheName string, key string) error
}
