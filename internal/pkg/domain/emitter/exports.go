package emitter

import (
	"aed-api-server/internal/pkg/domain/config"
	"context"
	"sync"
)

var (
	conf        *config.DomainEventConfig
	emitter     Emitter
	onceEmitter = sync.Once{}
	ctx         context.Context
)

func checkAndInit() {
	if conf == nil {
		panic("Must set config first")
	}

	if ctx == nil {
		panic("Must set context first")
	}

	onceEmitter.Do(func() {
		e, err := NewKafkaEmitter(ctx, conf)
		if err != nil {
			panic(err)
		}

		emitter = e
	})
}

func SetConfig(c *config.DomainEventConfig) {
	conf = c
}

func SetContext(c context.Context) {
	ctx = c
}

// Start 启动Emitter
func Start() {
	checkAndInit()
	emitter.Start()
}

// Emit 发送事件
func Emit(events ...DomainEvent) error {
	if emitter == nil {
		panic("Must start Emitter first")
	}
	return emitter.Emit(events...)
}

// On 开启事件监听
// @Param evt 空结构体对象，注意不是指针
// @Param handlers 处理函数
func On(evt DomainEvent, handlers ...DomainEventHandler) {
	checkAndInit()
	emitter.On(evt, handlers...)
}

// Off 撤销事件监听
// @Param evt 空结构体对象，注意不是指针
// @Param handlers 处理函数
func Off(evt DomainEvent, handlers ...DomainEventHandler) {
	if emitter == nil {
		panic("Must start Emitter first")
	}
	checkAndInit()
	emitter.Off(evt, handlers...)
}

func Stop() {
	if emitter == nil {
		panic("Must start Emitter first")
	}
	emitter.Close()
}
