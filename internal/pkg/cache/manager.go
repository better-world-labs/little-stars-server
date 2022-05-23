package cache

import (
	"aed-api-server/internal/pkg/base"
	"fmt"
	"log"
	"sync"
)

var once sync.Once
var instance Cache

func GetManager() Cache {
	once.Do(func() {
		instance = NewManager(NewLocalLRUCache(1024 * 10))
	})
	return instance
}

// Manager 缓存行为编排，使用链式结构灵活支持多级缓存
type Manager struct {

	// multiLevelChain 多级缓存链，缓存优先级由前到后
	multiLevelChain []Cache
}

// NewManager 由以先后顺序作为缓存命中优先级的多个 Cache 初始化
func NewManager(cache ...Cache) Cache {
	return &Manager{
		multiLevelChain: cache,
	}
}

func (m *Manager) Put(cacheName string, key string, value interface{}) error {
	lowestLevel := m.multiLevelChain[len(m.multiLevelChain)-1]

	err := lowestLevel.Put(cacheName, key, value)
	if err != nil {
		return base.WrapError("CacheManager", "Put cache error", err)
	}

	return nil
}

func (m *Manager) Get(cacheName string, key string) (interface{}, bool, error) {
	for i := 0; i < len(m.multiLevelChain); i++ {
		c := m.multiLevelChain[i]

		item, exists, err := c.Get(cacheName, key)
		if err != nil || !exists {
			log.Printf("%v", base.WrapError("CacheManager", fmt.Sprintf("not found from cache %d, cacheName=%s, key=%s", i, cacheName, key), err))
			continue
		}

		for j := i - 1; j >= 0; j-- {
			err := m.multiLevelChain[j].Put(cacheName, key, item)
			if err != nil {
				log.Printf("%v", base.WrapError("CacheManager", "put higher cache error", err))
			}
		}

		return item, true, nil
	}

	return nil, false, nil
}

func (m *Manager) Evict(cacheName string, key string) error {
	for _, c := range m.multiLevelChain {
		err := c.Evict(cacheName, key)
		if err != nil {
			return base.WrapError("CacheManager", "cache evict error", err)
		}
	}

	return nil
}
