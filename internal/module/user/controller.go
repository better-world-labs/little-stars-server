package user

import (
	"aed-api-server/internal/interfaces"
	"aed-api-server/internal/interfaces/entities"
	"aed-api-server/internal/pkg"
	"aed-api-server/internal/pkg/config"
	"aed-api-server/internal/pkg/global"
	"aed-api-server/internal/pkg/location"
	"aed-api-server/internal/pkg/response"
	"aed-api-server/internal/pkg/utils"
	"crypto/md5"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"gitlab.openviewtech.com/openview-pub/gopkg/log"
	"net/http"
	"strconv"
	"strings"
	"time"
)

type Controller struct {
	service     Service
	backendConf config.Backend
}

func NewController(backend config.Backend, s Service) *Controller {
	return &Controller{backendConf: backend, service: s}
}

func (con *Controller) WechatAppLogin(c *gin.Context) {
	//var body LoginDTO
	//err := c.BindJSON(&body)
	//if err != nil {
	//	log.DefaultLogger().Errorf("bind json error: %v", err)
	//	response.ReplyError(c, err)
	//	return
	//}
	//
	//res, token, err := con.service.WechatAppLogin(body.Code)
	//if err != nil {
	//	log.DefaultLogger().Errorf("do login error: %v", err)
	//	response.ReplyError(c, err)
	//	return
	//}

	response.ReplyError(c, errors.New("not implements"))
}

func (con *Controller) WechatMiniProgramLogin(c *gin.Context) {
	var body LoginCommand
	err := c.BindJSON(&body)
	log.DefaultLogger().Infof("mobileCode: %s", body.MobileCode)

	if err != nil {
		log.DefaultLogger().Errorf("bind json error: %v", err)
		response.ReplyError(c, err)
		return
	}

	res, token, _, err := con.service.WechatMiniProgramLogin(body)
	if err != nil {
		log.DefaultLogger().Errorf("do login error: %v", err)
		response.ReplyError(c, err)
		return
	}

	log.Info("Login: uid=", res.Uid)
	response.ReplyOK(c, &entities.UserDTO{
		ID:       res.ID,
		Uid:      res.Uid,
		Nickname: res.Nickname,
		Token:    token,
		Avatar:   res.Avatar,
		Mobile:   res.Mobile,
		Openid:   res.Openid,
	})
}

func (con *Controller) WechatMiniProgramLoginSimple(c *gin.Context) {
	var body SimpleLoginCommand
	err := c.BindJSON(&body)

	if err != nil {
		log.DefaultLogger().Errorf("bind json error: %v", err)
		response.ReplyError(c, err)
		return
	}

	res, token, _, err := con.service.WechatMiniProgramLoginSimple(body)
	if err != nil {
		log.DefaultLogger().Errorf("do login error: %v", err)
		response.ReplyError(c, err)
		return
	}

	log.Info("SimpleLogin: uid=", res.Uid)
	response.ReplyOK(c, &entities.UserDTO{
		ID:       res.ID,
		Uid:      res.Uid,
		Nickname: res.Nickname,
		Token:    token,
		Avatar:   res.Avatar,
		Mobile:   res.Mobile,
		Openid:   res.Openid,
	})
}

func (con *Controller) UpdatePosition(c *gin.Context) {
	accountID := c.MustGet(pkg.AccountIDKey).(int64)
	var body Position
	body.AccountID = accountID
	err := c.BindJSON(&body)
	if err != nil {
		log.DefaultLogger().Errorf("bind json error: %v", err)
		response.ReplyError(c, err)
		return
	}

	err = con.service.UpdatePosition(&body)
	if err != nil {
		log.DefaultLogger().Errorf("update position error: %v", err)
		response.ReplyError(c, err)
		return
	}

	//调用任务接口生成任务
	go interfaces.S.Task.GenJobsByUserLocation(accountID, location.Coordinate{
		Latitude:  body.Latitude,
		Longitude: body.Longitude,
	})

	response.ReplyOK(c, nil)
}

func (con *Controller) UpdateMobile(c *gin.Context) {
	accountID := c.MustGet(pkg.AccountIDKey).(int64)
	var body MobileDTO
	err := c.BindJSON(&body)
	if err != nil {
		log.DefaultLogger().Errorf("bind json error: %v", err)
		response.ReplyApiError(c, http.StatusBadRequest, response.StatusInvalidParam, err.Error())
		return
	}

	err = con.service.UpdateMobile(accountID, body.Mobile)
	if err != nil {
		response.ReplyError(c, err)
		return
	}

	response.ReplyOK(c, nil)
}

