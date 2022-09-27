package config

import (
	"github.com/magiconair/properties"
	"gitlab.openviewtech.com/openview-pub/gopkg/conf"
)

func LoadConfig(baseDir string, envArr ...string) (*AppConfig, error) {
	x, _, err := LoadConfigX(baseDir, envArr...)
	return x, err
}

func LoadConfigX(baseDir string, envArr ...string) (*AppConfig, *properties.Properties, error) {
	var env = "local"
	if len(envArr) > 0 {
		env = envArr[0]
	}
	pros := conf.Init(baseDir, env)

	var config AppConfig
	err := pros.Decode(&config)
	return &config, pros, err
}
