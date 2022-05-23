package oss

import (
	"aed-api-server/internal/pkg"
	"aed-api-server/internal/pkg/response"

	"github.com/gin-gonic/gin"
	"gitlab.openviewtech.com/openview-pub/gopkg/log"
)

type Controller struct {
	service Service
}

func NewController(s Service) *Controller {
	return &Controller{service: s}
}

// 应用服务器上传
func (con *Controller) Upload(c *gin.Context) {
	fileHeader, err := c.FormFile("files")
	if err != nil {
		log.DefaultLogger().Errorf("get file error: %v", err)
		response.ReplyError(c, err)
		return
	}
	url, err := con.service.OssUpload(fileHeader)
	if err != nil {
		log.DefaultLogger().Errorf("OssUpload error: %v", err)
		response.ReplyError(c, err)
		return
	}

	response.ReplyOK(c, &UploadDto{
		Origin: url,
	})
}

// 获取阿里云直传Token
func (con *Controller) GetUploadToken(c *gin.Context) {
	accountID := c.MustGet(pkg.AccountIDKey).(int64)
	policy, err := con.service.GetUploadToken(accountID)
	if err != nil {
		log.DefaultLogger().Errorf("GetUploadToken error: %v", err)
		response.ReplyError(c, err)
		return
	}
	response.ReplyOK(c, policy)
}
