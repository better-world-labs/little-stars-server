package task

import (
	"aed-api-server/internal/interfaces"
	"aed-api-server/internal/pkg/config"
	"aed-api-server/internal/pkg/db"
	"time"
)

const testUserId = 1
const testTaskId = 100
const testDeviceId = "a472dc08-c618-49c8-9ca6-c252365e1afa"
const _1d = 24 * time.Hour
const pageSize = 10

func InitDbAndConfig() func() {
	c, err := config.LoadConfig("../../../config-local.yaml")
	if err != nil {
		panic("get config error")
	}
	interfaces.InitConfig(c)
	db.InitEngine(c.Database)

	return func() {
		println("close db")
		engine := db.GetEngine()
		if engine != nil {
			engine.Close()
		}
	}
}

type MPicketCondition struct{}

func (p MPicketCondition) IsPicketNone(deviceId string) bool {
	return true
}
func (p MPicketCondition) IsLastTwiceConflict(deviceId string) bool {
	return true
}
func (p MPicketCondition) IsLastTwiceFalse(deviceId string) bool {
	return true
}

func mockPicketCondition() {
	interfaces.S.PicketCondition = MPicketCondition{}
}
