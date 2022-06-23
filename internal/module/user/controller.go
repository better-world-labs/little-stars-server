package user

import (
	"aed-api-server/internal/interfaces"
	"aed-api-server/internal/interfaces/entities"
	service2 "aed-api-server/internal/interfaces/service"
	"aed-api-server/internal/pkg"
	"aed-api-server/internal/pkg/config"
	"aed-api-server/internal/pkg/location"
	"aed-api-server/internal/pkg/response"
	"crypto/md5"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"gitlab.openviewtech.com/openview-pub/gopkg/route"
	"strconv"
	"strings"
	"time"
)

type Controller struct {
	Service     service2.UserServiceOld `inject:"-"`
	backendConf config.Backend
}

func NewController(backend config.Backend) *Controller {
	return &Controller{backendConf: backend}
}

func (con *Controller) MountAuthRouter(r *route.Router) {
	r.PUT("/user/position", con.UpdatePosition)
	r.PUT("/user/mobile", con.UpdateMobile)
	r.PUT("/user/info", con.UpdateUserInfo)
	r.GET("/user/info", con.GetUserInfo)
	r.GET("/user/infos", con.GetUserInfos)
	r.POST("/user/check-token", con.CheckUserToken)
	r.POST("/user/events", con.DealReportedEvents)
	r.GET("/user/charity-card", con.GetUserCharityCard)
	r.POST("/user/message-subscribe", con.ReportSubscribeMessageSettings)
	r.GET("/user/message-subscribe/last", con.GetLastReportSubscribeMessageSettings)
}

func (con *Controller) MountNoAuthRouter(r *route.Router) {
	r.POST("/user/wechat/app/login", con.WechatAppLogin)
	r.POST("/user/wechat/mini-program/login", con.WechatMiniProgramLogin)
	r.POST("/v2/user/wechat/mini-program/login", con.WechatMiniProgramLoginV2)
	r.POST("/user/wechat/mini-program/login-simple", con.WechatMiniProgramLoginSimple)
	r.POST("/user/generate-uid", con.GenerateUid)
	r.GET("/user-counts", con.CountRegisteredUsers)
}

func (con *Controller) MountGinEngineRouter(r *route.Router) {
	r.POST("/admin-api/user/login", con.Login)
}

func (con *Controller) WechatAppLogin(c *gin.Context) (interface{}, error) {
	return nil, errors.New("nod implements")
}

func (con *Controller) WechatMiniProgramLogin(c *gin.Context) (interface{}, error) {
	var body entities.LoginCommand
	err := c.ShouldBindJSON(&body)
	log.Infof("mobileCode: %s", body.MobileCode)

	if err != nil {
		return nil, err
	}

	res, token, _, err := con.Service.WechatMiniProgramLogin(body)
	if err != nil {
		return nil, err
	}

	log.Info("Login: uid=", res.Uid)
	return &entities.UserDTO{
		ID:       res.ID,
		Uid:      res.Uid,
		Nickname: res.Nickname,
		Token:    token,
		Avatar:   res.Avatar,
		Mobile:   res.Mobile,
		Openid:   res.Openid,
	}, nil
}

func (con *Controller) WechatMiniProgramLoginV2(c *gin.Context) (interface{}, error) {
	var body entities.LoginCommandV2
	err := c.ShouldBindJSON(&body)

	if err != nil {
		return nil, err
	}

	res, token, _, err := con.Service.WechatMiniProgramLoginV2(body)
	if err != nil {
		return nil, err
	}

	log.Info("Login: uid=", res.Uid)
	return &entities.UserDTO{
		ID:       res.ID,
		Uid:      res.Uid,
		Nickname: res.Nickname,
		Token:    token,
		Avatar:   res.Avatar,
		Mobile:   res.Mobile,
		Openid:   res.Openid,
	}, nil
}

