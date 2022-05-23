package oss

// 阿里云oss配置
type Config struct {
	Endpoint        string `yaml:"endpoint"`
	AccesskeyId     string `yaml:"accesskey-id"`
	AccesskeySecret string `yaml:"accesskey-secret"`
	BucketName      string `yaml:"bucket-name"`
	Domain          string `yaml:"domain"`
	UploadDir       string `yaml:"upload-dir"`
}

// oss直传配置
type ConfigStruct struct {
	Expiration string     `json:"expiration"`
	Conditions [][]string `json:"conditions"`
}

// oss直传token结构体
type PolicyToken struct {
	AccessKeyId string `json:"accessid"`
	Host        string `json:"host"`
	Expire      int64  `json:"expire"`
	Signature   string `json:"signature"`
	Policy      string `json:"policy"`
	Directory   string `json:"dir"`
	Callback    string `json:"callback"`
	FileName    string `json:"filename"`
}
