package interfaces

import (
	"aed-api-server/internal/server/config"
)

// 存放 配置相关的结构

var _c *config.AppConfig

func InitConfig(c *config.AppConfig) {
	_c = c
}

// GetConfig 获取配置信息
func GetConfig() *config.AppConfig {
	return _c
}
