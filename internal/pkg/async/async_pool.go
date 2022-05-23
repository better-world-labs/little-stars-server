package async

import (
	"log"
	"sync"
)

var TaskPool *Pool

// Pool Go routine Pool
type Pool struct {
	goRoutinePoolSize int
	maxTaskQueueSize  int
	asyncTaskQueue    chan func()
	started           bool
	group             sync.WaitGroup
}

func NewPool(c Config) *Pool {
	return &Pool{
		goRoutinePoolSize: c.GoRoutinePoolSize,
		maxTaskQueueSize:  c.MaxTaskQueueSize,
		asyncTaskQueue:    make(chan func(), c.MaxTaskQueueSize),
		group:             sync.WaitGroup{},
	}
}

func (c *Pool) Start() {
	c.group.Add(c.goRoutinePoolSize)

	for i := 0; i < c.goRoutinePoolSize; i++ {
		go func() {
			defer func() {
				err := recover()
				if err != nil {
					log.Printf("task run error: %v", recover())
				}
			}()

			for task := range c.asyncTaskQueue {
				task()
			}

			c.group.Done()
		}()
	}

	log.Printf("async task pool started")
	c.started = true
}

func (c *Pool) AddSyncTask(task func()) (err error) {
	if !c.started {
		err = ErrorCientStoped
	} else {
		c.asyncTaskQueue <- task
	}

	return
}
func (c *Pool) Stop() {
	close(c.asyncTaskQueue)
	c.started = false
	c.group.Wait()
}
