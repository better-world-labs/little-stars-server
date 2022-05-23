package merit_tree

import (
	"aed-api-server/internal/interfaces"
	"aed-api-server/internal/interfaces/service"
	"aed-api-server/internal/pkg"
	"aed-api-server/internal/pkg/global"
	"aed-api-server/internal/pkg/response"
	"github.com/gin-gonic/gin"
	"time"
)

type Controller struct {
}

func (Controller) GetUserMeritTreeInfo(c *gin.Context) {
	info, err := interfaces.S.MeritTree.GetTreeInfo(c.MustGet(pkg.AccountIDKey).(int64))
	if err != nil {
		response.ReplyError(c, err)
		return
	}

	response.ReplyOK(c, info)
}

func (Controller) GetUserMeritTreeBubblesCount(c *gin.Context) {
	type Rst struct {
		Count int `json:"count"`
	}

	count, err := interfaces.S.MeritTree.GetTreeBubblesCount(c.MustGet(pkg.AccountIDKey).(int64))
	if err != nil {
		response.ReplyError(c, err)
		return
	}

	response.ReplyOK(c, &Rst{
		count,
	})
}

func (Controller) ReadUserMeritTreeBubblesCount(c *gin.Context) {
	err := interfaces.S.MeritTree.ReadTreeBubblesCount(c.MustGet(pkg.AccountIDKey).(int64))
	if err != nil {
		response.ReplyError(c, err)
		return
	}
	response.ReplyOK(c, nil)
}

func (Controller) AcceptBubblePoints(c *gin.Context) {
	type AcceptBubblePointsParams struct {
		BubbleId int64 `form:"bubbleId"`
	}

	var params AcceptBubblePointsParams
	if err := c.ShouldBindJSON(&params); err != nil {
		response.ReplyError(c, err)
		return
	}

	rst, err := interfaces.S.MeritTree.ReceiveBubblePoints(c.MustGet(pkg.AccountIDKey).(int64), params.BubbleId)
	if err != nil {
		response.ReplyError(c, err)
		return
	}
	response.ReplyOK(c, rst)
}

func (Controller) GetWalkConvertInfo(c *gin.Context) {
	var req service.WalkConvertInfoReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ReplyError(c, err)
		return
	}
	userId := c.MustGet(pkg.AccountIDKey).(int64)
	info, err := interfaces.S.Walk.GetWalkConvertInfo(userId, &req)
	if err != nil {
		response.ReplyError(c, err)
		return
	}
	response.ReplyOK(c, info)
}

func (Controller) ConvertWalkToPoints(c *gin.Context) {
	type Req struct {
		TodayWalk int `json:"todayWalk" form:"todayWalk"`
	}
	var req Req
	userId := c.MustGet(pkg.AccountIDKey).(int64)

	if err := c.ShouldBindJSON(&req); err != nil {
		response.ReplyError(c, err)
		return
	}

	rst, err := interfaces.S.Walk.ConvertWalkToPoints(userId, req.TodayWalk)
	if err != nil {
		response.ReplyError(c, err)
		return
	}
	response.ReplyOK(c, rst)
}

func (Controller) SignEarly(c *gin.Context) {
	userId := c.MustGet(pkg.AccountIDKey).(int64)

	rst, bErr := interfaces.S.Early.SignEarly(userId)
	if bErr != nil {
		if bErr != response.ErrorSignEarlyTimeNotAllowed &&
			bErr != response.ErrorSignEarlyTodayAlreadySignedYet {
			response.ReplyError(c, bErr)
			return
		}

		record, exists, err := interfaces.S.Early.GetLatestRecord(userId)
		if err != nil {
			response.ReplyError(c, err)
			return
		}

		errResponse := map[string]interface{}{
			"currentTime": global.FormattedTime(time.Now()),
			"days":        0,
		}

		if exists {
			errResponse["days"] = record.Days
		}

		response.ReplyErrorWithData(c, bErr, errResponse)
		return
	}

	response.ReplyOK(c, rst)
}
