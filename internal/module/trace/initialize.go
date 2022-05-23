package trace

import (
	"aed-api-server/internal/interfaces"
	"aed-api-server/internal/module/user"
	"aed-api-server/internal/pkg/config"
)

func Init(c *config.AppConfig) {
	initService(c)
}

func initService(c *config.AppConfig) {
	interfaces.S.Trace = NewTraceService(user.NewWechatClient(&c.Wechat))
}
