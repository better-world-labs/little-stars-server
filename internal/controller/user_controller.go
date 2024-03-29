package controller

import (
	"aed-api-server/internal/interfaces"
	"aed-api-server/internal/interfaces/entities"
	"aed-api-server/internal/interfaces/service"
	"aed-api-server/internal/pkg"
	"aed-api-server/internal/pkg/location"
	page "aed-api-server/internal/pkg/query"
	"aed-api-server/internal/pkg/response"
	"aed-api-server/internal/pkg/utils"
	"aed-api-server/internal/server/config"
	"aed-api-server/internal/service/user"
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

type UserController struct {
	Service     service.UserServiceOld `inject:"-"`
	User        service.UserService    `inject:"-"`
	BackendConf config.Backend         `conf:"backend"`
}

//go:inject-component
func NewUserController() *UserController {
	return &UserController{}
}

func (con *UserController) MountAuthRouter(r *route.Router) {
	r.PUT("/user/position", con.UpdatePosition)
	r.PUT("/user/mobile", con.UpdateMobile)
	r.PUT("/user/info", con.UpdateUserInfo)
	r.GET("/user/info", con.GetUserInfo)
	r.GET("/user/infos", con.GetUserInfos)
	r.POST("/user/check-token", con.CheckUserToken)
	r.GET("/user/id", con.GetUserId)
	r.POST("/user/events", con.DealReportedEvents)
	r.GET("/user/charity-card", con.GetUserCharityCard)
	r.GET("/user/stat", con.GetUserStat)
	r.POST("/user/message-subscribe", con.ReportSubscribeMessageSettings)
	r.GET("/user/message-subscribe/last", con.GetLastReportSubscribeMessageSettings)
}

func (con *UserController) MountAdminRouter(r *route.Router) {
	r.POST("/user/update-avatar", con.UpdateAvatar)
	r.GET("/users", con.ListUsers)
}

func (con *UserController) MountNoAuthRouter(r *route.Router) {
	r.POST("/user/wechat/app/login", con.WechatAppLogin)
	r.POST("/user/wechat/mini-program/login", con.WechatMiniProgramLogin)
	r.POST("/v2/user/wechat/mini-program/login", con.WechatMiniProgramLoginV2)
	r.POST("/user/wechat/mini-program/login-simple", con.WechatMiniProgramLoginSimple)
	r.POST("/user/generate-uid", con.GenerateUid)
	r.GET("/user-counts", con.CountRegisteredUsers)
}

func (con *UserController) MountGinEngineRouter(r *route.Router) {
	r.POST("/admin-api/user/login", con.Login)
}

func (con *UserController) WechatAppLogin(c *gin.Context) (interface{}, error) {
	return nil, errors.New("nod implements")
}

func (con *UserController) WechatMiniProgramLogin(c *gin.Context) (interface{}, error) {
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

func (con *UserController) WechatMiniProgramLoginV2(c *gin.Context) (interface{}, error) {
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

func (con *UserController) WechatMiniProgramLoginSimple(c *gin.Context) (interface{}, error) {
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

func (con *UserController) UpdatePosition(c *gin.Context) (interface{}, error) {
	accountID := c.MustGet(pkg.AccountIDKey).(int64)
	var body entities.Position
	body.AccountID = accountID
	err := c.ShouldBindJSON(&body)
	if err != nil {
		return nil, err
	}

	err = con.User.UpdatePosition(&body)
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

func (con *UserController) UpdateMobile(c *gin.Context) (interface{}, error) {
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

func (con *UserController) UpdateUserInfo(c *gin.Context) (interface{}, error) {
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

func (con *UserController) GetUserInfo(c *gin.Context) (interface{}, error) {
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

func (con *UserController) GetUserInfos(c *gin.Context) (interface{}, error) {
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

func (con *UserController) CheckUserToken(c *gin.Context) (interface{}, error) {
	userId := utils.GetContextUserId(c)
	u, err := con.User.GetUserInfo(userId)
	if err != nil || u == nil || u.SessionKey == "" {
		return nil, response.ErrorInvalidToken
	}

	return map[string]interface{}{
		"userId": c.MustGet(pkg.AccountIDKey),
	}, nil
}
func (con *UserController) GetUserId(c *gin.Context) (interface{}, error) {
	return map[string]interface{}{
		"userId": c.MustGet(pkg.AccountIDKey),
	}, nil
}

func (con *UserController) Login(context *gin.Context) (interface{}, error) {
	var loginCommand struct {
		Username          string `json:"username" binding:"required"`
		EncryptedPassword string `json:"encryptedPassword" binding:"required"`
	}

	err := context.ShouldBindJSON(&loginCommand)
	if err != nil {
		return nil, err
	}

	encryptedPassword := fmt.Sprintf("%x", md5.Sum([]byte(con.BackendConf.Password)))
	fmt.Printf("%s", encryptedPassword)
	if strings.ToLower(loginCommand.EncryptedPassword) == encryptedPassword {
		token, err := user.SignToken(con.BackendConf.Id)
		if err != nil {
			return nil, err
		}

		return map[string]interface{}{
			"id":    con.BackendConf.Id,
			"token": token,
		}, nil
	}

	return nil, errors.New("wrong username or password")
}

func (con *UserController) GenerateUid(context *gin.Context) (interface{}, error) {
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

func (con *UserController) DealReportedEvents(context *gin.Context) (interface{}, error) {
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

	con.User.DealUserReportEvents(userId, event.Key, event.Params)
	return nil, nil
}

func (con *UserController) CountRegisteredUsers(c *gin.Context) (interface{}, error) {
	countUser, err := con.Service.CountUser()
	if err != nil {
		return nil, nil
	}

	return map[string]interface{}{
		"count": countUser,
	}, nil
}

func (con *UserController) GetUserCharityCard(c *gin.Context) (interface{}, error) {
	userId := c.MustGet(pkg.AccountIDKey).(int64)
	url, err := con.Service.GetUserCharityCard(userId)
	if err != nil {
		return nil, err
	}
	return map[string]interface{}{
		"url": url,
	}, nil
}
func (con *UserController) ReportSubscribeMessageSettings(c *gin.Context) (interface{}, error) {
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

func (con *UserController) GetLastReportSubscribeMessageSettings(c *gin.Context) (interface{}, error) {
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

func (con *UserController) GetUserStat(ctx *gin.Context) (interface{}, error) {
	userId := utils.GetContextUserId(ctx)
	stat, err := con.User.GetUserAboutStat(userId)
	if err != nil {
		return nil, err
	}

	return stat, nil
}

func (con *UserController) UpdateAvatar(ctx *gin.Context) (interface{}, error) {
	count, err := con.Service.UpdateUsersAvatar()
	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"updated": count,
	}, nil
}

func (con *UserController) ListUsers(ctx *gin.Context) (interface{}, error) {
	p, err := page.BindPageQuery(ctx)
	if err != nil {
		return nil, err
	}

	keyword := ctx.Query("keyword")
	users, err := con.User.PageUsers(*p, keyword)
	if err != nil {
		return nil, err
	}

	return users, nil
}
