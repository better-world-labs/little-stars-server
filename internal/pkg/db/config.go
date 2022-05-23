package db

import "time"

type MysqlConfig struct {
	DriverName   string        `yaml:"driver-name"`
	Dsn          string        `yaml:"dsn"`
	MaxIdleCount int           `yaml:"max-idle-count"`
	MaxOpen      int           `yaml:"max-open"`
	MaxLifetime  time.Duration `yaml:"max-lifetime"`
	MaxIdleTime  time.Duration `yaml:"maxI-idle-time"`
}
