package oss

import (
	"mime/multipart"
)

// Service oss阿里云存储
type Service interface {
	// oss 文件上传
	// @Param file multipart
	OssUpload(fileheader *multipart.FileHeader) (string, error)
	// oss获取直传token
	GetUploadToken(prefix string, userId int64) (interface{}, error)
}
