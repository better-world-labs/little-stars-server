package merit_tree

import (
	"aed-api-server/internal/interfaces"
	"aed-api-server/internal/interfaces/entities"
	"aed-api-server/internal/interfaces/service"
	"aed-api-server/internal/pkg"
	"aed-api-server/internal/pkg/global"
	"aed-api-server/internal/pkg/response"
	"github.com/gin-gonic/gin"
	"gitlab.openviewtech.com/openview-pub/gopkg/route"
	"time"
)

type Controller struct {
	TreasureChest service.TreasureChestService `inject:"-"`
}

func NewController() *Controller {
	return &Controller{}
}

func (c Controller) MountAdminRouter(r *route.Router) {
	g := r.Group("/merit-tree")
	g.POST("/treasure-chest", c.CreateTreasureChest)
}

func (c Controller) MountAuthRouter(r *route.Router) {
	g := r.Group("/merit-tree")
	g.GET("", c.GetUserMeritTreeInfo)
	g.GET("/bubbles/count", c.GetUserMeritTreeBubblesCount)
	g.PUT("/bubbles/count", c.ReadUserMeritTreeBubblesCount)
	g.PUT("/bubbles", c.AcceptBubblePoints)
	g.POST("/walk-points", c.GetWalkConvertInfo)
	g.PUT("/walk-points", c.ConvertWalkToPoints)
	g.PUT("/sign-early", c.SignEarly)
	g.GET("/treasure-chest", c.GetUserTreasureChest)
}

func (Controller) GetUserMeritTreeInfo(c *gin.Context) (interface{}, error) {
	info, err := interfaces.S.MeritTree.GetTreeInfo(c.MustGet(pkg.AccountIDKey).(int64))
	if err != nil {
		return nil, err
	}

	return info, nil
}

func (Controller) GetUserMeritTreeBubblesCount(c *gin.Context) (interface{}, error) {
	type Rst struct {
		Count int `json:"count"`
	}

	count, err := interfaces.S.MeritTree.GetTreeBubblesCount(c.MustGet(pkg.AccountIDKey).(int64))
	if err != nil {
		return nil, err
	}

	return &Rst{
		count,
	}, nil
}

func (Controller) ReadUserMeritTreeBubblesCount(c *gin.Context) (interface{}, error) {
	err := interfaces.S.MeritTree.ReadTreeBubblesCount(c.MustGet(pkg.AccountIDKey).(int64))
	if err != nil {
		return nil, err
	}

	return nil, nil
}

func (Controller) AcceptBubblePoints(c *gin.Context) (interface{}, error) {
	type AcceptBubblePointsParams struct {
		BubbleId int64 `form:"bubbleId"`
	}

	var params AcceptBubblePointsParams
	if err := c.ShouldBindJSON(&params); err != nil {
		return nil, err
	}

	rst, err := interfaces.S.MeritTree.ReceiveBubblePoints(c.MustGet(pkg.AccountIDKey).(int64), params.BubbleId)
	if err != nil {
		return nil, err
	}

	return rst, nil
}

func (Controller) GetWalkConvertInfo(c *gin.Context) (interface{}, error) {
	var req service.WalkConvertInfoReq
	if err := c.ShouldBindJSON(&req); err != nil {
		return nil, err
	}
	userId := c.MustGet(pkg.AccountIDKey).(int64)
	info, err := interfaces.S.Walk.GetWalkConvertInfo(userId, &req)
	if err != nil {
		return nil, err
	}

	return info, nil
}

func (Controller) ConvertWalkToPoints(c *gin.Context) (interface{}, error) {
	type Req struct {
		TodayWalk int `json:"todayWalk" form:"todayWalk"`
	}
	var req Req
	userId := c.MustGet(pkg.AccountIDKey).(int64)

	if err := c.ShouldBindJSON(&req); err != nil {
		return nil, err
	}

	rst, err := interfaces.S.Walk.ConvertWalkToPoints(userId, req.TodayWalk)
	if err != nil {
		return nil, err
	}
	return rst, nil
}

func (Controller) SignEarly(c *gin.Context) (interface{}, error) {
	userId := c.MustGet(pkg.AccountIDKey).(int64)

	rst, bErr := interfaces.S.Early.SignEarly(userId)
	if bErr != nil {
		if bErr != response.ErrorSignEarlyTimeNotAllowed &&
			bErr != response.ErrorSignEarlyTodayAlreadySignedYet {
			return nil, bErr
		}

		record, exists, err := interfaces.S.Early.GetLatestRecord(userId)
		if err != nil {
			return nil, err
		}

		errResponse := map[string]interface{}{
			"currentTime": global.FormattedTime(time.Now()),
			"days":        0,
		}

		if exists {
			errResponse["days"] = record.Days
		}

		return errResponse, bErr
	}

	return rst, nil
}

func (c Controller) GetUserTreasureChest(ctx *gin.Context) (interface{}, error) {
	userId := ctx.MustGet(pkg.AccountIDKey).(int64)
	return c.TreasureChest.GetUserTreasureChest(userId)
}

func (c Controller) CreateTreasureChest(context *gin.Context) (interface{}, error) {
	r := new(entities.TreasureChestCreateRequest)
	if err := context.ShouldBindJSON(&r); err != nil {
		return nil, err
	}

	if err := c.TreasureChest.CreateTreasureChest(r); err != nil {
		return nil, err
	}
	return nil, nil
}
