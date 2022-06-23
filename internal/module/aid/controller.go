package aid

import (
	"aed-api-server/internal/interfaces/entities"
	service2 "aed-api-server/internal/interfaces/service"
	"aed-api-server/internal/module/aid/track"
	"aed-api-server/internal/module/speech"
	"aed-api-server/internal/pkg"
	"aed-api-server/internal/pkg/asserts"
	"aed-api-server/internal/pkg/feedback"
	"aed-api-server/internal/pkg/location"
	"aed-api-server/internal/pkg/query"
	"aed-api-server/internal/pkg/response"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"gitlab.openviewtech.com/openview-pub/gopkg/route"
	"html/template"
	"strconv"
)

type Controller struct {
	Service   service2.AidService `inject:"-"`
	PathToken speech.TokenService
}

func NewController() *Controller {
	return &Controller{PathToken: speech.NewTokenService()}
}

func (con *Controller) MountNoAuthRouter(r *route.Router) {
	r.GET("/aid/infos", con.ListHelpInfosPaged)
	r.GET("/aid/infos-hours", con.ListOneHoursInfos)
	r.GET("/aid/info", con.GetHelpInfo)
	r.POST("/aid/aid-called/:token/:aid", con.ActionAidCalled)
	r.GET("/c/:token/:aid", con.ActionAidCalledPage)
}

func (con *Controller) MountAuthRouter(r *route.Router) {
	r.POST("/aid/publish", con.PublishHelpInfo)
	r.POST("/aid/arrived", con.ActionArrived)
	r.POST("/aid/going-to-scene", con.ActionGoingToScene)
	r.POST("/aid/called", con.ActionCalled)
	r.GET("/aid/me/published", con.ListMyHelpInfosPaged)
	r.GET("/aid/me/participated", con.ListHelpInfosParticipatedPaged)
	r.GET("/aid/me/all", con.ListHelpInfosMyAll)
	r.GET("/aid/track", con.GeAidTrackForUser)
}

func (con *Controller) PublishHelpInfo(c *gin.Context) (interface{}, error) {
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

func (con *Controller) ListHelpInfosPaged(c *gin.Context) (interface{}, error) {
	var position *location.Coordinate
	lat := c.Query("latitude")
	lon := c.Query("longitude")

	if lat != "" && lon != "" {
		log.Infof("lon=%s, lat=%s", lon, lat)
		latF, err := strconv.ParseFloat(lat, 64)
		if err != nil {
			return nil, response.NewIllegalArgumentError(err.Error())
		}

		lonF, err := strconv.ParseFloat(lon, 64)
		if err != nil {
			return nil, response.NewIllegalArgumentError(err.Error())
		}

		position = &location.Coordinate{Longitude: lonF, Latitude: latF}
	}

	pageQuery, err := page.BindPageQuery(c)
	if err != nil {
		return nil, err
	}

	result, err := con.Service.ListHelpInfosPaged(pageQuery, position, &entities.HelpInfo{})
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (con *Controller) ListOneHoursInfos(c *gin.Context) (interface{}, error) {
	result, err := con.Service.ListOneHoursInfos()
	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"infos": result,
	}, nil
}

func (con *Controller) ListHelpInfosParticipatedPaged(c *gin.Context) (interface{}, error) {
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

func (con *Controller) ListMyHelpInfosPaged(c *gin.Context) (interface{}, error) {
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

func (con *Controller) ActionArrived(c *gin.Context) (interface{}, error) {
	accountID := c.MustGet(pkg.AccountIDKey).(int64)
	var param entities.ActionDTO
	err := c.ShouldBindJSON(&param)
	if err != nil {
		return nil, response.NewIllegalArgumentError(err.Error())
	}

	aidInteger, err := strconv.ParseInt(param.AidID, 10, 64)
	if err != nil {
		return nil, response.NewIllegalArgumentError(err.Error())
	}

	eventRst, err := con.Service.ActionArrived(accountID, aidInteger, param.Coordinate)
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

func (con *Controller) ActionCalled(c *gin.Context) (interface{}, error) {
	accountID := c.MustGet(pkg.AccountIDKey).(int64)
	var param entities.ActionDTO
	err := c.ShouldBindJSON(&param)
	if err != nil {
		return nil, response.NewIllegalArgumentError(err.Error())
	}

	aidInteger, err := strconv.ParseInt(param.AidID, 10, 64)
	if err != nil {
		return nil, response.NewIllegalArgumentError(err.Error())
	}

	err = con.Service.ActionCalled(accountID, aidInteger)
	if err != nil {
		return nil, err
	}

	return nil, nil
}

func (con *Controller) ActionGoingToScene(c *gin.Context) (interface{}, error) {
	accountID := c.MustGet(pkg.AccountIDKey).(int64)
	var param entities.ActionDTO
	err := c.ShouldBindJSON(&param)
	if err != nil {
		return nil, response.NewIllegalArgumentError(err.Error())
	}

	aidInteger, err := strconv.ParseInt(param.AidID, 10, 64)
	if err != nil {
		return nil, response.NewIllegalArgumentError(err.Error())
	}

	err = con.Service.ActionGoingToScene(accountID, aidInteger)
	if err != nil {
		return nil, err
	}

	return nil, nil
}

func (con *Controller) GeAidTrackForUser(c *gin.Context) (interface{}, error) {
	accountID := c.MustGet(pkg.AccountIDKey).(int64)
	query := struct {
		AidId int64 `form:"aidId"`
	}{}
	err := c.ShouldBindQuery(&query)
	if err != nil {
		return nil, response.NewIllegalArgumentError(err.Error())
	}

	t, err := track.GetService().GetUserTrack(query.AidId, accountID)
	if err != nil {
		return nil, err
	}

	return track.DTO{
		DeviceGot:    t.DeviceGot,
		SceneArrived: t.SceneArrived,
	}, nil
}

func (con *Controller) GetHelpInfo(c *gin.Context) (interface{}, error) {
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

func (con *Controller) ActionAidCalled(c *gin.Context) (interface{}, error) {
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

func (con *Controller) ActionAidCalledPage(c *gin.Context) (interface{}, error) {
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

func (con *Controller) ListHelpInfosMyAll(c *gin.Context) (interface{}, error) {
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

	return page.Result{List: append(participated.List.([]*entities.HelpInfoComposedDTO), published.List.([]*entities.HelpInfoComposedDTO)...), Total: participated.Total + published.Total}, nil
}
