package friends

import (
	"aed-api-server/internal/interfaces"
	"aed-api-server/internal/interfaces/entities"
	"aed-api-server/internal/interfaces/events"
	"aed-api-server/internal/pkg/domain/emitter"
	"fmt"
	"time"
)

func Init() {
	interfaces.S.Friends = NewService()
	emitter.On(&events.FirstLoginEvent{}, handleNewUserLogin)
	emitter.On(&entities.Trace{}, handleNewTraceCreated)
	initCron()
}

const (
	H = 0
	M = 0
	S = 0
)

var t *time.Timer

// TODO 实现全局定时任务调度器
func initCron() {
	duration := getDuration(H, M, S)
	t = time.NewTimer(duration)
	go run()
}

func getDuration(h int, m int, s int) time.Duration {
	now := time.Now()
	next := time.Date(now.Year(), now.Month(), now.Day(), h, m, s, 0, now.Location())
	if next.Sub(now) < 0 {
		next = next.Add(24 * time.Hour)
	}

	fmt.Printf("[friends] task run on %s\n", next)
	return next.Sub(time.Now())
}

func run() {
	for {
		<-t.C
		doCron()
		t.Reset(getDuration(H, M, S))
	}
}
