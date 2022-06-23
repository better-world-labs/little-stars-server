package controller

import (
	"aed-api-server/internal/interfaces/entities"
	"aed-api-server/internal/interfaces/events"
	activity2 "aed-api-server/internal/module/activity"
	"aed-api-server/internal/pkg"
	"aed-api-server/internal/pkg/feedback"
	"aed-api-server/internal/pkg/global"
	"aed-api-server/internal/pkg/response"
	"aed-api-server/internal/service/activity"
	"github.com/gin-gonic/gin"
	"gitlab.openviewtech.com/openview-pub/gopkg/route"
	"strconv"
	"strings"
	"time"
)

type ActivityController struct {
}

func NewActivityController() *ActivityController {
	return &ActivityController{}
}

func (con *ActivityController) MountNoAuthRouter(r *route.Router) {
	r.GET("/aid/activities-sorted", con.ListActivities)
	r.GET("/aid/activities-sorted/latest", con.GetLatestActivity)
	r.GET("/aid/activity", con.GetOneByID)
	r.GET("/aid/activities", con.GetManyByIDs)
}

func (con *ActivityController) MountAuthRouter(r *route.Router) {
	r.POST("/aid/going-to-learning", con.GoingToLearning)
	r.POST("/aid/scene-report", con.CreateScene)
	r.POST("/aid/going-to-device", con.GoingToDevice)
	r.POST("/aed/borrow", con.GetDevice)
}

func (con *ActivityController) ListActivities(c *gin.Context) (interface{}, error) {
	query := struct {
		Limit int   `form:"limit"`
		AidId int64 `form:"aidId"`
	}{}

	err := c.ShouldBindQuery(&query)
	if err != nil {
		return nil, response.NewIllegalArgumentError(err.Error())
	}
	res, err := activity2.GetService().ListLatestCategorySorted(query.AidId, query.Limit)
	if err != nil {
		return nil, err
	}

	last, err := activity2.GetService().GetLastUpdated(query.AidId)
	if err != nil {
		return nil, err
	}

	resolveActivityImages(res)
	data := map[string]interface{}{"activities": res}
	if last != nil {
		data["lastUpdated"] = last.Created
	}

	return data, nil
}

func (con *ActivityController) CreateScene(c *gin.Context) (interface{}, error) {
	accountID := c.MustGet(pkg.AccountIDKey).(int64)
	var dto RecordSceneReportDTO
	err := c.ShouldBindJSON(&dto)
	if err != nil {
		return nil, response.NewIllegalArgumentError(err.Error())
	}

	eventRst, err := activity2.GetService().SaveActivitySceneReport(events.NewSceneReportEvent(dto.AidID, accountID, dto.Description, dto.Images))
	if err != nil {
		return nil, err
	}

	back := feedback.NewValuableFeedBack()
	err = back.AddPointsEventRsts(eventRst)
	if err != nil {
		return nil, err
	}

	return back, nil
}

func (con *ActivityController) GetLatestActivity(c *gin.Context) (interface{}, error) {
	query := struct {
		AidId int64 `form:"aidId"`
	}{}

	err := c.ShouldBindQuery(&query)
	if err != nil {
		return nil, response.NewIllegalArgumentError(err.Error())
	}

	res, err := activity2.GetService().ListLatestCategorySorted(query.AidId, 1)
	resolveActivityImages(res)
	if len(res) == 0 {
		return nil, nil
	} else {
		return res[0], nil
	}
}

func (con *ActivityController) GetOneByID(c *gin.Context) (interface{}, error) {
	query := struct {
		Id int64 `form:"id"`
	}{}

	err := c.ShouldBindQuery(&query)
	if err != nil {
		return nil, response.NewIllegalArgumentError(err.Error())
	}

	res, err := activity2.GetService().GetOneByID(query.Id)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (con *ActivityController) GoingToDevice(c *gin.Context) (interface{}, error) {
	accountID := c.MustGet(pkg.AccountIDKey).(int64)
	req := new(activity2.GoingToDevice)
	err := c.ShouldBindJSON(req)
	if err != nil {
		return nil, err
	}

	err = activity2.GetService().SaveActivityGoingToGetDevice(events.NewGoingToGetDeviceEvent(req.AidId, accountID))
	if err != nil {
		return nil, err
	}

	return nil, nil
}

func (con *ActivityController) GetDevice(c *gin.Context) (interface{}, error) {
	accountID := c.MustGet(pkg.AccountIDKey).(int64)
	req := new(activity2.BorrowDevice)
	err := c.ShouldBindJSON(req)
	if err != nil {
		return nil, err
	}

	pointEvt, err := activity2.GetService().SaveActivityDeviceGot(events.NewDeviceGotEvent(req.AidId, accountID))
	if err != nil {
		return nil, err
	}

	back := feedback.NewValuableFeedBack()
	err = back.AddPointsEventRsts(pointEvt)
	if err != nil {
		return nil, err
	}

	return back, nil
}

func (con *ActivityController) GetManyByIDs(c *gin.Context) (interface{}, error) {
	ids, exists := c.GetQuery("ids")
	if !exists {
		return nil, response.NewIllegalArgumentError("invalid param")
	}

	var param []int64
	for _, e := range strings.Split(ids, ",") {
		i, err := strconv.ParseInt(e, 10, 64)
		if err != nil {
			return nil, response.NewIllegalArgumentError("invalid param")
		}

		param = append(param, i)
	}

	res, err := activity2.GetService().GetManyByIDs(param)
	if err != nil {
		return nil, err
	}

	resolveActivityImages(res)
	return map[string]interface{}{"activities": res}, nil
}

// 废弃
func (con *ActivityController) GoingToLearning(c *gin.Context) (interface{}, error) {
	accountID := c.MustGet(pkg.AccountIDKey).(int64)
	param := make(map[string]interface{}, 1)
	err := c.ShouldBindJSON(&param)
	if err != nil {
		return nil, response.NewIllegalArgumentError(err.Error())
	}

	aid, err := strconv.ParseInt(param["aidId"].(string), 10, 64)
	if err != nil {
		return nil, response.NewIllegalArgumentError(err.Error())
	}

	err = activity2.GetService().Create(&entities.Activity{
		HelpInfoID: aid,
		Class:      activity.ClassSkillLearning,
		UserID:     &accountID,
		Created:    global.FormattedTime(time.Now()),
	})

	if err != nil {
		return nil, err
	}

	return nil, nil
}

type RecordSceneReportDTO struct {
	AidID       int64    `json:"aidId,string"`
	Description string   `json:"description"`
	Images      []string `json:"images"`
}

func resolveActivityImages(activities []*entities.Activity) {
	for _, a := range activities {
		resolveActivityImage(a)
	}
}

func resolveActivityImage(a *entities.Activity) {
	if a.Class == activity.ClassSceneReport {
		if img, ok := a.Record["images"]; ok {
			switch img.(type) {
			case []interface{}:
				imgStrings := parseImagesToStrings(img.([]interface{}))
				a.Record["images"] = imgStrings
			}
		}
	}
}

func parseImagesToStrings(r []interface{}) []string {
	var images []string
	for _, img := range r {
		switch img.(type) {
		case map[string]interface{}:
			if origin, ok := img.(map[string]interface{})["origin"]; ok {
				images = append(images, origin.(string))
			}
		default:
			images = append(images, img.(string))
		}
	}

	return images
}
