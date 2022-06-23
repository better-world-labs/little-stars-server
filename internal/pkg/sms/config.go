package sms

type Config struct {
	AccessKeyID     string `properties:"access-key-id"`
	AccessKeySecret string `properties:"access-key-secret"`
	Enabled         bool   `properties:"enabled"`
}
