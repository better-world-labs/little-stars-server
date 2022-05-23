package exam

import (
	"aed-api-server/internal/interfaces"
	"aed-api-server/internal/interfaces/service"
	"aed-api-server/internal/pkg"
	"aed-api-server/internal/pkg/feedback"
	"aed-api-server/internal/pkg/response"
	"github.com/gin-gonic/gin"
	"strconv"
)

type (
	ParamStart struct {
		Type int `json:"type" binding:"required,min=1,max=2"`
	}
	ParamAnswer struct {
		ID      int64 `json:"id"`
		Answers []int `json:"answers"`
	}
)

type Controller struct {
	service service.ExamService
}

func NewController(cert service.CertService) *Controller {
	// TODO 临时方式，后边修改依赖编排方式
	interfaces.S.Exam.SetCertService(cert)
	return &Controller{service: interfaces.S.Exam}
}

func (con *Controller) Start(ctx *gin.Context) {
	userId := ctx.MustGet(pkg.AccountIDKey).(int64)
	var p ParamStart
	projectId := ctx.Param("projectId")

	err := ctx.BindJSON(&p)
	if err != nil {
		response.ReplyError(ctx, err)
		return
	}

	projectIdInt64, err := strconv.ParseInt(projectId, 10, 64)
	if err != nil {
		response.ReplyError(ctx, err)
		return
	}

	exam, err := con.service.Start(projectIdInt64, userId, p.Type)
	if err != nil {
		response.ReplyError(ctx, err)
		return
	}

	response.ReplyOK(ctx, NewStartExamVo(*exam))
}

func (con *Controller) Save(ctx *gin.Context) {
	userId := ctx.MustGet(pkg.AccountIDKey).(int64)
	examId := ctx.Param("examId")
	var p struct {
		Questions []ParamAnswer `json:"questions"`
	}
	err := ctx.BindJSON(&p)
	if err != nil {
		response.ReplyError(ctx, err)
		return
	}

	examIdInt, err := strconv.ParseInt(examId, 10, 64)
	if err != nil {
		response.ReplyError(ctx, err)
		return
	}

	var paper = map[int64][]int{}
	for _, q := range p.Questions {
		paper[q.ID] = q.Answers
	}

	err = con.service.Save(examIdInt, userId, paper)
	if err != nil {
		response.ReplyError(ctx, err)
		return
	}

	response.ReplyOK(ctx, nil)
}

func (con *Controller) Submit(ctx *gin.Context) {
	userId := ctx.MustGet(pkg.AccountIDKey).(int64)
	examId := ctx.Param("examId")
	var p struct {
		Questions []ParamAnswer `json:"questions"`
	}
	err := ctx.BindJSON(&p)
	if err != nil {
		response.ReplyError(ctx, err)
		return
	}

	examIdInt, err := strconv.ParseInt(examId, 10, 64)
	if err != nil {
		response.ReplyError(ctx, err)
		return
	}

	var paper = map[int64][]int{}
	for _, q := range p.Questions {
		paper[q.ID] = q.Answers
	}

	certImg, certNum, pointsRst, err := con.service.Submit(examIdInt, userId, paper)
	if err != nil {
		response.ReplyError(ctx, err)
		return
	}

	exam, exists, err := con.service.GetByID(examIdInt)
	if err != nil || !exists {
		response.ReplyError(ctx, err)
		return
	}

	back := feedback.NewValuableFeedBack()
	err = back.AddPointsEventRsts(pointsRst)
	if err != nil || !exists {
		response.ReplyError(ctx, err)
		return
	}

	back.Put("score", exam.Score)
	back.Put("passed", exam.CheckPass())
	back.Put("certId", certNum)
	back.Put("certImg", certImg)

	response.ReplyOK(ctx, back)
}

func (con *Controller) ListSubmitted(ctx *gin.Context) {
	userId := ctx.MustGet(pkg.AccountIDKey).(int64)
	query := struct {
		Latest int `form:"latest"`
		Type   int `form:"type" binding:"required,min=1,max=2"`
	}{}
	projectIdStr := ctx.Param("projectId")

	projectId, err := strconv.ParseInt(projectIdStr, 10, 64)
	if err != nil {
		response.ReplyError(ctx, err)
		return
	}

	err = ctx.BindQuery(&query)
	if err != nil {
		response.ReplyError(ctx, err)
		return
	}

	exams, err := con.service.ListLatestSubmitted(projectId, userId, query.Type, query.Latest)
	if err != nil {
		response.ReplyError(ctx, err)
		return
	}

	vos := make([]SimpleListInfoDto, len(exams))
	for i, e := range exams {
		vos[i] = ParseSimpleListInfoDto(e)
	}

	response.ReplyOK(ctx, map[string]interface{}{
		"exams": vos,
	})
}

func (con *Controller) GetUnSubmittedLatest(ctx *gin.Context) {
	userId := ctx.MustGet(pkg.AccountIDKey).(int64)
	var query struct {
		Type int `form:"type" binding:"required,min=1,max=2"`
	}

	projectIdStr := ctx.Param("projectId")
	projectId, err := strconv.ParseInt(projectIdStr, 10, 64)
	if err != nil {
		response.ReplyError(ctx, err)
		return
	}
	err = ctx.BindQuery(&query)
	if err != nil {
		response.ReplyError(ctx, err)
		return
	}

	exam, exists, err := con.service.GetLatestUnSubmitted(projectId, userId, query.Type)
	if err != nil {
		response.ReplyError(ctx, err)
		return
	}

	if !exists {
		response.ReplyOK(ctx, nil)
		return
	}

	response.ReplyOK(ctx, ParseSimpleListInfoWithQuestionDto(exam))
}

func (con *Controller) GetByID(ctx *gin.Context) {
	idStr := ctx.Param("examId")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		response.ReplyError(ctx, err)
		return
	}

	exam, exists, err := con.service.GetByID(id)
	if err != nil {
		response.ReplyError(ctx, err)
		return
	}

	if !exists {
		response.ReplyOK(ctx, nil)
		return
	}

	response.ReplyOK(ctx, ParseDetailDTo(exam))
}
