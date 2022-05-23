package clock_in

import (
	"aed-api-server/internal/interfaces"
	"aed-api-server/internal/interfaces/entities"
	"aed-api-server/internal/pkg"
	"aed-api-server/internal/pkg/feedback"
	"aed-api-server/internal/pkg/location"
	"aed-api-server/internal/pkg/response"
	"errors"
	"github.com/gin-gonic/gin"
)

type Controller struct{}

func (Controller) GetClockInStat(c *gin.Context) {
	stat, err := interfaces.S.ClockIn.GetDeviceClockInStat()
	if err != nil {
		response.ReplyError(c, err)
		return
	}
	response.ReplyOK(c, stat)
}

func (Controller) PostClockIn(c *gin.Context) {
	var clockInBaseInfo entities.ClockInBaseInfo
	if err := c.ShouldBind(&clockInBaseInfo); err != nil {
		response.ReplyError(c, err)
		return
	}
	rst, err := interfaces.S.ClockIn.DoDeviceClockIn(&clockInBaseInfo, c.MustGet(pkg.AccountIDKey).(int64))
	if err != nil {
		response.ReplyError(c, err)
		return
	}
	back := feedback.NewValuableFeedBack()
	err = back.AddPointsEventRsts(rst)
	if err != nil {
		response.ReplyError(c, err)
		return
	}

	response.ReplyOK(c, back)
}

func (Controller) GetDeviceClockInList(c *gin.Context) {
	deviceId, _ := c.GetQuery("deviceId")
	if len(deviceId) == 0 {
		response.ReplyError(c, errors.New("deviceId is required"))
		return
	}

	list, err := interfaces.S.ClockIn.GetDeviceClockInList(deviceId)
	if err != nil {
		response.ReplyError(c, err)
		return
	}
	response.ReplyOK(c, list)
}

func (Controller) GetDeviceLastClockIn(c *gin.Context) {
	deviceId, _ := c.GetQuery("deviceId")
	if len(deviceId) == 0 {
		response.ReplyError(c, errors.New("deviceId is required"))
		return
	}

	var from location.Coordinate
	err := c.ShouldBindQuery(&from)
	if err != nil {
		response.ReplyError(c, err)
		return
	}

	rst, err := interfaces.S.ClockIn.GetDeviceLastClockIn(from, deviceId)
	if err != nil {
		response.ReplyError(c, err)
		return
	}
	response.ReplyOK(c, rst)
}
