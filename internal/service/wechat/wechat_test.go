package wechat

import (
	"aed-api-server/internal/interfaces"
	"aed-api-server/internal/pkg/db"
	"aed-api-server/internal/server/config"
	"github.com/stretchr/testify/assert"
	"gitlab.openviewtech.com/openview-pub/gopkg/log"
	"testing"
)

func Test_getAccessTokenFromWx(t *testing.T) {
	c, err := config.LoadConfig("../../../")
	if err != nil {
		panic("get config error")
	}
	interfaces.InitConfig(c)
	db.InitEngine(c.Database)
	log.Init(c.Log)

	//wx, i, err := getAccessTokenFromWx()
	//assert.Nil(t, err)
	//
	//log.Info("1. wx=", wx, ",i=", i)
	//
	//wx, i, err = getAccessTokenFromWx()
	assert.Nil(t, err)

	//log.Info("2. wx=", wx, ",i=", i)
}
