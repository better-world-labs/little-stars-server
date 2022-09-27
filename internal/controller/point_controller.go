package controller

import (
	"aed-api-server/internal/controller/dto"
	"aed-api-server/internal/interfaces"
	"aed-api-server/internal/interfaces/entities"
	service2 "aed-api-server/internal/interfaces/service"
	"aed-api-server/internal/pkg"
	page "aed-api-server/internal/pkg/query"
	"aed-api-server/internal/pkg/response"
	"aed-api-server/internal/pkg/utils"
	"github.com/gin-gonic/gin"
	"gitlab.openviewtech.com/openview-pub/gopkg/route"
)

type PointsController struct {
	Service service2.PointsService `inject:"-"`
	User    service2.UserService   `inject:"-"`
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

func (con PointsController) MountAdminRouter(r *route.Router) {
	g := r.Group("/points")
	g.POST("/award", con.PointsAward)
	g.POST("/punish", con.PointsPunish)
	g.GET("/award-flow", con.GetAwardFlow)
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

func (con PointsController) PointsAward(ctx *gin.Context) (interface{}, error) {
	param := struct {
		UserId      int64  `json:"userId"`
		Points      int    `json:"points"`
		Description string `json:"description"`
		AutoReceive bool   `json:"autoReceive"`
	}{}

	err := ctx.ShouldBindJSON(&param)
	if err != nil {
		return nil, err
	}

	if param.Points <= 0 {
		return nil, response.NewIllegalArgumentError("invalid point")
	}

	return nil, con.Service.DealPoint(param.UserId, param.Points, param.Description, entities.PointsEventTypeActivityGive, param.AutoReceive)
}

func (con PointsController) PointsPunish(ctx *gin.Context) (interface{}, error) {
	param := struct {
		UserId      int64  `json:"userId"`
		Points      int    `json:"points"`
		Description string `json:"description"`
	}{}

	err := ctx.ShouldBindJSON(&param)
	if err != nil {
		return nil, err
	}

	if param.Points <= 0 {
		return nil, response.NewIllegalArgumentError("invalid point")
	}

	return nil, con.Service.DealPoint(param.UserId, -param.Points, param.Description, entities.PointsEventTypeActivityGive, true)
}

func (con PointsController) GetAwardFlow(ctx *gin.Context) (interface{}, error) {
	p, err := page.BindPageQuery(ctx)
	if err != nil {
		return nil, err
	}

	var query entities.AwardFlowQueryCommand
	err = ctx.ShouldBindQuery(&query)
	if err != nil {
		return nil, err
	}

	flows, err := con.Service.PageAwardPointFLows(*p, query)
	if err != nil {
		return nil, err
	}

	userIds := utils.Map(flows.List, func(f *entities.AwardPointFlow) int64 {
		return f.UserId
	})

	userMap, err := con.User.GetMapUserByIDs(userIds)
	if err != nil {
		return nil, err
	}

	return page.NewResult[*dto.AwardFlowDto](dto.ParseAwardFlowDtos(flows.List, userMap), flows.Total), nil
}
