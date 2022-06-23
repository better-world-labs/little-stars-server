package project

import (
	"aed-api-server/internal/interfaces"
	"aed-api-server/internal/interfaces/service"
	"aed-api-server/internal/pkg"
	"github.com/gin-gonic/gin"
)

func (Controller) GetProjectCourses(c *gin.Context) (interface{}, error) {
	type GetProjectCoursesVO struct {
		Courses []*service.Course `json:"courses"`
	}
	var param Param
	if err := c.ShouldBindUri(&param); err != nil {
		return nil, err
	}

	courses, err := interfaces.S.Course.GetCoursesByProjectId(param.ProjectId)
	if err != nil {
		return nil, err
	}

	return &GetProjectCoursesVO{
		Courses: courses,
	}, nil
}

func (Controller) GetLearntCourses(c *gin.Context) (interface{}, error) {
	type GetLearntCoursesVO struct {
		Courses []*service.Course `json:"courses"`
	}
	courses, err := interfaces.S.Course.GetUserLearntCourses(c.MustGet(pkg.AccountIDKey).(int64))
	if err != nil {
		return nil, err
	}

	return &GetLearntCoursesVO{
		Courses: courses,
	}, nil
}

func (Controller) GetProjectCoursesById(c *gin.Context) (interface{}, error) {
	var param Param
	if err := c.ShouldBindUri(&param); err != nil {
		return nil, err
	}

	course, err := interfaces.S.Course.GetCourseByCourseId(param.CourseId)
	if err != nil {
		return nil, err
	}

	return course, nil
}

func (Controller) LearntCourseById(c *gin.Context) (interface{}, error) {
	var param Param
	if err := c.ShouldBindUri(&param); err != nil {
		return nil, err
	}

	rst, err := interfaces.S.Course.LearntCourseByCourseId(param.CourseId, c.MustGet(pkg.AccountIDKey).(int64))
	if err != nil {
		return nil, err
	}

	return rst, nil
}

func (Controller) GetArticleById(c *gin.Context) (interface{}, error) {
	var param Param
	if err := c.ShouldBindUri(&param); err != nil {
		return nil, err
	}

	article, err := interfaces.S.Course.GetArticleById(param.ArticleId)
	if err != nil {
		return nil, err
	}

	return article, nil
}
