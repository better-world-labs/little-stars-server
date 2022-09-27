package controller

import (
	"aed-api-server/internal/interfaces"
	"aed-api-server/internal/interfaces/entities"
	"aed-api-server/internal/pkg"
	"aed-api-server/internal/pkg/global"
	"aed-api-server/internal/pkg/utils"
	"encoding/json"
	"errors"
	"github.com/gin-gonic/gin"
	"gitlab.openviewtech.com/openview-pub/gopkg/route"
)

type UserConfigV2Controller struct {
}

//go:inject-component
func UserNewConfigV2Controller() *UserConfigV2Controller {
	return &UserConfigV2Controller{}
}

func (c UserConfigV2Controller) MountAuthRouter(r *route.Router) {
	configR := r.Group("/v3")
	configR.GET("/configs", c.GetConfig)
	configR.POST("/get-config", c.GetConfigWithDefault)
	configR.PUT("/configs", c.PutConfig)
}

func (UserConfigV2Controller) PutConfig(c *gin.Context) (interface{}, error) {
	userId := c.MustGet(pkg.AccountIDKey).(int64)
	key := c.Query("key")

	var value struct {
		Value interface{} `json:"value" binding:"required"`
	}

	err := c.ShouldBindJSON(&value)
	if err != nil {
		return nil, err
	}

	_, err = interfaces.S.UserConfigV2.Put(userId, key, value)

	if err != nil {
		return nil, err
	}
	return nil, nil
}

func (UserConfigV2Controller) GetConfig(ctx *gin.Context) (interface{}, error) {
	userId := utils.GetContextUserId(ctx)
	key, exists := ctx.GetQuery("key")
	if !exists {
		return nil, errors.New("invalid param")
	}

	config, err := interfaces.S.UserConfigV2.Get(userId, key)
	if err != nil {
		return nil, err
	}

	if config == nil {
		return nil, nil
	}

	return map[string]interface{}{
		"value":     config.Value["value"],
		"createdAt": global.FormattedTime(config.CreatedAt),
	}, nil
}

func (c UserConfigV2Controller) GetConfigWithDefault(ctx *gin.Context) (interface{}, error) {
	userId := utils.GetContextUserId(ctx)
	var param struct {
		Key             string      `json:"key" binding:"required"`
		DefaultValue    interface{} `json:"defaultValue"`
		SystemConfigKey string      `json:"systemConfigKey"`
	}

	err := ctx.ShouldBindJSON(&param)
	var value *entities.UserConfig
	if err != nil {
		return nil, err
	}

	if param.SystemConfigKey != "" {
		config, err := interfaces.S.Config.GetConfig(param.SystemConfigKey)
		if err != nil {
			return nil, err
		}

		var v interface{}
		err = json.Unmarshal([]byte(config), &v)
		if err != nil {
			v = config
		}

		value, err = interfaces.S.UserConfigV2.GetOrDefault(userId, param.Key, v)
		if err != nil {
			return nil, err
		}
	} else if param.DefaultValue != nil {
		value, err = interfaces.S.UserConfigV2.GetOrDefault(userId, param.Key, param.DefaultValue)
		if err != nil {
			return nil, err
		}
	} else {
		value, err = interfaces.S.UserConfigV2.Get(userId, param.Key)
		if err != nil {
			return nil, err
		}
	}

	return map[string]interface{}{
		"value":     value.Value["value"],
		"createdAt": global.FormattedTime(value.CreatedAt),
	}, nil
}
