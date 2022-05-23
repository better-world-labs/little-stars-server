package point

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

func (con *Controller) QueryPointDetail(c *gin.Context) {
	accountID := c.MustGet(pkg.AccountIDKey).(int64)
	res, err := con.service.Detail(accountID)
	if err != nil {
		log.DefaultLogger().Errorf("GetPointDetail error: %v", err)
		response.ReplyError(c, err)
		return
	}

	response.ReplyOK(c, map[string]interface{}{"details": res})
}

func (con *Controller) TotalPoints(c *gin.Context) {
	accountID := c.MustGet(pkg.AccountIDKey).(int64)
	res, err := con.service.TotalPoints(accountID)
	if err != nil {
		log.DefaultLogger().Errorf("TotalPoints error: %v", err)
		response.ReplyError(c, err)
		return
	}

	response.ReplyOK(c, map[string]interface{}{"total": res})
}
