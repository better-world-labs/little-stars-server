package user_config

import (
	"aed-api-server/internal/interfaces"
	"aed-api-server/internal/pkg"
	"aed-api-server/internal/pkg/response"
	"errors"
	"github.com/gin-gonic/gin"
)

type Controller struct {
}

var SupportedKeys = []string{
	"DEVICE_PICKET",
	"STUDY",
}

func keyValidator(key string) bool {
	for _, it := range SupportedKeys {
		if it == key {
			return true
		}
	}
	return false
}

func (Controller) PutConfig(c *gin.Context) {
	userId := c.MustGet(pkg.AccountIDKey).(int64)
	key := c.Query("key")

	if !keyValidator(key) {
		response.ReplyError(c, errors.New("key only support DEVICE_PICKET、STUDY"))
		return
	}
	data, err := c.GetRawData()
	if err != nil {
		response.ReplyError(c, err)
		return
	}
	_, err = interfaces.S.UserConfig.PutConfig(userId, key, string(data))
	if err != nil {
		response.ReplyError(c, err)
		return
	}
	response.ReplyOK(c, nil)

}

func (Controller) GetConfig(c *gin.Context) {
	userId := c.MustGet(pkg.AccountIDKey).(int64)
	key := c.Query("key")

	if !keyValidator(key) {
		response.ReplyError(c, errors.New("key only support DEVICE_PICKET、STUDY"))
		return
	}

	config, err := interfaces.S.UserConfig.GetConfig(userId, key)
	if err != nil {
		response.ReplyError(c, err)
		return
	}
	config = "{\"code\":0,\"data\":" + config + "}"
	c.Header("Content-Type", "application/json; charset=utf-8")
	c.String(200, config)
}
