package oss

import (
	"fmt"
	"log"
	"mime/multipart"

	"github.com/aliyun/aliyun-oss-go-sdk/oss"

	"crypto/hmac"
	"crypto/sha1"
	"encoding/base64"
	"encoding/json"

	"hash"
	"io"
	"time"
)

type service struct {
	client *oss.Client
	c      *Config
}

const (
	ExpireTime int64 = 30
)

func NewService(c *Config) Service {
	client, err := oss.New(c.Endpoint, c.AccesskeyId, c.AccesskeySecret)
	if err != nil {
		panic(fmt.Sprintf("oss client create failed, msg error: %s", err))
	}
	return &service{client: client, c: c}
}

func (s service) OssUpload(fileheader *multipart.FileHeader) (string, error) {
	client := s.client

	file, err := fileheader.Open()
	if err != nil {
		log.Printf("open file error: %s", err)
		return "", err
	}

	bucket, err := client.Bucket(s.c.BucketName)
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

	return fmt.Sprintf("http://%s/%s", s.c.Domain, objectName), nil
}

func (s service) GetUploadToken(accountID int64) (interface{}, error) {
	now := time.Now().Unix()
	expire_end := now + ExpireTime
	var tokenExpire = get_gmt_iso8601(expire_end)

	// create post policy json
	var config ConfigStruct
	config.Expiration = tokenExpire
	var condition []string
	uploadDir := fmt.Sprintf("%s/%v/", s.c.UploadDir, accountID) // 上传目录 /aed/用户id/
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
	h := hmac.New(func() hash.Hash { return sha1.New() }, []byte(s.c.AccesskeySecret))
	io.WriteString(h, debyte)
	signedStr := base64.StdEncoding.EncodeToString(h.Sum(nil))

	var policyToken PolicyToken
	policyToken.AccessKeyId = s.c.AccesskeyId
	policyToken.Host = fmt.Sprintf("https://%s", s.c.Domain)
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
