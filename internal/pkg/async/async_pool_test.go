package async

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"sync"
	"testing"
)

func TestPool(t *testing.T) {
	client := NewPool(Config{
		GoRoutinePoolSize: 10,
		MaxTaskQueueSize:  10,
	})

	client.Start()

	group := sync.WaitGroup{}
	group.Add(1000)
	for i := 0; i < 1000; i++ {
		err := client.AddSyncTask(func() {
			fmt.Println("hahaha")
			//panic(ErrorCientStoped)
			group.Done()
		})

		assert.Nil(t, err)
	}

	group.Wait()
	client.Stop()
}
