package oss

import (
	"aed-api-server/internal/pkg"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"gitlab.openviewtech.com/openview-pub/gopkg/route"
)

type Controller struct {
	service Service
}

func NewController(s Service) *Controller {
	return &Controller{service: s}
}

func (con *Controller) MountAdminRouter(r *route.Router) {
	r.GET("/upload-token", con.GetUploadTokenAdmin)
}

func (con *Controller) MountAuthRouter(r *route.Router) {
	r.POST("/common/photo", con.Upload)
	r.GET("/common/upload_token", con.GetUploadToken)
}

// 应用服务器上传
func (con *Controller) Upload(c *gin.Context) (interface{}, error) {
	fileHeader, err := c.FormFile("files")
	if err != nil {
		log.Errorf("get file error: %v", err)
		return nil, err
	}
	url, err := con.service.OssUpload(fileHeader)
	if err != nil {
		log.Errorf("OssUpload error: %v", err)
		return nil, err
	}

	return &UploadDto{
		Origin: url,
	}, nil
}

// 获取阿里云直传Token
func (con *Controller) GetUploadToken(c *gin.Context) (interface{}, error) {
	accountID := c.MustGet(pkg.AccountIDKey).(int64)
	policy, err := con.service.GetUploadToken("", accountID)
	if err != nil {
		log.Errorf("GetUploadToken error: %v", err)
		return nil, err
	}
	return policy, nil
}

func (con *Controller) GetUploadTokenAdmin(context *gin.Context) (interface{}, error) {
	accountID := context.MustGet(pkg.AccountIDKey).(int64)
	if accountID < 0 {
		accountID = -accountID
	}
	policy, err := con.service.GetUploadToken("stars-admin/", accountID)
	if err != nil {
		log.Errorf("GetUploadToken error: %v", err)
		return nil, err
	}
	return policy, nil
}
