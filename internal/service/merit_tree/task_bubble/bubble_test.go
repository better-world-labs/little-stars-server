package task_bubble

import (
	"aed-api-server/internal/interfaces/events"
	"aed-api-server/internal/pkg/domain/emitter"
	"aed-api-server/internal/server/config"
	"context"
	"gitlab.openviewtech.com/openview-pub/gopkg/log"
	"testing"
	"time"
)

func TestEventEmitted(t *testing.T) {
	c, err := config.LoadConfig("../../../../")
	if err != nil {
		panic("get config error")
	}

	emitter.SetContext(context.Background())
	emitter.SetConfig(&c.Domain)
	emitter.Start()

	err = emitter.Emit(&events.FirstLoginEvent{
		UserId:  115,
		Openid:  "oyL7e5fPMvMauIzcC3t9rX7RV1dM",
		LoginAt: time.Now(),
	})
	if err != nil {
		log.Info("err", err)
	}

	time.Sleep(5 * time.Second)
}
