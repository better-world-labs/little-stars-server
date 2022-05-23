package test

import (
	"aed-api-server/internal/pkg/domain/emitter"
	"fmt"
	"github.com/stretchr/testify/require"
	"sync"
	"testing"
	"time"
)

type Event1 struct {
	Id        string        `json:"id"`
	Timestamp time.Duration `json:"timestamp"`
}

func (e *Event1) Decode(bytes []byte) (emitter.DomainEvent, error) {
	panic("implement me")
}

func (e *Event1) Encode() ([]byte, error) {
	panic("implement me")
}

func TestHandlerRegistry(t *testing.T) {
	//fmt.Printf("handleTimeTrick1 = %p", handlerTimeTick1)
	//fmt.Printf("handleTimeTrick2 = %p", handlerTimeTick2)
	//fmt.Printf("handleTimeTrick3 = %p", handlerTimeTick3)
	//fmt.Printf("handleTimeTrick4 = %p", handlerTimeTick4)
	//fmt.Printf("handleTimeTrick5 = %p", handlerTimeTick5)
	//fmt.Printf("handleTimeTrick6 = %p", handlerTimeTick6)
	i := 0
	for ; i < 100000; i++ {
		//fmt.Printf("----------------------------start----------------------------------------------\n")
		registry := emitter.NewHandlerRegistry()
		g := sync.WaitGroup{}
		g.Add(5)

		go func() {
			registry.Register(&TimeTick{}, handlerTimeTick1)
			registry.Register(&TimeTick{}, handlerTimeTick2)
			g.Done()
		}()

		go func() {
			registry.Register(&TimeTick{}, handlerTimeTick3)
			registry.Register(&TimeTick{}, handlerTimeTick4)
			g.Done()
		}()

		go func() {
			registry.Register(&TimeTick{}, handlerTimeTick5)
			registry.Register(&TimeTick{}, handlerTimeTick6)
			g.Done()
		}()

		go func() {
			registry.Register(&Event1{}, handlerEvent11)
			registry.Register(&Event1{}, handlerEvent12)
			g.Done()
		}()

		go func() {
			registry.Register(&Event1{}, handlerEvent12)
			registry.Register(&Event1{}, handlerEvent11)
			g.Done()
		}()
		g.Wait()

		keeper, exists := registry.Get(emitter.GetStructType(TimeTick{}))
		require.True(t, exists)
		if len(keeper.Handlers()) != 6 {
			fmt.Println(keeper)
			break
		}

		require.Equal(t, 6, len(keeper.Handlers()))

		keeper, exists = registry.Get(emitter.GetStructType(TimeTick{}))
		require.True(t, exists)

		g.Add(3)
		go func() {
			registry.Delete(&TimeTick{}, handlerTimeTick1)
			registry.Delete(&TimeTick{}, handlerTimeTick2)
			g.Done()
		}()

		go func() {
			registry.Delete(&TimeTick{}, handlerTimeTick3)
			registry.Delete(&TimeTick{}, handlerTimeTick4)
			g.Done()
		}()

		go func() {
			registry.Delete(&TimeTick{}, handlerTimeTick5)
			registry.Delete(&TimeTick{}, handlerTimeTick6)
			g.Done()
		}()

		g.Wait()
		_, exists = registry.Get(emitter.GetStructType(&TimeTick{}))
		require.False(t, exists)

		g.Add(1)
		go func() {
			registry.Register(&TimeTick{}, handlerTimeTick5)
			registry.Register(&TimeTick{}, handlerTimeTick6)
			g.Done()
		}()

		g.Wait()
		h, exists := registry.Get(emitter.GetStructType(&TimeTick{}))
		require.True(t, exists)
		require.Equal(t, 2, len(h.Handlers()))
	}

	require.Equal(t, 100000, i)
}

func handlerEvent12(event emitter.DomainEvent) error {
	return nil
}

func handlerEvent11(event emitter.DomainEvent) error {
	return nil
}

func handlerTimeTick6(event emitter.DomainEvent) error {
	return nil
}

func handlerTimeTick5(event emitter.DomainEvent) error {
	return nil
}

func handlerTimeTick4(event emitter.DomainEvent) error {
	return nil
}

func handlerTimeTick3(event emitter.DomainEvent) error {
	return nil
}

func handlerTimeTick2(event emitter.DomainEvent) error {
	return nil
}

func handlerTimeTick1(event emitter.DomainEvent) error {

	return nil
}
