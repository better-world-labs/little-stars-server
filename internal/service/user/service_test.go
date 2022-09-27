package user

import (
	"aed-api-server/internal/interfaces"
	"aed-api-server/internal/interfaces/entities"
	"aed-api-server/internal/pkg/cache"
	"aed-api-server/internal/pkg/db"
	conf "aed-api-server/internal/server/config"
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestService_BatchGetLastUserEventByType(t *testing.T) {
	c, err := conf.LoadConfig("../../../")
	if err != nil {
		panic("get config error")
	}
	interfaces.InitConfig(c)
	db.InitEngine(c.Database)

	service := Service{}
	byType, err := service.BatchGetLastUserEventByType([]int64{50}, entities.UserEventTypeGetWalkStep)
	assert.Nil(t, err)

	fmt.Printf("byType:%v", byType)
}

func TestService_GetUserEncryptKey(t *testing.T) {
	c, err := conf.LoadConfig("../../../")
	if err != nil {
		panic("get config error")
	}
	interfaces.InitConfig(c)
	db.InitEngine(c.Database)
	cache.InitPool(c.Redis)

	service := Service{
		//Wechat:  mock.NewWechatMock(),
		Encrypt: NewCryptKeyCache(),
	}
	service.Encrypt.Env = "local"

	key, err := service.GetUserEncryptKey(50, 3)
	assert.Nil(t, err)
	fmt.Printf("%s\n", key.EncryptKey)
}
