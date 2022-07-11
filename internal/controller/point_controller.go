package controller

import (
	"aed-api-server/internal/controller/dto"
	"aed-api-server/internal/interfaces"
	service2 "aed-api-server/internal/interfaces/service"
	"aed-api-server/internal/pkg"
	"aed-api-server/internal/pkg/utils"
	"github.com/gin-gonic/gin"
	"gitlab.openviewtech.com/openview-pub/gopkg/route"
)

type PointsController struct {
	Service service2.PointsService `inject:"-"`
}

//go:inject-component
func NewPointsController() *PointsController {
	return &PointsController{}
}

func (con PointsController) MountAuthRouter(r *route.Router) {
	g := r.Group("/points")
	g.GET("/details", con.GetPointDetail)
	g.GET("/total", con.GetUserPointsCount)
	g.GET("/strategies", con.GetPointsStrategies)
}

func (PointsController) GetPointsStrategies(c *gin.Context) (interface{}, error) {
	strategies, err := interfaces.S.PointsScheduler.GetPointStrategies()
	if err != nil {
		return nil, err
	}
	return strategies, nil
}

func (PointsController) GetPointDetail(c *gin.Context) (interface{}, error) {
	userId := c.MustGet(pkg.AccountIDKey).(int64)
	records, err := interfaces.S.Points.GetUserPointsRecords(userId)
	if err != nil {
		return nil, err
	}
	return struct {
		Details []*dto.UserPointsRecordDto `json:"details"`
	}{dto.ParseDtos(records)}, nil
}

func (PointsController) GetUserPointsCount(c *gin.Context) (interface{}, error) {
	userId := c.MustGet(pkg.AccountIDKey).(int64)
	all, err := utils.PromiseAll(func() (interface{}, error) {
		return interfaces.S.Points.GetUserTotalPoints(userId)
	}, func() (interface{}, error) {
		return interfaces.S.Points.GetUserDonatedPoints(userId)
	})

	if err != nil {
		return nil, err
	}

	total := all[0].(int)
	donated := all[1].(int)

	return struct {
		Total   int `json:"total"`
		Donated int `json:"donated"`
	}{total, donated}, nil
}
