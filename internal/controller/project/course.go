package project

import (
	"aed-api-server/internal/interfaces"
	"aed-api-server/internal/interfaces/service"
	"aed-api-server/internal/pkg"
	"aed-api-server/internal/pkg/response"
	"github.com/gin-gonic/gin"
)

func (Controller) GetProjectCourses(c *gin.Context) {
	type GetProjectCoursesVO struct {
		Courses []*service.Course `json:"courses"`
	}
	var param Param
	if err := c.ShouldBindUri(&param); err != nil {
		response.ReplyError(c, err)
	}

	courses, err := interfaces.S.Course.GetCoursesByProjectId(param.ProjectId)
	if err != nil {
		response.ReplyError(c, err)
		return
	}

	response.ReplyOK(c, &GetProjectCoursesVO{
		Courses: courses,
	})
}

func (Controller) GetLearntCourses(c *gin.Context) {
	type GetLearntCoursesVO struct {
		Courses []*service.Course `json:"courses"`
	}
	courses, err := interfaces.S.Course.GetUserLearntCourses(c.MustGet(pkg.AccountIDKey).(int64))
	if err != nil {
		response.ReplyError(c, err)
		return
	}
	response.ReplyOK(c, &GetLearntCoursesVO{
		Courses: courses,
	})
}

func (Controller) GetProjectCoursesById(c *gin.Context) {
	var param Param
	if err := c.ShouldBindUri(&param); err != nil {
		response.ReplyError(c, err)
		return
	}

	course, err := interfaces.S.Course.GetCourseByCourseId(param.CourseId)
	if err != nil {
		response.ReplyError(c, err)
		return
	}

	response.ReplyOK(c, course)
}

func (Controller) LearntCourseById(c *gin.Context) {
	var param Param
	if err := c.ShouldBindUri(&param); err != nil {
		response.ReplyError(c, err)
		return
	}

	rst, err := interfaces.S.Course.LearntCourseByCourseId(param.CourseId, c.MustGet(pkg.AccountIDKey).(int64))
	if err != nil {
		response.ReplyError(c, err)
		return
	}
	response.ReplyOK(c, rst)
}

func (Controller) GetArticleById(c *gin.Context) {
	var param Param
	if err := c.ShouldBindUri(&param); err != nil {
		response.ReplyError(c, err)
		return
	}

	article, err := interfaces.S.Course.GetArticleById(param.ArticleId)
	if err != nil {
		response.ReplyError(c, err)
		return
	}
	response.ReplyOK(c, article)
}