func (con *Controller) WechatMiniProgramLoginSimple(c *gin.Context) (interface{}, error) {
	var body entities.SimpleLoginCommand
	err := c.ShouldBindJSON(&body)

	if err != nil {
		return nil, err
	}

	res, token, _, err := con.Service.WechatMiniProgramLoginSimple(body)
	if err != nil {
		return nil, err
	}

	log.Info("SimpleLogin: uid=", res.Uid)
	return &entities.UserDTO{
		ID:       res.ID,
		Uid:      res.Uid,
		Nickname: res.Nickname,
		Token:    token,
		Avatar:   res.Avatar,
		Mobile:   res.Mobile,
		Openid:   res.Openid,
	}, nil
}

func (con *Controller) UpdatePosition(c *gin.Context) (interface{}, error) {
	accountID := c.MustGet(pkg.AccountIDKey).(int64)
	var body entities.Position
	body.AccountID = accountID
	err := c.ShouldBindJSON(&body)
	if err != nil {
		return nil, err
	}

	err = con.Service.UpdatePosition(&body)
	if err != nil {
		return nil, err
	}

	//调用任务接口生成任务
	go interfaces.S.Task.GenJobsByUserLocation(accountID, location.Coordinate{
		Latitude:  body.Latitude,
		Longitude: body.Longitude,
	})

	return nil, nil
}

func (con *Controller) UpdateMobile(c *gin.Context) (interface{}, error) {
	accountID := c.MustGet(pkg.AccountIDKey).(int64)
	var body entities.MobileCommand
	err := c.ShouldBindJSON(&body)
	if err != nil {
		return nil, err
	}

	err = con.Service.UpdateMobile(accountID, body)
	if err != nil {
		return nil, err
	}

	account, err := con.Service.GetUserByID(accountID)
	if err != nil {
		return nil, err
	}

	return &entities.UserDTO{
		ID:       account.ID,
		Uid:      account.Uid,
		Nickname: account.Nickname,
		Avatar:   account.Avatar,
		Mobile:   account.Mobile,
		Openid:   account.Openid,
	}, nil
}

func (con *Controller) UpdateUserInfo(c *gin.Context) (interface{}, error) {
	accountID := c.MustGet(pkg.AccountIDKey).(int64)

	var body entities.UserDTO
	err := c.ShouldBindJSON(&body)
	if err != nil {
		return nil, err
	}

	err = con.Service.UpdateUserInfo(&entities.User{
		ID:       accountID,
		Nickname: body.Nickname,
		Mobile:   body.Mobile,
		Avatar:   body.Avatar,
	})
	if err != nil {
		return nil, err
	}

	return nil, nil
}

func (con *Controller) GetUserInfo(c *gin.Context) (interface{}, error) {
	accountID := c.MustGet(pkg.AccountIDKey).(int64)

	a, err := con.Service.GetUserByID(accountID)
	if err != nil {
		return nil, err
	}

	return &entities.UserDTO{
		ID:       a.ID,
		Nickname: a.Nickname,
		Mobile:   a.Mobile,
		Avatar:   a.Avatar,
		Uid:      a.Uid,
	}, nil
}

func (con *Controller) GetUserInfos(c *gin.Context) (interface{}, error) {
	param, exists := c.GetQuery("id")
	if !exists {
		return nil, response.NewIllegalArgumentError("invalid param")
	}

	log.Infof("get user infos; ids=%s", param)
	var ids []int64
	for _, e := range strings.Split(param, ",") {
		i, err := strconv.ParseInt(e, 10, 64)
		if err != nil {
			return nil, response.NewIllegalArgumentError(err.Error())
		}

		ids = append(ids, i)
	}

	accounts, err := con.Service.ListUserByIDs(ids)
	if err != nil {
		return nil, response.NewIllegalArgumentError(err.Error())
	}

	dto := make([]*entities.SimpleUser, 0)
	for _, account := range accounts {
		dto = append(dto, &entities.SimpleUser{ID: account.ID, Nickname: account.Nickname, Avatar: account.Avatar})
	}

	return map[string]interface{}{"users": dto}, nil
}

func (con *Controller) CheckUserToken(c *gin.Context) (interface{}, error) {
	return nil, nil
}

