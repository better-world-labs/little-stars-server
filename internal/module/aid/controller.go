package aid

import (
	"aed-api-server/internal/interfaces/entities"
	"aed-api-server/internal/module/aid/track"
	"aed-api-server/internal/module/speech"
	"aed-api-server/internal/pkg"
	"aed-api-server/internal/pkg/asserts"
	"aed-api-server/internal/pkg/feedback"
	"aed-api-server/internal/pkg/global"
	"aed-api-server/internal/pkg/location"
	"aed-api-server/internal/pkg/query"
	"aed-api-server/internal/pkg/response"
	"aed-api-server/internal/pkg/utils"
	"fmt"
	"github.com/gin-gonic/gin"
	"gitlab.openviewtech.com/openview-pub/gopkg/log"
	"html/template"
	"strconv"
)

type Controller struct {
	service   Service
	pathToken speech.TokenService
}

func NewController(service Service) Controller {
	return Controller{service: service, pathToken: speech.NewTokenService()}
}
func (con *Controller) PublishHelpInfo(c *gin.Context) {
	accountID := c.MustGet(pkg.AccountIDKey).(int64)
	var body PublishDTO
	err := c.BindJSON(&body)
	if err != nil {
		log.DefaultLogger().Errorf("bind json error: %v", err)
		response.ReplyError(c, err)
		return
	}

	id, eventRst, err := con.service.PublishHelpInfo(accountID, &body)
	if err != nil {
		log.DefaultLogger().Errorf("bind json error: %v", err)
		response.ReplyError(c, err)
	}

	feedBack := feedback.NewValuableFeedBack()
	err = feedBack.AddPointsEventRsts(eventRst)
	if err != nil {
		log.DefaultLogger().Errorf("bind json error: %v", err)
		response.ReplyError(c, err)
	}

	feedBack.Put("id", strconv.FormatInt(id, 10))
	response.ReplyOK(c, feedBack)
}

func (con *Controller) ListHelpInfosPaged(c *gin.Context) {
	var position *location.Coordinate
	lat := c.Query("latitude")
	lon := c.Query("longitude")

	if lat != "" && lon != "" {
		log.DefaultLogger().Infof("lon=%s, lat=%s", lon, lat)
		latF, err := strconv.ParseFloat(lat, 64)
		utils.MustNil(err, global.ErrorInvalidParam)
		lonF, err := strconv.ParseFloat(lon, 64)
		utils.MustNil(err, global.ErrorInvalidParam)
		position = &location.Coordinate{Longitude: lonF, Latitude: latF}
	}

	pageQuery, err := page.BindPageQuery(c)
	if err != nil {
		log.DefaultLogger().Errorf("bind query error: %v", err)
		response.ReplyError(c, err)
		return
	}

	result, err := con.service.ListHelpInfosPaged(pageQuery, position, &entities.HelpInfo{})
	if err != nil {
		log.DefaultLogger().Errorf("bind query error: %v", err)
		response.ReplyError(c, err)
		return
	}

	response.ReplyOK(c, result)
}

func (con *Controller) ListOneHoursInfos(c *gin.Context) {
	result, err := con.service.ListOneHoursInfos()
	if err != nil {
		log.DefaultLogger().Errorf("bind query error: %v", err)
		response.ReplyError(c, err)
		return
	}

	response.ReplyOK(c, map[string]interface{}{
		"infos": result,
	})
}

func (con *Controller) ListHelpInfosParticipatedPaged(c *gin.Context) {
	accountID := c.MustGet(pkg.AccountIDKey).(int64)
	pageQuery, err := page.BindPageQuery(c)
	if err != nil {
		log.DefaultLogger().Errorf("bind query error: %v", err)
		response.ReplyError(c, err)
		return
	}

	result, err := con.service.ListHelpInfosParticipatedPaged(pageQuery, accountID)
	response.ReplyOK(c, result)
}

func (con *Controller) ListMyHelpInfosPaged(c *gin.Context) {
	accountID := c.MustGet(pkg.AccountIDKey).(int64)
	pageQuery, err := page.BindPageQuery(c)
	if err != nil {
		log.DefaultLogger().Errorf("bind query error: %v", err)
		response.ReplyError(c, err)
		return
	}

	result, err := con.service.ListHelpInfosPaged(pageQuery, nil, &entities.HelpInfo{Publisher: accountID})
	response.ReplyOK(c, result)
}

func (con *Controller) ActionArrived(c *gin.Context) {
	accountID := c.MustGet(pkg.AccountIDKey).(int64)
	var param ActionDTO
	err := c.BindJSON(&param)
	utils.MustNil(err, global.ErrorInvalidParam)

	aidInteger, err := strconv.ParseInt(param.AidID, 10, 64)
	utils.MustNil(err, global.ErrorInvalidParam)

	eventRst, err := con.service.ActionArrived(accountID, aidInteger, param.Coordinate)
	utils.MustNil(err, err)

	back := feedback.NewValuableFeedBack()
	err = back.AddPointsEventRsts(eventRst)
	utils.MustNil(err, err)

	response.ReplyOK(c, back)
}

func (con *Controller) ActionCalled(c *gin.Context) {
	accountID := c.MustGet(pkg.AccountIDKey).(int64)
	var param ActionDTO
	err := c.BindJSON(&param)
	utils.MustNil(err, global.ErrorInvalidParam)

	aidInteger, err := strconv.ParseInt(param.AidID, 10, 64)
	utils.MustNil(err, global.ErrorInvalidParam)

	err = con.service.ActionCalled(accountID, aidInteger)
	utils.MustNil(err, global.ErrorInvalidParam)

	response.ReplyOK(c, nil)
}

