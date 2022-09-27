package config

type RocketConf struct {
	EndPoint   string `properties:"endpoint"`
	InstanceId string `properties:"instanceId"`
	AccessKey  string `properties:"accessKey"`
	SecretKey  string `properties:"secretKey"`
	Topic      string `properties:"topic"`
	GroupId    string `properties:"groupId"`
	DebugTag   string `properties:"debug"`
}
