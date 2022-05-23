package project

import (
	"aed-api-server/internal/interfaces"
	"aed-api-server/internal/pkg"
	"aed-api-server/internal/pkg/response"
	_ "aed-api-server/internal/service/project" //关联service
	"github.com/gin-gonic/gin"
)

type Controller struct {
}

type Param struct {
	ProjectId int64 `uri:"projectId"`
	CourseId  int64 `uri:"courseId"`
	ArticleId int64 `uri:"articleId"`
}

func (Controller) GetProjectById(c *gin.Context) {
	var param Param
	if err := c.ShouldBindUri(&param); err != nil {
		response.ReplyError(c, err)
		return
	}

	project, err := interfaces.S.Project.GetProjectById(param.ProjectId)
	if err != nil {
		response.ReplyError(c, err)
		return
	}
	response.ReplyOK(c, project)
}

func (Controller) CheckVideoCompleted(c *gin.Context) {
	type CheckVideoCompleted struct {
		Completed bool `json:"completed"`
	}
	var param Param
	if err := c.ShouldBindUri(&param); err != nil {
		response.ReplyError(c, err)
		return
	}

	completed, err := interfaces.S.Project.IsProjectVideoCompleted(param.ProjectId, c.MustGet(pkg.AccountIDKey).(int64))
	if err != nil {
		response.ReplyError(c, err)
		return
	}

	response.ReplyOK(c, &CheckVideoCompleted{
		Completed: completed,
	})
}

func (Controller) CompletedProjectVideo(c *gin.Context) {
	var param Param
	if err := c.ShouldBindUri(&param); err != nil {
		response.ReplyError(c, err)
		return
	}

	rst, err := interfaces.S.Project.CompletedProjectVideo(param.ProjectId, c.MustGet(pkg.AccountIDKey).(int64))
	if err != nil {
		response.ReplyError(c, err)
		return
	}
	response.ReplyOK(c, rst)
}

func (Controller) GetProjectUserLevel(c *gin.Context) {
	var param Param
	if err := c.ShouldBindUri(&param); err != nil {
		response.ReplyError(c, err)
		return
	}

	level, err := interfaces.S.Project.GetUserProjectLevel(param.ProjectId, c.MustGet(pkg.AccountIDKey).(int64))
	if err != nil {
		response.ReplyError(c, err)
		return
	}
	response.ReplyOK(c, level)
}
