package point

import (
	"aed-api-server/internal/interfaces"
	"aed-api-server/internal/pkg"
	"aed-api-server/internal/pkg/response"
	"aed-api-server/internal/pkg/utils"
	"github.com/gin-gonic/gin"
)

type Controller struct {
}

func (Controller) GetPointsStrategies(c *gin.Context) {
	strategies, err := interfaces.S.PointsScheduler.GetPointStrategies()
	if err != nil {
		response.ReplyError(c, err)
		return
	}
	response.ReplyOK(c, strategies)
}

func (Controller) GetPointDetail(c *gin.Context) {
	userId := c.MustGet(pkg.AccountIDKey).(int64)
	records, err := interfaces.S.Points.GetUserPointsRecords(userId)
	if err != nil {
		response.ReplyError(c, err)
		return
	}
	response.ReplyOK(c, struct {
		Details []*UserPointsRecordDto `json:"details"`
	}{parseDtos(records)})
}

func (Controller) GetUserPointsCount(c *gin.Context) {
	userId := c.MustGet(pkg.AccountIDKey).(int64)
	all, err := utils.PromiseAll(func() (interface{}, error) {
		return interfaces.S.Points.GetUserTotalPoints(userId)
	}, func() (interface{}, error) {
		return interfaces.S.Points.GetUserDonatedPoints(userId)
	})

	if err != nil {
		response.ReplyError(c, err)
		return
	}

	total := all[0].(int)
	donated := all[1].(int)

	response.ReplyOK(c, struct {
		Total   int `json:"total"`
		Donated int `json:"donated"`
	}{total, donated})
}
