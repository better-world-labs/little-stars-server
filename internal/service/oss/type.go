package oss

// 阿里云oss配置
type Config struct {
	Endpoint        string `properties:"endpoint"`
	AccesskeyId     string `properties:"accesskey-id"`
	AccesskeySecret string `properties:"accesskey-secret"`
	BucketName      string `properties:"bucket-name"`
	Domain          string `properties:"domain"`
	UploadDir       string `properties:"upload-dir"`
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
