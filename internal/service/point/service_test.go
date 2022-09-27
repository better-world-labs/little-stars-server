package point

import (
	"aed-api-server/internal/interfaces"
	"aed-api-server/internal/interfaces/entities"
	"aed-api-server/internal/pkg/db"
	"aed-api-server/internal/pkg/global"
	"aed-api-server/internal/server/config"
	"github.com/stretchr/testify/assert"
	"gitlab.openviewtech.com/openview-pub/gopkg/log"
	"testing"
	"time"
)

var userId int64 = 49
var addPoints = 100

func InitDbAndConfig() func() {
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

func Test_GetUserTotalPoints(t *testing.T) {
	t.Cleanup(InitDbAndConfig())
	s := NewService()
	points, err := s.GetUserTotalPoints(userId)
	assert.Nil(t, err)
	assert.True(t, points > 0)
	log.Info(points)
}

func Test_GetUnReceivePoints(t *testing.T) {
	t.Cleanup(InitDbAndConfig())
	s := NewService()
	_, err := s.GetUnReceivePoints(userId)
	assert.Nil(t, err)
}

func Test_insertPoints_ReceivePoints(t *testing.T) {
	s := NewService()
	t.Cleanup(InitDbAndConfig())
	now := time.Now()
	expired := now.Add(10 * time.Second)
	err := insertPoints(userId, addPoints, entities.PointsEventTypeActivityGive, "", entities.PointsEventParams{
		RefTable:   "test",
		RefTableId: 100,
	}, expired)
	assert.Nil(t, err)

	points, err := s.GetUnReceivePoints(userId)
	assert.Nil(t, err)
	n := len(points)
	assert.True(t, n > 0)
	flow := points[0]
	assert.Equal(t, flow.Points, addPoints)

	json, _ := global.FormattedTime(expired).MarshalJSON()
	marshalJSON, _ := flow.ExpiredAt.MarshalJSON()
	assert.Equal(t, json, marshalJSON)

	err = s.ReceivePoints(userId, flow.Id)
	assert.Nil(t, err)

	points, err = s.GetUnReceivePoints(userId)
	assert.Nil(t, err)
	assert.Equal(t, len(points), n-1)
}
