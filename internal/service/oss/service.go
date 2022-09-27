package oss

import (
	"aed-api-server/internal/interfaces/service"
	"crypto/hmac"
	"crypto/sha1"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/aliyun/aliyun-oss-go-sdk/oss"
	"hash"
	"io"
	"log"
	"mime/multipart"
	"time"
)

type ossService struct {
	client *oss.Client

	Endpoint        string `conf:"alioss.endpoint"`
	AccesskeyId     string `conf:"alioss.accesskey-id"`
	AccesskeySecret string `conf:"alioss.accesskey-secret"`
	BucketName      string `conf:"alioss.bucket-name"`
	Domain          string `conf:"alioss.domain"`
	UploadDir       string `conf:"alioss.upload-dir"`
}

const (
	ExpireTime int64 = 30
)

//go:inject-component
func NewOssService() service.OssService {
	return &ossService{}
}

func (s *ossService) getClient() *oss.Client {
	if s.client == nil {
		client, _ := oss.New(s.Endpoint, s.AccesskeyId, s.AccesskeySecret)
		s.client = client
	}
	return s.client
}

func (s *ossService) Upload(path string, reader io.Reader) (string, error) {
	client := s.getClient()

	bucket, err := client.Bucket(s.BucketName)
	if err != nil {
		log.Printf("oss bucket create failed, msg error: %s", err)
		return "", err
	}

	path = fmt.Sprintf("%s/%s", s.UploadDir, path)
	err = bucket.PutObject(path, reader)
	if err != nil {
		log.Printf("oss upload file failed, msg error: %s", err)
		return "", err
	}

	return fmt.Sprintf("https://%s/%s", s.Domain, path), nil
}

func (s *ossService) OssUpload(fileheader *multipart.FileHeader) (string, error) {
	client := s.getClient()

	file, err := fileheader.Open()
	if err != nil {
		log.Printf("open file error: %s", err)
		return "", err
	}

	bucket, err := client.Bucket(s.BucketName)
	if err != nil {
		log.Printf("oss bucket create failed, msg error: %s", err)
		return "", err
	}

	objectName := fileheader.Filename
	err = bucket.PutObject(objectName, file)
	if err != nil {
		log.Printf("oss upload file failed, msg error: %s", err)
		return "", err
	}

	return fmt.Sprintf("http://%s/%s", s.Domain, objectName), nil
}

func (s *ossService) GetUploadToken(prefix string, accountID int64) (interface{}, error) {
	now := time.Now().Unix()
	expire_end := now + ExpireTime
	var tokenExpire = get_gmt_iso8601(expire_end)

	// create post policy json
	var config ConfigStruct
	config.Expiration = tokenExpire
	var condition []string
	uploadDir := fmt.Sprintf("%s%s/%v/", prefix, s.UploadDir, accountID) // 上传目录 /aed/用户id/
	condition = append(condition, "starts-with")
	condition = append(condition, "$key")
	condition = append(condition, uploadDir)
	config.Conditions = append(config.Conditions, condition)

	// calculate signature
	result, err := json.Marshal(config)
	if err != nil {
		return "", nil
	}
	debyte := base64.StdEncoding.EncodeToString(result)
	h := hmac.New(func() hash.Hash { return sha1.New() }, []byte(s.AccesskeySecret))
	io.WriteString(h, debyte)
	signedStr := base64.StdEncoding.EncodeToString(h.Sum(nil))

	var policyToken PolicyToken
	policyToken.AccessKeyId = s.AccesskeyId
	policyToken.Host = fmt.Sprintf("https://%s", s.Domain)
	policyToken.Expire = expire_end
	policyToken.Signature = string(signedStr)
	policyToken.Directory = uploadDir
	policyToken.Policy = string(debyte)
	policyToken.FileName = fmt.Sprintf("%v_%v", accountID, time.Now().UnixNano())

	return policyToken, nil
}

func get_gmt_iso8601(expire_end int64) string {
	var tokenExpire = time.Unix(expire_end, 0).UTC().Format("2006-01-02T15:04:05Z")
	return tokenExpire
}
