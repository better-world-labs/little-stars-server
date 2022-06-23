package db

import "time"

type MysqlConfig struct {
	DriverName   string        `properties:"driver-name"`
	Dsn          string        `properties:"dsn"`
	MaxIdleCount int           `properties:"max-idle-count"`
	MaxOpen      int           `properties:"max-open"`
	MaxLifetime  time.Duration `properties:"max-lifetime"`
}
