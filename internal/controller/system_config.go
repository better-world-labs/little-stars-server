package controller

import (
	"aed-api-server/internal/interfaces"
	"aed-api-server/internal/pkg/response"
	"github.com/gin-gonic/gin"
)

type ConfigController struct {
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

func (ConfigController) PutConfig(c *gin.Context) {
	key := c.Query("key")

	data, err := c.GetRawData()
	if err != nil {
		response.ReplyError(c, err)
		return
	}

	err = interfaces.S.Config.PutConfig(key, string(data))
	if err != nil {
		response.ReplyError(c, err)
		return
	}
	response.ReplyOK(c, nil)

}

func (ConfigController) GetConfig(c *gin.Context) {
	key := c.Query("key")

	config, err := interfaces.S.Config.GetConfig(key)
	if err != nil {
		response.ReplyError(c, err)
		return
	}
	c.Header("Content-Type", "application/json; charset=utf-8")
	config = "{\"code\":0,\"data\":" + config + "}"
	c.String(200, config)
}

func (ConfigController) GetAllConfig(c *gin.Context) {
	config, err := interfaces.S.Config.GetAllConfig()
	if err != nil {
		response.ReplyError(c, err)
		return
	}
	c.Header("Content-Type", "application/json; charset=utf-8")
	config = "{\"code\":0,\"data\":" + config + "}"
	c.String(200, config)
}
