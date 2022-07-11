package service

import (
	"mime/multipart"
)

// OssService oss阿里云存储
type OssService interface {
	// OssUpload 文件上传
	// @Param file multipart
	OssUpload(fileHeader *multipart.FileHeader) (string, error)

	// GetUploadToken oss获取直传token
	GetUploadToken(prefix string, userId int64) (interface{}, error)
}
