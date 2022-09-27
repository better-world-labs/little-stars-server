package controller

import (
	"aed-api-server/internal/interfaces/entities"
	"aed-api-server/internal/interfaces/service"
	"aed-api-server/internal/pkg"
	"aed-api-server/internal/pkg/asserts"
	"aed-api-server/internal/pkg/feedback"
	"aed-api-server/internal/pkg/location"
	"aed-api-server/internal/pkg/query"
	"aed-api-server/internal/pkg/response"
	"aed-api-server/internal/service/aid"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"gitlab.openviewtech.com/openview-pub/gopkg/route"
	"html/template"
	"strconv"
)

type AidController struct {
	Service   service.AidService   `inject:"-"`
	PathToken service.TokenService `inject:"-"`
}

//go:inject-component
func NewAidController() *AidController {
	return &AidController{}
}

func (con *AidController) MountNoAuthRouter(r *route.Router) {
	r.GET("/aid/infos", con.ListHelpInfosPaged)
	r.GET("/aid/infos-hours", con.ListOneHoursInfos)
	r.GET("/aid/info", con.GetHelpInfo)
	r.POST("/aid/aid-called/:token/:aid", con.ActionAidCalled)
	r.GET("/p/c/:token/:aid", con.ActionAidCalledPage)
}

func (con *AidController) MountAuthRouter(r *route.Router) {
	r.POST("/aid/publish", con.PublishHelpInfo)
	r.POST("/aid/exercise/publish", con.PublishHelpInfoExercise)
	r.POST("/aid/arrived", con.ActionArrived)
	r.POST("/aid/exercise/arrived", con.ActionNPCArrived)
	r.POST("/aid/going-to-scene", con.ActionGoingToScene)
	r.POST("/aid/called", con.ActionCalled)
	r.GET("/aid/me/published", con.ListMyHelpInfosPaged)
	r.GET("/aid/me/participated", con.ListHelpInfosParticipatedPaged)
	r.GET("/aid/me/all", con.ListHelpInfosMyAll)
	r.GET("/aid/track", con.GeAidTrackForUser)
}

func (con *AidController) PublishHelpInfo(c *gin.Context) (interface{}, error) {
	accountID := c.MustGet(pkg.AccountIDKey).(int64)
	var body entities.PublishDTO
	err := c.ShouldBindJSON(&body)
	if err != nil {
		return nil, err
	}

	id, eventRst, err := con.Service.PublishHelpInfo(accountID, &body)
	if err != nil {
		return nil, err
	}

	feedBack := feedback.NewValuableFeedBack()
	err = feedBack.AddPointsEventRsts(eventRst)
	if err != nil {
		return nil, err
	}

	feedBack.Put("id", strconv.FormatInt(id, 10))
	return feedBack, nil
}

func (con *AidController) PublishHelpInfoExercise(c *gin.Context) (interface{}, error) {
	accountID := c.MustGet(pkg.AccountIDKey).(int64)
	var body entities.PublishDTO
	err := c.ShouldBindJSON(&body)
	if err != nil {
		return nil, err
	}

	id, npc, err := con.Service.PublishHelpInfoExercise(accountID, &body)
	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"id":  id,
		"npc": npc,
	}, nil
}

func (con *AidController) ListHelpInfosPaged(c *gin.Context) (interface{}, error) {
	param := struct {
		*location.Coordinate

		Exercise bool `form:"exercise"`
	}{}

	err := c.ShouldBindQuery(&param)
	if err != nil {
		return nil, err
	}

	pageQuery, err := page.BindPageQuery(c)
	if err != nil {
		return nil, err
	}

	result, err := con.Service.ListHelpInfosPaged(pageQuery, param.Coordinate, &entities.HelpInfo{Exercise: param.Exercise})
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (con *AidController) ListOneHoursInfos(c *gin.Context) (interface{}, error) {
	result, err := con.Service.ListOneHoursInfos()
	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"infos": result,
	}, nil
}

