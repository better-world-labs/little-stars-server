package config

import (
	"gopkg.in/yaml.v2"
	"os"
)

func LoadConfig(filename string) (*AppConfig, error) {
	var config AppConfig
	f, err := os.Open(filename)
	if err != nil {
		return nil, err
	}

	err = yaml.NewDecoder(f).Decode(&config)
	if err != nil {
		return nil, err
	}

	return &config, nil
}
