package controller

import (
	"aed-api-server/internal/interfaces"
	"aed-api-server/internal/interfaces/entities"
	"aed-api-server/internal/pkg"
	"aed-api-server/internal/pkg/response"
	"github.com/gin-gonic/gin"
	"io"
	"log"
	"time"
)

type FeedbackController struct{}

func (FeedbackController) SubmitFeedback(c *gin.Context) {
	var feedback entities.Feedback

	if err := c.ShouldBindJSON(&feedback); err != nil {
		response.ReplyError(c, err)
		return
	}
	userId := c.MustGet(pkg.AccountIDKey).(int64)
	err := interfaces.S.Feedback.SubmitFeedback(userId, &feedback)
	if err != nil {
		response.ReplyError(c, err)
		return
	}

	response.ReplyOK(c, nil)
}

func (FeedbackController) ExportFeedback(c *gin.Context) {
	type Params struct {
		BeginDate time.Time `form:"beginDate"`
		EndDate   time.Time `form:"endDate"`
	}

	var params Params

	if err := c.ShouldBindQuery(&params); err != nil {
		response.ReplyError(c, err)
		return
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
}
