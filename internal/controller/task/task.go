package task

import (
	"aed-api-server/internal/interfaces"
	"aed-api-server/internal/pkg"
	"aed-api-server/internal/pkg/location"
	"aed-api-server/internal/pkg/response"
	"github.com/gin-gonic/gin"
)

type Controller struct{}

//GetUserTasks 获取任务列表
func (Controller) GetUserTasks(c *gin.Context) {
	type Param struct {
		PageSize       int     `form:"pageSize"`
		Status         string  `form:"status"`
		Cursor         string  `form:"cursor"`
		IncludeExpired bool    `form:"includeExpired"`
		Longitude      float64 `form:"longitude"`
		Latitude       float64 `form:"latitude"`
	}

	var param Param
	err := c.ShouldBindQuery(&param)
	if err != nil {
		response.ReplyError(c, err)
		return
	}

	result, err := interfaces.S.Task.GetUserTasks(
		c.MustGet(pkg.AccountIDKey).(int64),
		param.PageSize,
		param.Status,
		param.Cursor,
		param.IncludeExpired,
		location.Coordinate{
			Longitude: param.Longitude,
			Latitude:  param.Latitude,
		},
	)
	if err != nil {
		response.ReplyError(c, err)
		return
	}
	response.ReplyOK(c, result)
}

//ReadTask 将任务设置为已读
func (Controller) ReadTask(c *gin.Context) {
	err := interfaces.S.Task.ReadUserTask(
		c.MustGet(pkg.AccountIDKey).(int64),
		c.Param("jobId"),
	)

	if err != nil {
		response.ReplyError(c, err)
		return
	}

	response.ReplyOK(c, nil)
}

func (Controller) GetUserTaskStat(c *gin.Context) {
	res, err := interfaces.S.Task.GetUserTaskStat(c.MustGet(pkg.AccountIDKey).(int64))
	if err != nil {
		response.ReplyError(c, err)
		return
	}
	response.ReplyOK(c, res)
}

func (Controller) FindPicketTaskByDeviceId(c *gin.Context) {
	userTask, err := interfaces.S.Task.FindUserTaskByUserIdAndDeviceId(
		c.MustGet(pkg.AccountIDKey).(int64),
		c.Param("deviceId"),
	)
	if err != nil {
		response.ReplyError(c, err)
		return
	}
	response.ReplyOK(c, userTask)
}
