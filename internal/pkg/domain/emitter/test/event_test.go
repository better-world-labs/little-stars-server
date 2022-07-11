package test

import (
	"aed-api-server/internal/pkg/domain/config"
	"aed-api-server/internal/pkg/domain/emitter"
	"context"
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/atomic"
	"log"
	"sync"
	"testing"
	"time"
)

func TestEmitter(t *testing.T) {
	var tick = TimeTick{
		Id:   "1",
		Tick: 11,
		Time: time.Now(),
	}
	var executeFlag1, executeFlag2, executeFlag3 atomic.Int32

	// 值接收？指针接收？减少订阅方心智负担
	var evtHandler = func(event emitter.DomainEvent) error {
		//g.Done()
		executeFlag1.Inc()
		timeTick := event.(*TimeTick)
		assert.Equal(t, timeTick.Tick, tick.Tick)
		log.Printf("handler handleMessage")
		return nil
	}

	var evtHandler2 = func(event emitter.DomainEvent) error {
		//g.Done()
		executeFlag2.Inc()
		timeTick := event.(*TimeTick)
		assert.Equal(t, timeTick.Tick, tick.Tick)
		log.Printf("handler 2handleMessage")
		return nil
	}

	var evtHandler3 = func(event emitter.DomainEvent) error {
		//g.Done()
		executeFlag3.Inc()
		timeTick := event.(*TimeTick)
		assert.Equal(t, timeTick.Tick, tick.Tick)
		log.Printf(" handler3 handleMessage")

		//return errors.New("xxx")
		return nil
	}

	ctx, cancel := context.WithCancel(context.Background())
	emitter.SetContext(ctx)
	emitter.SetConfig(&config.DomainEventConfig{
		Server:          "localhost:9092",
		Topic:           "star-local-domain-emitter",
		GroupId:         "service-name",
		DeadLetterTopic: "star-local-dead-letter",
	})

	//订阅事件
	emitter.On(&TimeTick{}, evtHandler)
	emitter.On(&TimeTick{}, evtHandler2)
	emitter.On(&TimeTick{}, evtHandler3)

	//发送事件
	go func() {
		for i := 0; i < 1000; i++ {
			err := emitter.Emit(&tick)
			require.Nil(t, err)
		}
	}()

	g := sync.WaitGroup{}
	g.Add(1)
	go func() {
		emitter.Start()
		g.Done()
	}()

	//g.Wait()

	time.Sleep(10 * time.Second)

	fmt.Printf("executeFlag1 = %d\n", executeFlag1)
	fmt.Printf("executeFlag2 = %d\n", executeFlag2)
	fmt.Printf("executeFlag3 = %d\n", executeFlag3)

	assert.Equal(t, int32(1000), executeFlag1.Load())
	assert.Equal(t, int32(1000), executeFlag2.Load())
	assert.Equal(t, int32(1000), executeFlag3.Load())

	emitter.Off(&TimeTick{}, evtHandler)

	cancel()
	g.Wait()
	emitter.Stop()
}
