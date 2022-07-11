package controller

import (
	"aed-api-server/internal/interfaces"
	"github.com/gin-gonic/gin"
	"gitlab.openviewtech.com/openview-pub/gopkg/route"
)

type StatController struct{}

//go:inject-component
func NewStatController() *StatController {
	return &StatController{}
}

func (c StatController) MountAdminRouter(r *route.Router) {
	r.GET("/stat/kpi", c.KpiStat)
	r.GET("/stat/points/top", c.PointTop)
}

func (StatController) KpiStat(c *gin.Context) (interface{}, error) {
	stat, err := interfaces.S.Stat.DoKipStat()
	if err != nil {
		return nil, err
	}
	return stat, nil
}

func (StatController) PointTop(c *gin.Context) (interface{}, error) {
	stat, err := interfaces.S.Stat.StatPointsTop()
	if err != nil {
		return nil, err
	}
	return stat, nil
}
