package project

import (
	"aed-api-server/internal/interfaces"
	"aed-api-server/internal/pkg/db"
	"aed-api-server/internal/server/config"
	"github.com/stretchr/testify/assert"
	"testing"
)

func initDbAndConfig() func() {
	c, err := config.LoadConfig("../../../")
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

func Test_UpdateUserProjectLevel(t *testing.T) {
	t.Cleanup(initDbAndConfig())
	service := Service{}
	err := service.UpdateUserProjectLevel(10, 20, 3)
	assert.Nil(t, err)
}
