package project

import (
	"aed-api-server/internal/interfaces"
	"aed-api-server/internal/pkg"
	_ "aed-api-server/internal/service/project" //关联service
	"github.com/gin-gonic/gin"
	"gitlab.openviewtech.com/openview-pub/gopkg/route"
)

type Controller struct {
}

type Param struct {
	ProjectId int64 `uri:"projectId"`
	CourseId  int64 `uri:"courseId"`
	ArticleId int64 `uri:"articleId"`
}

func NewController() *Controller {
	return &Controller{}
}

func (c Controller) MountAuthRouter(r *route.Router) {
	projectC := r.Group("/projects")
	projectC.GET("/:projectId/check-video-completed", c.CheckVideoCompleted)
	projectC.GET("/:projectId/courses/learnt", c.GetLearntCourses)
	projectC.PUT("/:projectId/video/completed", c.CompletedProjectVideo)
	projectC.PUT("/courses/:courseId/learnt", c.LearntCourseById)
	projectC.GET("/:projectId/level", c.GetProjectUserLevel)
}

func (c Controller) MountNoAuthRouter(r *route.Router) {
	projectR := r.Group("/projects")
	projectR.GET("/:projectId", c.GetProjectById)
	projectR.GET("/:projectId/courses", c.GetProjectCourses)
	projectR.GET("/courses/:courseId", c.GetProjectCoursesById)
	projectR.GET("/courses/articles/:articleId", c.GetArticleById)
}

func (Controller) GetProjectById(c *gin.Context) (interface{}, error) {
	var param Param
	if err := c.ShouldBindUri(&param); err != nil {
		return nil, err
	}

	project, err := interfaces.S.Project.GetProjectById(param.ProjectId)
	if err != nil {
		return nil, err
	}

	return project, nil
}

func (Controller) CheckVideoCompleted(c *gin.Context) (interface{}, error) {
	type CheckVideoCompleted struct {
		Completed bool `json:"completed"`
	}
	var param Param
	if err := c.ShouldBindUri(&param); err != nil {
		return nil, err
	}

	completed, err := interfaces.S.Project.IsProjectVideoCompleted(param.ProjectId, c.MustGet(pkg.AccountIDKey).(int64))
	if err != nil {
		return nil, err
	}

	return &CheckVideoCompleted{
		Completed: completed,
	}, nil
}

func (Controller) CompletedProjectVideo(c *gin.Context) (interface{}, error) {
	var param Param
	if err := c.ShouldBindUri(&param); err != nil {
		return nil, err
	}

	rst, err := interfaces.S.Project.CompletedProjectVideo(param.ProjectId, c.MustGet(pkg.AccountIDKey).(int64))
	if err != nil {
		return nil, err
	}

	return rst, nil
}

func (Controller) GetProjectUserLevel(c *gin.Context) (interface{}, error) {
	var param Param
	if err := c.ShouldBindUri(&param); err != nil {
		return nil, err
	}

	level, err := interfaces.S.Project.GetUserProjectLevel(param.ProjectId, c.MustGet(pkg.AccountIDKey).(int64))
	if err != nil {
		return nil, err
	}

	return level, nil
}