func (con *Controller) ActionGoingToScene(c *gin.Context) {
	accountID := c.MustGet(pkg.AccountIDKey).(int64)
	var param ActionDTO
	err := c.BindJSON(&param)
	utils.MustNil(err, global.ErrorInvalidParam)

	aidInteger, err := strconv.ParseInt(param.AidID, 10, 64)
	utils.MustNil(err, global.ErrorInvalidParam)

	err = con.service.ActionGoingToScene(accountID, aidInteger)
	utils.MustNil(err, err)

	response.ReplyOK(c, nil)
}

func (con *Controller) GeAidTrackForUser(c *gin.Context) {
	accountID := c.MustGet(pkg.AccountIDKey).(int64)
	aidId, exists := c.GetQuery("aidId")
	utils.MustTrue(exists, global.ErrorInvalidParam)

	i, err := strconv.ParseInt(aidId, 10, 64)
	utils.MustNil(err, global.ErrorInvalidParam)

	t, err := track.GetService().GetUserTrack(i, accountID)
	response.ReplyOK(c, track.DTO{
		DeviceGot:    t.DeviceGot,
		SceneArrived: t.SceneArrived,
	})
}

func (con *Controller) GetHelpInfo(c *gin.Context) {
	var position *location.Coordinate
	id, existsId := c.GetQuery("id")
	lat, existsLat := c.GetQuery("latitude")
	lon, existsLon := c.GetQuery("longitude")

	if !existsId {
		response.ReplyError(c, global.ErrorInvalidParam)
		return
	}

	if existsLat && existsLon {
		latF, err := strconv.ParseFloat(lat, 64)
		utils.MustNil(err, global.ErrorInvalidParam)
		lonF, err := strconv.ParseFloat(lon, 64)
		utils.MustNil(err, global.ErrorInvalidParam)
		position = &location.Coordinate{
			Longitude: lonF,
			Latitude:  latF,
		}
	}
	i, err := strconv.ParseInt(id, 10, 64)
	utils.MustNil(err, global.ErrorInvalidParam)

	composed, exists, err := con.service.GetHelpInfoComposedByID(i, position)
	utils.MustNil(err, err)

	if exists {
		response.ReplyOK(c, composed)
	} else {
		response.ReplyOK(c, nil)
	}
}

func (con *Controller) ActionAidCalled(c *gin.Context) {
	token := c.Param("token")
	aid := c.Param("aid")
	var p Call120RequestDto
	err := c.BindJSON(&p)
	utils.MustNil(err, err)

	aidId, err := strconv.ParseInt(aid, 10, 64)
	utils.MustNil(err, err)

	ok, err := con.pathToken.ValidateToken(token, p.MobileLast4)
	if err != nil {
		log.DefaultLogger().Errorf("ValidateToken error: %v", err)
		response.ReplyError(c, global.ErrorInvalidParam)
		return
	}

	if !ok {
		log.DefaultLogger().Errorf("ValidateToken %t", ok)
		response.ReplyError(c, global.ErrorLinkInvalid)
		return
	}

	err = con.service.Action120Called(aidId)
	if err != nil {
		log.DefaultLogger().Errorf("Action120Called error: %v", err)
		response.ReplyError(c, global.ErrorInvalidParam)
		return
	}

	_, err = con.pathToken.RemoveToken(token)
	if err != nil {
		log.DefaultLogger().Errorf("remove token %s error", token)
		return
	}

	response.ReplyOK(c, nil)
}

func (con *Controller) ActionAidCalledPage(c *gin.Context) {
	token := c.Param("token")
	aid := c.Param("aid")
	tpl := template.New("ActionAidCalled")
	html, _ := asserts.GetResource("called_trigger.html")

	tpl, err := tpl.Parse(string(html))
	utils.MustNil(err, err)

	c.Status(200)
	c.Header("Content-Type", "text/html;charset=utf-8")
	err = tpl.Execute(c.Writer, map[string]interface{}{"token": token, "aid": aid})
	if err != nil {
		log.DefaultLogger().Errorf("tpl execute error: %v", err)
	}
}

func (con *Controller) ListHelpInfosMyAll(c *gin.Context) {
	accountID := c.MustGet(pkg.AccountIDKey).(int64)
	pageQuery, err := page.BindPageQuery(c)
	if err != nil {
		log.DefaultLogger().Errorf("bind query error: %v", err)
		response.ReplyError(c, err)
		return
	}

	participated, err := con.service.ListHelpInfosParticipatedPaged(pageQuery, accountID)
	if err != nil {
		response.ReplyError(c, err)
		return
	}

	published, err := con.service.ListHelpInfosPaged(pageQuery, nil, &entities.HelpInfo{Publisher: accountID})
	if err != nil {
		response.ReplyError(c, err)
		return
	}

	response.ReplyOK(c, page.Result{List: append(participated.List.([]*HelpInfoComposedDTO), published.List.([]*HelpInfoComposedDTO)...), Total: participated.Total + published.Total})
}

func ResponseHTML(c *gin.Context, status int, message string) {
	c.Writer.Header().Set("Content-Type", "text/html;charset=utf-8")
	c.Status(status)
	_, err := c.Writer.WriteString(fmt.Sprintf(`
		<html>
          <h1>%s</h1>
        </html>
    `, message))
	if err != nil {
		log.DefaultLogger().Errorf("write html error: %s", err)
	}
}
