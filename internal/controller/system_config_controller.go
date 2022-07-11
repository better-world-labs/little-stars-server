package controller

import (
	"aed-api-server/internal/interfaces"
	"github.com/gin-gonic/gin"
	"gitlab.openviewtech.com/openview-pub/gopkg/route"
)

type SystemConfigController struct {
}

//go:inject-component
func SystemNewConfigController() *SystemConfigController {
	return &SystemConfigController{}
}

var SystemSupportedKeys = []string{
	"DEVICE_PICKET",
	"STUDY",
}

func (c SystemConfigController) MountNoAuthRouter(r *route.Router) {
	g := r.Group("/system/configs")
	g.GET("/", c.GetConfig)
	g.GET("/all", c.GetAllConfig)
}

func (c SystemConfigController) MountAuthRouter(r *route.Router) {
	r.POST("/system/configs", c.PutConfig)
}

func (c SystemConfigController) MountAdminRouter(r *route.Router) {
	r.POST("/system/configs", c.PutConfig)
}

func systemKeyValidator(key string) bool {
	for _, it := range SystemSupportedKeys {
		if it == key {
			return true
		}
	}
	return false
}

func (SystemConfigController) PutConfig(c *gin.Context) (interface{}, error) {
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

func (SystemConfigController) GetConfig(c *gin.Context) (interface{}, error) {
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

func (SystemConfigController) GetAllConfig(c *gin.Context) (interface{}, error) {
	config, err := interfaces.S.Config.GetAllConfig()
	if err != nil {
		return nil, err
	}
	c.Header("Content-Type", "application/json; charset=utf-8")
	config = "{\"code\":0,\"data\":" + config + "}"
	c.String(200, config)

	return nil, nil
}