func (con *AidController) ListHelpInfosParticipatedPaged(c *gin.Context) (interface{}, error) {
	accountID := c.MustGet(pkg.AccountIDKey).(int64)
	pageQuery, err := page.BindPageQuery(c)
	if err != nil {
		return nil, err
	}

	result, err := con.Service.ListHelpInfosParticipatedPaged(pageQuery, accountID)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (con *AidController) ListMyHelpInfosPaged(c *gin.Context) (interface{}, error) {
	accountID := c.MustGet(pkg.AccountIDKey).(int64)
	pageQuery, err := page.BindPageQuery(c)
	if err != nil {
		return nil, err
	}

	result, err := con.Service.ListHelpInfosPaged(pageQuery, nil, &entities.HelpInfo{Publisher: accountID})
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (con *AidController) ActionArrived(c *gin.Context) (interface{}, error) {
	accountID := c.MustGet(pkg.AccountIDKey).(int64)
	var param entities.ActionDTO
	err := c.ShouldBindJSON(&param)
	if err != nil {
		return nil, response.NewIllegalArgumentError(err.Error())
	}

	eventRst, err := con.Service.ActionArrived(accountID, param.AidID, param.Coordinate)
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

func (con *AidController) ActionCalled(c *gin.Context) (interface{}, error) {
	accountID := c.MustGet(pkg.AccountIDKey).(int64)
	var param entities.ActionDTO
	err := c.ShouldBindJSON(&param)
	if err != nil {
		return nil, response.NewIllegalArgumentError(err.Error())
	}

	err = con.Service.ActionCalled(accountID, param.AidID)
	if err != nil {
		return nil, err
	}

	return nil, nil
}

func (con *AidController) ActionGoingToScene(c *gin.Context) (interface{}, error) {
	accountID := c.MustGet(pkg.AccountIDKey).(int64)
	var param entities.ActionDTO
	err := c.ShouldBindJSON(&param)
	if err != nil {
		return nil, response.NewIllegalArgumentError(err.Error())
	}

	err = con.Service.ActionGoingToScene(accountID, param.AidID)
	if err != nil {
		return nil, err
	}

	return nil, nil
}

func (con *AidController) GeAidTrackForUser(c *gin.Context) (interface{}, error) {
	accountID := c.MustGet(pkg.AccountIDKey).(int64)
	query := struct {
		AidId int64 `form:"aidId"`
	}{}
	err := c.ShouldBindQuery(&query)
	if err != nil {
		return nil, response.NewIllegalArgumentError(err.Error())
	}

	t, err := aid.GetService().GetUserTrack(query.AidId, accountID)
	if err != nil {
		return nil, err
	}

	type DTO struct {
		DeviceGot    bool `json:"deviceGot"`
		SceneArrived bool `json:"sceneArrived"`
	}
	return DTO{
		DeviceGot:    t.DeviceGot,
		SceneArrived: t.SceneArrived,
	}, nil
}

func (con *AidController) GetHelpInfo(c *gin.Context) (interface{}, error) {
	var position *location.Coordinate
	id, existsId := c.GetQuery("id")
	lat, existsLat := c.GetQuery("latitude")
	lon, existsLon := c.GetQuery("longitude")

	if !existsId {
		return nil, response.NewIllegalArgumentError("invalid param")
	}

	if existsLat && existsLon {
		latF, err := strconv.ParseFloat(lat, 64)
		if err != nil {
			return nil, response.NewIllegalArgumentError(err.Error())
		}

		lonF, err := strconv.ParseFloat(lon, 64)
		if err != nil {
			return nil, response.NewIllegalArgumentError(err.Error())
		}
		position = &location.Coordinate{
			Longitude: lonF,
			Latitude:  latF,
		}
	}
	i, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		return nil, response.NewIllegalArgumentError(err.Error())
	}

	composed, exists, err := con.Service.GetHelpInfoComposedByID(i, position)
	if err != nil {
		return nil, response.NewIllegalArgumentError(err.Error())
	}

	if exists {
		return composed, nil
	} else {
		return nil, nil
	}
}

func (con *AidController) ActionAidCalled(c *gin.Context) (interface{}, error) {
	token := c.Param("token")
	aid := c.Param("aid")
	var p entities.Call120RequestDto
	err := c.ShouldBindJSON(&p)
	if err != nil {
		return nil, response.NewIllegalArgumentError(err.Error())
	}

	aidId, err := strconv.ParseInt(aid, 10, 64)
	if err != nil {
		return nil, response.NewIllegalArgumentError(err.Error())
	}

	ok, err := con.PathToken.ValidateToken(token, p.MobileLast4)
	if err != nil {
		return nil, err
	}

	if !ok {
		return nil, err
	}

	err = con.Service.Action120Called(aidId)
	if err != nil {
		return nil, err
	}

	_, err = con.PathToken.RemoveToken(token)
	if err != nil {
		return nil, err
	}

	return nil, nil
}

func (con *AidController) ActionAidCalledPage(c *gin.Context) (interface{}, error) {
	token := c.Param("token")
	aid := c.Param("aid")
	tpl := template.New("ActionAidCalled")
	html, _ := asserts.GetResource("called_trigger.html")

	tpl, err := tpl.Parse(string(html))
	if err != nil {
		return nil, err
	}

	c.Status(200)
	c.Header("Content-Type", "text/html;charset=utf-8")
	err = tpl.Execute(c.Writer, map[string]interface{}{"token": token, "aid": aid})
	if err != nil {
		log.Errorf("tpl execute error: %v", err)
	}

	return nil, nil
}

func (con *AidController) ListHelpInfosMyAll(c *gin.Context) (interface{}, error) {
	accountID := c.MustGet(pkg.AccountIDKey).(int64)
	pageQuery, err := page.BindPageQuery(c)
	if err != nil {
		return nil, err
	}

	participated, err := con.Service.ListHelpInfosParticipatedPaged(pageQuery, accountID)
	if err != nil {
		return nil, err
	}

	published, err := con.Service.ListHelpInfosPaged(pageQuery, nil, &entities.HelpInfo{Publisher: accountID})
	if err != nil {
		return nil, err
	}

	return page.NewResult[*entities.HelpInfoComposedDTO](append(participated.List, published.List...), participated.Total+published.Total), nil
}

func (con *AidController) ActionNPCArrived(ctx *gin.Context) (interface{}, error) {
	var param entities.ActionDTO

	err := ctx.ShouldBindJSON(&param)
	if err != nil {
		return nil, err
	}

	return nil, con.Service.ActionNPCArrived(param.AidID)
}
