package controller

import (
	"aed-api-server/internal/interfaces"
	"aed-api-server/internal/interfaces/service"
	"aed-api-server/internal/pkg"
	"aed-api-server/internal/pkg/location"
	"github.com/gin-gonic/gin"
	"gitlab.openviewtech.com/openview-pub/gopkg/route"
)

type TaskController struct {
	TaskService service.TaskService `inject:"-"`
}

//go:inject-component
func NewTaskController() *TaskController {
	return &TaskController{}
}

func (c TaskController) MountAuthRouter(r *route.Router) {
	taskR := r.Group("/task-jobs")
	r.GET("/devices/:deviceId/picket-task", c.FindPicketTaskByDeviceId)
	taskR.GET("", c.GetUserTasks)
	taskR.GET("/count", c.GetUserTaskStat)
	taskR.PUT("/:jobId/read", c.ReadTask)
	taskR.GET("/job", c.findByLink)
}

//GetUserTasks 获取任务列表
func (TaskController) GetUserTasks(c *gin.Context) (interface{}, error) {
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
		return nil, err
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
		return nil, err
	}

	return result, nil
}

//ReadTask 将任务设置为已读
func (TaskController) ReadTask(c *gin.Context) (interface{}, error) {
	err := interfaces.S.Task.ReadUserTask(
		c.MustGet(pkg.AccountIDKey).(int64),
		c.Param("jobId"),
	)

	if err != nil {
		return nil, err
	}

	return nil, nil
}

func (TaskController) GetUserTaskStat(c *gin.Context) (interface{}, error) {
	res, err := interfaces.S.Task.GetUserTaskStat(c.MustGet(pkg.AccountIDKey).(int64))
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (TaskController) FindPicketTaskByDeviceId(c *gin.Context) (interface{}, error) {
	userTask, err := interfaces.S.Task.FindUserTaskByUserIdAndDeviceId(
		c.MustGet(pkg.AccountIDKey).(int64),
		c.Param("deviceId"),
	)
	if err != nil {
		return nil, err
	}
	return userTask, nil
}

func (c TaskController) findByLink(context *gin.Context) (interface{}, error) {
	link := context.Query("link")
	if link == "" {
		return nil, nil
	}
	userId := context.MustGet(pkg.AccountIDKey).(int64)
	job, err := c.TaskService.FindJobByPageLink(userId, link)
	if err != nil {
		return nil, err
	}
	if job == nil {
		return nil, nil
	}
	return map[string]interface{}{
		"taskId": job.TaskId,
		"points": job.Points,
	}, nil
}
