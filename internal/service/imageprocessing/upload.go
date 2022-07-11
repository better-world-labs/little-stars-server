package imageprocessing

import (
	"aed-api-server/internal/interfaces"
	"aed-api-server/internal/pkg/db"
	"aed-api-server/internal/pkg/utils"
	"errors"
	"fmt"
	"github.com/aliyun/aliyun-oss-go-sdk/oss"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"io"
	"sync"
	"time"
)

//生成文件上传Oss
//key 在阿里云上存储的key
//genFile 生成文件的函数
//使用方法: 参考 Test_upload 方法
func upload(key string, genFile func(writer *io.PipeWriter)) error {
	r, w := io.Pipe()

	config := interfaces.GetConfig()
	client, err := oss.New(config.AliOss.Endpoint, config.AliOss.AccesskeyId, config.AliOss.AccesskeySecret)
	if err != nil {
		return err
	}

	bucket, err := client.Bucket(config.AliOss.BucketName)
	if err != nil {
		return err
	}

	utils.Go(func() {
		defer func() {
			err := w.Close()
			if err != nil {
				log.Error("upload writer close error:", err)
			}
		}()
		defer func() {
			if err := recover(); err != nil {
				errString := fmt.Sprintf("handle panic: %v, %s", err, utils.PanicTrace(2))
				log.Error(errString)
				err = r.CloseWithError(errors.New(errString))
				if err != nil {
					log.Error("upload reader close error:", err)
				}
			}
		}()

		genFile(w)
	})

	opts := []oss.Option{oss.ContentType("image/jpeg")}
	object, err := bucket.DoPutObject(&oss.PutObjectRequest{
		ObjectKey: key,
		Reader:    r,
	}, opts)

	log.Infof("upload result:%v", object)
	if err != nil {
		return err
	}
	return nil
}

//var group sync.WaitGroup
var chanMap = make(map[string]chan bool)

func getChan(key string) (chan bool, bool) {
	var mutex sync.Mutex
	mutex.Lock()
	defer mutex.Unlock()

	if chanMap[key] == nil {
		chanMap[key] = make(chan bool, 1)
		return chanMap[key], true
	} else {
		return chanMap[key], false
	}
}

func LookUpAndGenPic(
	c *gin.Context,
	makeKey func() (key string, url string, expired *time.Time),
	genFile func(writer *io.PipeWriter),
) error {
	key, saveUrl, expired := makeKey()
	if key == "" {
		panic("key cannot be empty")
	}

	url, err := getImageRedirect(key)
	if err != nil {
		return err
	}

	if url == "" {
		if saveUrl != "" {
			url = saveUrl
		} else {
			url = key
		}
		channel, isNew := getChan(url)
		if isNew {
			utils.Go(func() {
				defer func() {
					close(channel)
					chanMap[key] = nil
				}()
				err = upload(url, genFile)
				if err != nil {
					log.Errorf("%s pic gen error:%v", key, err)
					return
				}
				savePicKeyToDb(key, saveUrl, expired)
			})
		}
		<-channel
	}
	if err != nil {
		return err
	}
	c.Redirect(302, "https://openview-oss.oss-cn-chengdu.aliyuncs.com/"+url)
	return nil
}

type imageRedirect struct {
	Key     string
	Url     string
	Expired *time.Time
}

func savePicKeyToDb(key string, url string, expired *time.Time) {
	if url == "" {
		url = key
	}
	redirect := imageRedirect{Key: key, Url: url, Expired: expired}
	_, err := db.Insert("image_redirect", redirect)
	if err != nil {
		log.Errorf("savePicKeyToDb error:%v", err)
	}
}

func getImageRedirect(key string) (string, error) {
	var bin imageRedirect
	existed, err := db.Table("image_redirect").
		Select("url").Where("`key`=? and (expired > now() or expired is null)", key).Get(&bin)
	if err != nil {
		return "", err
	}
	if existed {
		return bin.Url, nil
	}
	return "", nil
}
