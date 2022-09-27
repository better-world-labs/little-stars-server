package oss

import "testing"

func TestUpload(t *testing.T) {
	service := NewOssService().(*ossService)
	service.C = &Config{
		Endpoint:        "oss-cn-chengdu.aliyuncs.com",
		AccesskeyId:     "LTAI5tN8fcJYPbrQGetxFL7c",
		AccesskeySecret: "4o40jlAWKPboMFbk7MajzMOBX7fNbj",
		BucketName:      "openview-oss",
		Domain:          "openview-oss.oss-cn-chengdu.aliyuncs.com",
		UploadDir:       "aed-dev",
	}

}
