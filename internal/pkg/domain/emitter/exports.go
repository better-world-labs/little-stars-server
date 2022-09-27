package emitter

import (
	"aed-api-server/internal/pkg/domain/config"
	"github.com/magiconair/properties"
	log "github.com/sirupsen/logrus"
	"sync"
)

type NewEmitterFunc func(conf interface{}) (Emitter, error)

var (
	emitter       Emitter
	onceEmitter   = sync.Once{}
	confMap       = make(map[string]interface{})
	newEmitterMap = make(map[string]NewEmitterFunc)
)

func init() {
	confMap["domain.event.kafka"] = &config.KafkaConfig{}
	confMap["domain.event.rocket"] = &config.RocketConf{}

	newEmitterMap["domain.event.kafka"] = NewKafkaEmitter
	newEmitterMap["domain.event.rocket"] = NewRocketEmitter
}

func SetConfig(p *properties.Properties) {
	onceEmitter.Do(func() {
		for k := range confMap {
			pro := p.FilterStripPrefix(k + ".")
			err := pro.Decode(confMap[k])
			if err == nil {
				fn := newEmitterMap[k]
				e, err := fn(confMap[k])
				if err != nil {
					panic("init emitter err:" + err.Error())
				}
				emitter = e
				return
			} else {
				log.Warnf("cannot init emitter with config:%v", err)
			}
		}
		panic("do not init emitter with config")
	})
}

func SetKafkaDirectly(kafkaConfig *config.KafkaConfig) {
	onceEmitter.Do(func() {
		kafkaEmitter, err := NewKafkaEmitter(kafkaConfig)
		if err != nil {
			panic("init emitter err:" + err.Error())
		}
		emitter = kafkaEmitter
	})
}

func SetRocketDirectly(conf config.RocketConf) {
	onceEmitter.Do(func() {
		kafkaEmitter, err := NewRocketEmitter(conf)
		if err != nil {
			panic("init emitter err:" + err.Error())
		}
		emitter = kafkaEmitter
	})
}

// Start 启动Emitter
func Start() {
	if emitter == nil {
		panic("must call SetConfig First")
	}

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
func On(evt DomainEvent, handlers ...DomainEventHandler) Emitter {
	return emitter.On(evt, handlers...)
}

// Off 撤销事件监听
// @Param evt 空结构体对象，注意不是指针
// @Param handlers 处理函数
func Off(evt DomainEvent, handlers ...DomainEventHandler) Emitter {
	if emitter == nil {
		panic("Must start Emitter first")
	}
	return emitter.Off(evt, handlers...)
}

func Stop() {
	if emitter == nil {
		panic("Must start Emitter first")
	}
	emitter.Close()
}