func (con *Controller) Login(context *gin.Context) (interface{}, error) {
	var loginCommand struct {
		Username          string `json:"username" binding:"required"`
		EncryptedPassword string `json:"encryptedPassword" binding:"required"`
	}

	err := context.ShouldBindJSON(&loginCommand)
	if err != nil {
		return nil, err
	}

	encryptedPassword := fmt.Sprintf("%x", md5.Sum([]byte(con.backendConf.Password)))
	fmt.Printf("%s", encryptedPassword)
	if strings.ToLower(loginCommand.EncryptedPassword) == encryptedPassword {
		token, err := SignToken(con.backendConf.Id)
		if err != nil {
			return nil, err
		}

		return map[string]interface{}{
			"id":    con.backendConf.Id,
			"token": token,
		}, nil
	}

	return nil, errors.New("wrong username or password")
}

func (con *Controller) GenerateUid(context *gin.Context) (interface{}, error) {
	err := con.Service.GenerateUidForExistsUser()
	if err != nil {
		return nil, err
	}

	return nil, nil
}

// JTime 时间支持ISO、时间戳、时间戳字符串 的json解析
type JTime time.Time

func (t *JTime) UnmarshalJSON(str []byte) error {
	var date time.Time
	err := json.Unmarshal(str, &date)
	if err != nil {
		var n int64
		err = json.Unmarshal(str, &n)
		if err != nil {
			var mSecStr string
			err = json.Unmarshal(str, &mSecStr)
			if err != nil {
				return err
			}
			n, err = strconv.ParseInt(mSecStr, 10, 64)
			if err != nil {
				return err
			}
		}
		date = time.UnixMilli(n)
	}
	*t = JTime(date)
	return nil
}

func (con *Controller) DealReportedEvents(context *gin.Context) (interface{}, error) {
	type UserEvent struct {
		Key    string        `json:"eventType"`
		Time   JTime         `json:"time"`
		Params []interface{} `json:"params"`
	}

	var event UserEvent
	if err := context.ShouldBindJSON(&event); err != nil {
		return nil, err
	}
	log.Infof("time %v\n", time.Time(event.Time))
	userId := context.MustGet(pkg.AccountIDKey).(int64)

	con.Service.DealUserEvents(userId, event.Key, event.Params)
	return nil, nil
}

func (con *Controller) CountRegisteredUsers(c *gin.Context) (interface{}, error) {
	countUser, err := con.Service.CountUser()
	if err != nil {
		return nil, nil
	}

	return map[string]interface{}{
		"count": countUser,
	}, nil
}

func (con *Controller) GetUserCharityCard(c *gin.Context) (interface{}, error) {
	userId := c.MustGet(pkg.AccountIDKey).(int64)
	url, err := con.Service.GetUserCharityCard(userId)
	if err != nil {
		return nil, err
	}
	return map[string]interface{}{
		"url": url,
	}, nil
}
func (con *Controller) ReportSubscribeMessageSettings(c *gin.Context) (interface{}, error) {
	type Request struct {
		Key           string                               `json:"key"`
		Templates     []*entities.SubscribeTemplateSetting `json:"templates"`
		Subscriptions *entities.SubscriptionsSetting       `json:"subscriptions"`
	}

	var req Request
	if err := c.ShouldBindJSON(&req); err != nil {
		return nil, err
	}

	err := interfaces.S.SubscribeMsg.Report(c.MustGet(pkg.AccountIDKey).(int64), req.Key, req.Templates, req.Subscriptions)
	if err != nil {
		return nil, err
	}

	return nil, nil
}

func (con *Controller) GetLastReportSubscribeMessageSettings(c *gin.Context) (interface{}, error) {
	key, b := c.GetQuery("key")
	if !b {
		return nil, errors.New("key cannot be empty")
	}

	templates, subscriptionsSetting, at, err := interfaces.S.SubscribeMsg.GetLastReport(c.MustGet(pkg.AccountIDKey).(int64), key)
	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"reportAt":      at,
		"templates":     templates,
		"subscriptions": subscriptionsSetting,
	}, nil
}
