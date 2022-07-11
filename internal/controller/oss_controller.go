package controller

import (
	"aed-api-server/internal/interfaces/service"
	"aed-api-server/internal/pkg"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"gitlab.openviewtech.com/openview-pub/gopkg/route"
)

type OssController struct {
	Oss service.OssService `inject:"-"`
}

//go:inject-component
func NewOssController() *OssController {
	return &OssController{}
}

func (con *OssController) MountAdminRouter(r *route.Router) {
	r.GET("/upload-token", con.GetUploadTokenAdmin)
}

func (con *OssController) MountAuthRouter(r *route.Router) {
	r.POST("/common/photo", con.Upload)
	r.GET("/common/upload_token", con.GetUploadToken)
}

// 应用服务器上传
func (con *OssController) Upload(c *gin.Context) (interface{}, error) {
	fileHeader, err := c.FormFile("files")
	if err != nil {
		log.Errorf("get file error: %v", err)
		return nil, err
	}
	url, err := con.Oss.OssUpload(fileHeader)
	if err != nil {
		log.Errorf("OssUpload error: %v", err)
		return nil, err
	}

	return map[string]interface{}{
		"Origin": url,
	}, nil
}

// 获取阿里云直传Token
func (con *OssController) GetUploadToken(c *gin.Context) (interface{}, error) {
	accountID := c.MustGet(pkg.AccountIDKey).(int64)
	policy, err := con.Oss.GetUploadToken("", accountID)
	if err != nil {
		log.Errorf("GetUploadToken error: %v", err)
		return nil, err
	}
	return policy, nil
}

func (con *OssController) GetUploadTokenAdmin(context *gin.Context) (interface{}, error) {
	accountID := context.MustGet(pkg.AccountIDKey).(int64)
	if accountID < 0 {
		accountID = -accountID
	}
	policy, err := con.Oss.GetUploadToken("stars-admin/", accountID)
	if err != nil {
		log.Errorf("GetUploadToken error: %v", err)
		return nil, err
	}
	return policy, nil
}