func (con *Controller) UpdateUserInfo(c *gin.Context) {
	accountID := c.MustGet(pkg.AccountIDKey).(int64)

	var body entities.UserDTO
	err := c.BindJSON(&body)
	if err != nil {
		log.DefaultLogger().Errorf("bind json error: %v", err)
		response.ReplyApiError(c, http.StatusBadRequest, response.StatusInvalidParam, err.Error())
		return
	}

	err = con.service.UpdateUserInfo(&User{
		ID:       accountID,
		Nickname: body.Nickname,
		Mobile:   body.Mobile,
		Avatar:   body.Avatar,
	})
	if err != nil {
		response.ReplyError(c, err)
		return
	}

	response.ReplyOK(c, nil)
}

func (con *Controller) GetUserInfo(c *gin.Context) {
	accountID := c.MustGet(pkg.AccountIDKey).(int64)

	a, err := con.service.GetUserByID(accountID)
	utils.MustNil(err, global.ErrorUnknown)

	response.ReplyOK(c, &entities.UserDTO{
		ID:       a.ID,
		Nickname: a.Nickname,
		Mobile:   a.Mobile,
		Avatar:   a.Avatar,
		Uid:      a.Uid,
	})
}

func (con *Controller) GetUserInfos(c *gin.Context) {
	param, exists := c.GetQuery("id")
	utils.MustTrue(exists, global.ErrorInvalidParam)

	log.DefaultLogger().Infof("get user infos; ids=%s", param)
	var ids []int64
	for _, e := range strings.Split(param, ",") {
		i, err := strconv.ParseInt(e, 10, 64)
		utils.MustNil(err, global.ErrorInvalidParam)
		ids = append(ids, i)
	}

	accounts, err := con.service.ListUserByIDs(ids)
	utils.MustNil(err, global.ErrorUnknown)

	dto := make([]*entities.SimpleUser, 0)
	for _, account := range accounts {
		dto = append(dto, &entities.SimpleUser{ID: account.ID, Nickname: account.Nickname, Avatar: account.Avatar})
	}

	response.ReplyOK(c, map[string]interface{}{"users": dto})
}

func (con *Controller) CheckUserToken(c *gin.Context) {
	response.ReplyOK(c, nil)
}

func (con *Controller) Login(context *gin.Context) {
	var loginCommand struct {
		Username          string `json:"username" binding:"required"`
		EncryptedPassword string `json:"encryptedPassword" binding:"required"`
	}

	err := context.BindJSON(&loginCommand)
	if err != nil {
		response.ReplyError(context, err)
		return
	}

	encryptedPassword := fmt.Sprintf("%x", md5.Sum([]byte(con.backendConf.Password)))
	if strings.ToLower(loginCommand.EncryptedPassword) == encryptedPassword {
		token, err := SignToken(con.backendConf.Id)
		if err != nil {
			response.ReplyError(context, err)
			return
		}

		response.ReplyOK(context, map[string]interface{}{
			"id":    con.backendConf.Id,
			"token": token,
		})
		return
	}

	response.ReplyError(context, errors.New("wrong username or password"))
}

func (con *Controller) GenerateUid(context *gin.Context) {
	err := con.service.GenerateUidForExistsUser()
	if err != nil {
		response.ReplyError(context, err)
		return
	}

	response.ReplyOK(context, nil)
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

func (con *Controller) DealReportedEvents(context *gin.Context) {
	type UserEvent struct {
		EventType string `json:"eventType"`
		Time      JTime  `json:"time"`
	}

	var event UserEvent
	if err := context.BindJSON(&event); err != nil {
		response.ReplyError(context, err)
		return
	}
	log.DefaultLogger().Infof("time %v\n", time.Time(event.Time))
	userId := context.MustGet(pkg.AccountIDKey).(int64)

	con.service.DealUserEvents(userId, event.EventType)
	response.ReplyOK(context, nil)
}

func (con *Controller) CountRegisteredUsers(c *gin.Context) {
	countUser, err := con.service.CountUser()
	if err != nil {
		response.ReplyError(c, err)
		return
	}

	response.ReplyOK(c, map[string]interface{}{
		"count": countUser,
	})
}

func (con *Controller) GetUserCharityCard(c *gin.Context) {
	userId := c.MustGet(pkg.AccountIDKey).(int64)
	url, err := con.service.GetUserCharityCard(userId)
	if err != nil {
		response.ReplyError(c, err)
		return
	}
	response.ReplyOK(c, map[string]interface{}{
		"url": url,
	})
}
