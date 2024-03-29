package controller

import (
	"aed-api-server/internal/interfaces"
	"aed-api-server/internal/pkg"
	"aed-api-server/internal/pkg/global"
	"github.com/gin-gonic/gin"
	"gitlab.openviewtech.com/openview-pub/gopkg/route"
)

type UserConfigController struct {
}

//go:inject-component
func UserNewConfigController() *UserConfigController {
	return &UserConfigController{}
}

func (c UserConfigController) MountAuthRouter(r *route.Router) {
	configR := r.Group("/configs")
	configR.GET("", c.GetConfig)
	configR.PUT("", c.PutConfig)

	v2 := r.Group("/v2/configs")
	v2.GET("", GetConfigV2)
}

func GetConfigV2(c *gin.Context) (interface{}, error) {
	userId := c.MustGet(pkg.AccountIDKey).(int64)
	key := c.Query("key")

	config, err := interfaces.S.UserConfig.GetConfig(userId, key)
	if err != nil {
		return nil, err
	}

	if config == nil {
		return nil, nil
	}

	return map[string]interface{}{
		"value":     config.Value,
		"createdAt": global.FormattedTime(config.CreatedAt),
	}, nil
}

func (UserConfigController) PutConfig(c *gin.Context) (interface{}, error) {
	userId := c.MustGet(pkg.AccountIDKey).(int64)
	key := c.Query("key")

	data, err := c.GetRawData()
	if err != nil {
		return nil, err
	}
	_, err = interfaces.S.UserConfig.PutConfig(userId, key, string(data))
	if err != nil {
		return nil, err
	}
	return nil, nil
}

func (UserConfigController) GetConfig(c *gin.Context) (interface{}, error) {
	userId := c.MustGet(pkg.AccountIDKey).(int64)
	key := c.Query("key")

	config, err := interfaces.S.UserConfig.GetConfig(userId, key)
	if err != nil {
		return nil, err
	}
	res := "{\"code\":0,\"data\":" + config.Value + "}"
	c.Header("Content-Type", "application/json; charset=utf-8")
	c.String(200, res)

	return nil, nil
}
