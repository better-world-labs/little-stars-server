package controller

import (
	service2 "aed-api-server/internal/interfaces/service"
	"aed-api-server/internal/pkg"
	"aed-api-server/internal/pkg/utils"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"gitlab.openviewtech.com/openview-pub/gopkg/route"
)

type SkillController struct {
	Service service2.SkillService `inject:"-"`
}

func NewSkillController() *SkillController {
	return &SkillController{}
}

func (con *SkillController) MountNoAuthRouter(r *route.Router) {
	// 为旧数据生成存证
	r.GET("/skill/create-cert-evidences", con.CreateCertEvidences)
}

func (con *SkillController) MountAuthRouter(r *route.Router) {
	r.GET("/skill/cert", con.MyCertificate)
	r.GET("/skill/certs", con.ListCerts)
}

func (con *SkillController) MyCertificate(c *gin.Context) (interface{}, error) {
	accountID := c.MustGet(pkg.AccountIDKey).(int64)
	list := con.Service.MyCertificate(accountID)
	return list, nil
}

func (con *SkillController) ListCerts(c *gin.Context) (interface{}, error) {
	certs := con.Service.ListCerts()
	return certs, nil
}

func (con *SkillController) CreateCertEvidences(c *gin.Context) (interface{}, error) {
	log.Infof("create certs ...")
	err := con.Service.CreateEvidences()
	utils.MustNil(err, err)
	return nil, nil
}
