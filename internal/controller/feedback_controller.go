package controller

import (
	"aed-api-server/internal/interfaces"
	"aed-api-server/internal/interfaces/entities"
	"aed-api-server/internal/pkg"
	"github.com/gin-gonic/gin"
	"gitlab.openviewtech.com/openview-pub/gopkg/route"
	"io"
	"log"
	"time"
)

type FeedbackController struct{}

//go:inject-component
func NewFeedbackController() *FeedbackController {
	return &FeedbackController{}
}

func (c FeedbackController) MountAuthRouter(r *route.Router) {
	r.Group("/user-feedbacks").
		POST("/", c.SubmitFeedback)
}

func (c FeedbackController) MountAdminRouter(r *route.Router) {
	r.GET("/user-feedbacks/excel", c.ExportFeedback)
}

func (FeedbackController) SubmitFeedback(c *gin.Context) (interface{}, error) {
	var feedback entities.Feedback

	if err := c.ShouldBindJSON(&feedback); err != nil {
		return nil, err
	}
	userId := c.MustGet(pkg.AccountIDKey).(int64)
	err := interfaces.S.Feedback.SubmitFeedback(userId, &feedback)
	if err != nil {
		return nil, err
	}

	return nil, nil
}

func (FeedbackController) ExportFeedback(c *gin.Context) (interface{}, error) {
	type Params struct {
		BeginDate time.Time `form:"beginDate"`
		EndDate   time.Time `form:"endDate"`
	}

	var params Params

	if err := c.ShouldBindQuery(&params); err != nil {
		return nil, nil
	}

	fileName := "用户反馈导出.xlsx"
	c.Header("Content-Type", "application/octet-stream")
	c.Header("Content-Disposition", "attachment; filename="+fileName)
	c.Header("Content-Transfer-Encoding", "binary")
	c.Header("Cache-Control", "no-cache")

	reader, writer := io.Pipe()
	interfaces.S.Feedback.ExportFeedback(params.BeginDate, params.EndDate, writer)

	if _, err := io.Copy(c.Writer, reader); err != nil {
		log.Fatal(err)
	}

	return nil, nil
}
