package exam

import (
	"aed-api-server/internal/interfaces/service"
	"aed-api-server/internal/pkg"
	"aed-api-server/internal/pkg/feedback"
	"github.com/gin-gonic/gin"
	"gitlab.openviewtech.com/openview-pub/gopkg/route"
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
	Service service.ExamService `inject:"-"`
}

func NewController() *Controller {
	return &Controller{}
}

func (con *Controller) MountAuthRouter(r *route.Router) {
	r.GET("/projects/:projectId/exams/submitted", con.ListSubmitted)
	r.GET("/projects/:projectId/exams/unsubmitted/latest", con.GetUnSubmittedLatest)
	r.POST("/projects/:projectId/exams", con.Start)
	r.POST("/projects/exams/:examId/save", con.Save)
	r.POST("/projects/exams/:examId/submit", con.Submit)
	r.GET("/projects/exams/:examId", con.GetByID)
}

func (con *Controller) Start(ctx *gin.Context) (interface{}, error) {
	userId := ctx.MustGet(pkg.AccountIDKey).(int64)
	var p ParamStart
	projectId := ctx.Param("projectId")

	err := ctx.ShouldBindJSON(&p)
	if err != nil {
		return nil, err
	}

	projectIdInt64, err := strconv.ParseInt(projectId, 10, 64)
	if err != nil {
		return nil, err
	}

	exam, err := con.Service.Start(projectIdInt64, userId, p.Type)
	if err != nil {
		return nil, err
	}

	return NewStartExamVo(*exam), nil
}

func (con *Controller) Save(ctx *gin.Context) (interface{}, error) {
	userId := ctx.MustGet(pkg.AccountIDKey).(int64)
	examId := ctx.Param("examId")
	var p struct {
		Questions []ParamAnswer `json:"questions"`
	}
	err := ctx.ShouldBindJSON(&p)
	if err != nil {
		return nil, err
	}

	examIdInt, err := strconv.ParseInt(examId, 10, 64)
	if err != nil {
		return nil, err
	}

	var paper = map[int64][]int{}
	for _, q := range p.Questions {
		paper[q.ID] = q.Answers
	}

	err = con.Service.Save(examIdInt, userId, paper)
	if err != nil {
		return nil, err
	}

	return nil, nil
}

func (con *Controller) Submit(ctx *gin.Context) (interface{}, error) {
	userId := ctx.MustGet(pkg.AccountIDKey).(int64)
	examId := ctx.Param("examId")
	var p struct {
		Questions []ParamAnswer `json:"questions"`
	}
	err := ctx.ShouldBindJSON(&p)
	if err != nil {
		return nil, err
	}

	examIdInt, err := strconv.ParseInt(examId, 10, 64)
	if err != nil {
		return nil, err
	}

	var paper = map[int64][]int{}
	for _, q := range p.Questions {
		paper[q.ID] = q.Answers
	}

	certImg, certNum, pointsRst, err := con.Service.Submit(examIdInt, userId, paper)
	if err != nil {
		return nil, err
	}

	exam, exists, err := con.Service.GetByID(examIdInt)
	if err != nil || !exists {
		return nil, err
	}

	back := feedback.NewValuableFeedBack()
	err = back.AddPointsEventRsts(pointsRst)
	if err != nil || !exists {
		return nil, err
	}

	back.Put("score", exam.Score)
	back.Put("passed", exam.CheckPass())
	back.Put("certId", certNum)
	back.Put("certImg", certImg)

	return back, nil
}

func (con *Controller) ListSubmitted(ctx *gin.Context) (interface{}, error) {
	userId := ctx.MustGet(pkg.AccountIDKey).(int64)
	query := struct {
		Latest int `form:"latest"`
		Type   int `form:"type" binding:"required,min=1,max=2"`
	}{}
	projectIdStr := ctx.Param("projectId")

	projectId, err := strconv.ParseInt(projectIdStr, 10, 64)
	if err != nil {
		return nil, err
	}

	err = ctx.ShouldBindQuery(&query)
	if err != nil {
		return nil, err
	}

	exams, err := con.Service.ListLatestSubmitted(projectId, userId, query.Type, query.Latest)
	if err != nil {
		return nil, err
	}

	vos := make([]SimpleListInfoDto, len(exams))
	for i, e := range exams {
		vos[i] = ParseSimpleListInfoDto(e)
	}

	return map[string]interface{}{
		"exams": vos,
	}, nil
}

func (con *Controller) GetUnSubmittedLatest(ctx *gin.Context) (interface{}, error) {
	userId := ctx.MustGet(pkg.AccountIDKey).(int64)
	var query struct {
		Type int `form:"type" binding:"required,min=1,max=2"`
	}

	projectIdStr := ctx.Param("projectId")
	projectId, err := strconv.ParseInt(projectIdStr, 10, 64)
	if err != nil {
		return nil, err
	}
	err = ctx.ShouldBindQuery(&query)
	if err != nil {
		return nil, err
	}

	exam, exists, err := con.Service.GetLatestUnSubmitted(projectId, userId, query.Type)
	if err != nil {
		return nil, err
	}

	if !exists {
		return nil, nil
	}

	return ParseSimpleListInfoWithQuestionDto(exam), nil
}

func (con *Controller) GetByID(ctx *gin.Context) (interface{}, error) {
	idStr := ctx.Param("examId")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		return nil, err
	}

	exam, exists, err := con.Service.GetByID(id)
	if err != nil {
		return nil, err
	}

	if !exists {
		return nil, nil
	}

	return ParseDetailDTo(exam), nil
}
