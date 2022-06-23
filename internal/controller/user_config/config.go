package user_config

import (
	"aed-api-server/internal/interfaces"
	"aed-api-server/internal/pkg"
	"errors"
	"github.com/gin-gonic/gin"
	"gitlab.openviewtech.com/openview-pub/gopkg/route"
)

type Controller struct {
}

func NewController() *Controller {
	return &Controller{}
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

func (c Controller) MountAuthRouter(r *route.Router) {
	configR := r.Group("/configs")
	configR.GET("", c.GetConfig)
	configR.PUT("", c.PutConfig)

}

func (Controller) PutConfig(c *gin.Context) (interface{}, error) {
	userId := c.MustGet(pkg.AccountIDKey).(int64)
	key := c.Query("key")

	if !keyValidator(key) {
		return nil, errors.New("key only support DEVICE_PICKET、STUDY")
	}
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

func (Controller) GetConfig(c *gin.Context) (interface{}, error) {
	userId := c.MustGet(pkg.AccountIDKey).(int64)
	key := c.Query("key")

	if !keyValidator(key) {
		return nil, errors.New("key only support DEVICE_PICKET、STUDY")
	}

	config, err := interfaces.S.UserConfig.GetConfig(userId, key)
	if err != nil {
		return nil, err
	}
	config = "{\"code\":0,\"data\":" + config + "}"
	c.Header("Content-Type", "application/json; charset=utf-8")
	c.String(200, config)

	return nil, nil
}
