package emitter

import (
	log "github.com/sirupsen/logrus"
	"reflect"
	"sync"
)

// HandlerRegistry  Handler 注册器
// 线程安全，目前锁粒度较粗，但是目前场景基本没什么影响，慢慢优化
// TODO Register 时，若 key 不存在，新创建的 HandlerKeeper CAS 替换到 map 中, 失败则重新获取 key再次尝试
// TODO 若 key 存在，则对 Key 加锁修改 HandlerKeeper 中的 Slice
type HandlerRegistry struct {
	keepers map[string]*HandlerKeeper
	rwMutex sync.RWMutex // 需要优化为 key 粒度的锁而非全局锁
}

func NewHandlerRegistry() *HandlerRegistry {
	return &HandlerRegistry{
		keepers: make(map[string]*HandlerKeeper),
		rwMutex: sync.RWMutex{},
	}
}

func (k *HandlerRegistry) Register(key DomainEvent, handler DomainEventHandler) {
	eventType := GetStructType(key)
	log.Infof("Register for EventType %s", eventType)

	k.rwMutex.Lock()
	defer k.rwMutex.Unlock()
	handlerKeeper, ok := k.keepers[eventType]
	if ok {
		k.doAppend(handlerKeeper, handler)
		return
	}

	handlerKeeper, ok = k.keepers[eventType]
	if ok {
		k.doAppend(handlerKeeper, handler)
		return
	}

	keeper := k.createKeeper(key, handler)
	k.keepers[eventType] = keeper
}

func (k *HandlerRegistry) doAppend(keeper *HandlerKeeper, handler DomainEventHandler) {
	keeper.slice = append(keeper.slice, handler)
}

func (k *HandlerRegistry) createKeeper(decoder Decoder, handler DomainEventHandler) *HandlerKeeper {
	return &HandlerKeeper{
		slice:   []DomainEventHandler{handler},
		decoder: decoder,
	}
}

func (k *HandlerRegistry) Get(key string) (*HandlerKeeper, bool) {
	k.rwMutex.RLock()
	defer k.rwMutex.RUnlock()

	keeper, exists := k.keepers[key]
	return keeper, exists
}

func (k *HandlerRegistry) Delete(evt DomainEvent, handler DomainEventHandler) {
	eventType := GetStructType(evt)
	log.Infof("Delete handler for EventType %s", eventType)

	k.rwMutex.Lock()
	defer k.rwMutex.Unlock()

	if keeper, exists := k.keepers[eventType]; exists {
		of1 := reflect.ValueOf(handler)
		for i, h := range keeper.slice {
			of2 := reflect.ValueOf(h)
			if of1.Pointer() == of2.Pointer() {
				keeper.slice = append(keeper.slice[:i], keeper.slice[i+1:]...)
				break
			}
		}

		if len(keeper.slice) == 0 {
			delete(k.keepers, eventType)
		}
	}
}

func (k *HandlerRegistry) GetEventTypes() []string {
	types := make([]string, 0, len(k.keepers))
	for t := range k.keepers {
		types = append(types, t)
	}
	return types
}
