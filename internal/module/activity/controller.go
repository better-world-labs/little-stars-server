package activity

import (
	"aed-api-server/internal/interfaces/entities"
	"aed-api-server/internal/interfaces/events"
	"aed-api-server/internal/pkg"
	"aed-api-server/internal/pkg/feedback"
	"aed-api-server/internal/pkg/global"
	"aed-api-server/internal/pkg/response"
	"aed-api-server/internal/pkg/utils"
	"aed-api-server/internal/service/activity"
	"github.com/gin-gonic/gin"
	"strconv"
	"strings"
	"time"
)

type Controller struct {
}

func NewController() Controller {
	return Controller{}
}

func (con *Controller) ListActivities(c *gin.Context) {
	aid, exists := c.GetQuery("aidId")
	limit := c.DefaultQuery("limit", "0")
	utils.MustTrue(exists, global.ErrorInvalidParam)

	i, err := strconv.ParseInt(aid, 10, 64)
	utils.MustNil(err, global.ErrorInvalidParam)

	l, err := strconv.Atoi(limit)
	utils.MustNil(err, global.ErrorInvalidParam)

	res, err := GetService().ListLatestCategorySorted(i, l)
	utils.MustNil(err, err)
	last, err := GetService().GetLastUpdated(i)
	utils.MustNil(err, err)

	resolveActivityImages(res)
	data := map[string]interface{}{"activities": res}
	if last != nil {
		data["lastUpdated"] = last.Created
	}

	response.ReplyOK(c, data)
}

func (con *Controller) CreateScene(c *gin.Context) {
	accountID := c.MustGet(pkg.AccountIDKey).(int64)
	var dto RecordSceneReportDTO
	err := c.BindJSON(&dto)
	utils.MustNil(err, global.ErrorInvalidParam)

	eventRst, err := GetService().SaveActivitySceneReport(events.NewSceneReportEvent(dto.AidID, accountID, dto.Description, dto.Images))
	utils.MustNil(err, err)

	back := feedback.NewValuableFeedBack()
	err = back.AddPointsEventRsts(eventRst)
	utils.MustNil(err, err)

	response.ReplyOK(c, back)
}

func (con *Controller) GetLatestActivity(c *gin.Context) {
	aid, exists := c.GetQuery("aidId")
	utils.MustTrue(exists, global.ErrorInvalidParam)

	i, err := strconv.ParseInt(aid, 10, 64)
	utils.MustNil(err, global.ErrorInvalidParam)

	res, err := GetService().ListLatestCategorySorted(i, 1)
	resolveActivityImages(res)
	if len(res) == 0 {
		response.ReplyOK(c, nil)
	} else {
		response.ReplyOK(c, res[0])
	}
}

func (con *Controller) GetOneByID(c *gin.Context) {
	id, exists := c.GetQuery("id")
	utils.MustTrue(exists, global.ErrorInvalidParam)

	i, err := strconv.ParseInt(id, 10, 64)
	utils.MustNil(err, global.ErrorInvalidParam)

	res, err := GetService().GetOneByID(i)
	utils.MustNil(err, err)

	response.ReplyOK(c, res)
}

func (con *Controller) GoingToDevice(c *gin.Context) {
	accountID := c.MustGet(pkg.AccountIDKey).(int64)
	req := new(GoingToDevice)
	err := c.BindJSON(req)
	if err != nil {
		response.ReplyError(c, err)
		return
	}

	err = GetService().SaveActivityGoingToGetDevice(events.NewGoingToGetDeviceEvent(req.AidId, accountID))
	if err != nil {
		response.ReplyError(c, err)
		return
	}

	//producer.Publish(record.CreateRecordGoingToGetDevice(accountID, int64(utils.ToInt(req.AidId))))
	response.ReplyOK(c, nil)
}

func (con *Controller) GetDevice(c *gin.Context) {
	accountID := c.MustGet(pkg.AccountIDKey).(int64)
	req := new(BorrowDevice)
	err := c.BindJSON(req)
	if err != nil {
		response.ReplyError(c, err)
		return
	}
	pointEvt, err := GetService().SaveActivityDeviceGot(events.NewDeviceGotEvent(req.AidId, accountID))
	if err != nil {
		response.ReplyError(c, err)
		return
	}

	back := feedback.NewValuableFeedBack()
	err = back.AddPointsEventRsts(pointEvt)
	if err != nil {
		response.ReplyError(c, err)
		return
	}

	response.ReplyOK(c, back)
}

func (con *Controller) GetManyByIDs(c *gin.Context) {
	ids, exists := c.GetQuery("ids")
	utils.MustTrue(exists, global.ErrorInvalidParam)

	var param []int64
	for _, e := range strings.Split(ids, ",") {
		i, err := strconv.ParseInt(e, 10, 64)
		utils.MustNil(err, global.ErrorInvalidParam)
		param = append(param, i)
	}

	res, err := GetService().GetManyByIDs(param)
	utils.MustNil(err, err)

	resolveActivityImages(res)
	response.ReplyOK(c, map[string]interface{}{"activities": res})
}

// 废弃
func (con *Controller) GoingToLearning(c *gin.Context) {
	accountID := c.MustGet(pkg.AccountIDKey).(int64)
	param := make(map[string]interface{}, 1)
	err := c.BindJSON(&param)
	utils.MustNil(err, global.ErrorInvalidParam)

	aid, err := strconv.ParseInt(param["aidId"].(string), 10, 64)
	utils.MustNil(err, err)
	err = GetService().Create(&entities.Activity{
		HelpInfoID: aid,
		Class:      activity.ClassSkillLearning,
		UserID:     &accountID,
		Created:    global.FormattedTime(time.Now()),
	})

	utils.MustNil(err, err)
	response.ReplyOK(c, nil)
}
