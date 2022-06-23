package controller

import (
	"aed-api-server/internal/interfaces"
	"github.com/gin-gonic/gin"
	"gitlab.openviewtech.com/openview-pub/gopkg/route"
)

type ConfigController struct {
}

func NewConfigController() *ConfigController {
	return &ConfigController{}
}

var SupportedKeys = []string{
	"DEVICE_PICKET",
	"STUDY",
}

func (c ConfigController) MountNoAuthRouter(r *route.Router) {
	g := r.Group("/system/configs")
	g.GET("/", c.GetConfig)
	g.GET("/all", c.GetAllConfig)
}

func (c ConfigController) MountAuthRouter(r *route.Router) {
	r.POST("/system/configs", c.PutConfig)
}

func (c ConfigController) MountAdminRouter(r *route.Router) {
	r.POST("/system/configs", c.PutConfig)
}

func keyValidator(key string) bool {
	for _, it := range SupportedKeys {
		if it == key {
			return true
		}
	}
	return false
}

func (ConfigController) PutConfig(c *gin.Context) (interface{}, error) {
	key := c.Query("key")

	data, err := c.GetRawData()
	if err != nil {
		return nil, err
	}

	err = interfaces.S.Config.PutConfig(key, string(data))
	if err != nil {
		return nil, err
	}

	return nil, nil
}

func (ConfigController) GetConfig(c *gin.Context) (interface{}, error) {
	key := c.Query("key")

	config, err := interfaces.S.Config.GetConfig(key)
	if err != nil {
		return nil, err
	}
	c.Header("Content-Type", "application/json; charset=utf-8")
	config = "{\"code\":0,\"data\":" + config + "}"
	c.String(200, config)

	return nil, nil
}

func (ConfigController) GetAllConfig(c *gin.Context) (interface{}, error) {
	config, err := interfaces.S.Config.GetAllConfig()
	if err != nil {
		return nil, err
	}
	c.Header("Content-Type", "application/json; charset=utf-8")
	config = "{\"code\":0,\"data\":" + config + "}"
	c.String(200, config)

	return nil, nil
}
