package point_test

import (
	"aed-api-server/internal/interfaces"
	"aed-api-server/internal/interfaces/entities"
	"aed-api-server/internal/interfaces/events"
	"aed-api-server/internal/pkg/config"
	"aed-api-server/internal/pkg/db"
	conf "aed-api-server/internal/pkg/domain/config"
	"aed-api-server/internal/pkg/domain/emitter"
	"aed-api-server/internal/service/point"
	"context"
	"github.com/stretchr/testify/assert"
	"gitlab.openviewtech.com/openview-pub/gopkg/log"
	"testing"
	"time"
)

func initMiddlewares() func() {
	c, err := config.LoadConfig("../../../")
	if err != nil {
		panic("get config error")
	}
	interfaces.InitConfig(c)
	db.InitEngine(c.Database)
	log.Init(c.Log)

	return func() {
		println("close db")
		engine := db.GetEngine()
		if engine != nil {
			engine.Close()
		}
	}
}

func TestDealPointsEvent(t *testing.T) {
	t.Cleanup(initMiddlewares())
	rst, err := interfaces.S.PointsScheduler.DealPointsEvent(&events.PointsEvent{
		PointsEventType: entities.PointsEventTypeLearntCourse,
		UserId:          49,
		Params: entities.PointsEventParams{
			RefTable:   "test",
			RefTableId: 1,
		},
	})
	assert.Nil(t, err)
	assert.Equal(t, int64(49), rst.UserId)
	assert.Equal(t, 50, rst.PeckingPointsChange)
}

func TestEmitter(t *testing.T) {
	t.Cleanup(initMiddlewares())
	emitter.SetContext(context.Background())
	emitter.SetConfig(&conf.DomainEventConfig{
		Server:          "kafka-star-dev.openviewtech.com:9092",
		Topic:           "test-domain-event",
		GroupId:         "star-service",
		DeadLetterTopic: "test-dead",
	})

	point.InitEventHandler()
	emitter.Start()

	err := emitter.Emit(&events.PointsEvent{
		PointsEventType: entities.PointsEventTypeLearntCourse,
		UserId:          49,
		Params: entities.PointsEventParams{
			RefTable:   "test",
			RefTableId: 1,
		},
	})

	assert.Nil(t, err)
	time.Sleep(10 * time.Second)
}
