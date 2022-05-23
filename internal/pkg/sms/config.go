package sms

type Config struct {
	AccessKeyID     string `yaml:"access-key-id"`
	AccessKeySecret string `yaml:"access-key-secret"`
	Enabled         bool   `yaml:"enabled"`
}
