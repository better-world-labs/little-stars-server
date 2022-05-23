package skill

import (
	"aed-api-server/internal/pkg"
	"aed-api-server/internal/pkg/response"
	"aed-api-server/internal/pkg/utils"
	"github.com/gin-gonic/gin"
	"gitlab.openviewtech.com/openview-pub/gopkg/log"
)

type Controller struct {
	service Service
}

func NewController(s Service) *Controller {
	return &Controller{service: s}
}

func (con *Controller) MyCertificate(c *gin.Context) {
	accountID := c.MustGet(pkg.AccountIDKey).(int64)
	list := con.service.MyCertificate(accountID)

	response.ReplyOK(c, list)
}

func (con *Controller) ListCerts(c *gin.Context) {
	certs := con.service.ListCerts()
	response.ReplyOK(c, certs)
}

func (con *Controller) CreateCertEvidences(c *gin.Context) {
	log.DefaultLogger().Infof("create certs ...")
	err := con.service.CreateEvidences()
	utils.MustNil(err, err)
	response.ReplyOK(c, nil)
}
