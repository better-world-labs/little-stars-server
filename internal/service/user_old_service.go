package service

import (
	"aed-api-server/internal/interfaces"
	"aed-api-server/internal/interfaces/entities"
	"aed-api-server/internal/interfaces/events"
	service2 "aed-api-server/internal/interfaces/service"
	"aed-api-server/internal/module/user"
	"aed-api-server/internal/pkg/cache"
	"aed-api-server/internal/pkg/crypto"
	"aed-api-server/internal/pkg/db"
	"aed-api-server/internal/pkg/domain/emitter"
	"aed-api-server/internal/pkg/response"
	"aed-api-server/internal/pkg/utils"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
	"time"

	log "github.com/sirupsen/logrus"
)

const CacheName = "user"

type userServiceOld struct {
	Client        user.WechatClient             `inject:"-"`
	TreasureChest service2.TreasureChestService `inject:"-"`
}

func NewUserServiceOld() service2.UserServiceOld {
	return &userServiceOld{}
}

func (s *userServiceOld) WechatAppLogin(code string) (*entities.User, string, error) {
	var info user.OAuthInfoDTO
	if err := s.Client.GetAccessToken(code, &info.OAuthAccessTokenDTO); err != nil {
		return nil, "", err
	}
	if err := s.Client.GetUserInfo(info.AccessToken, info.OpenID, &info); err != nil {
		return nil, "", err
	}

	account := entities.User{
		Nickname: info.Nickname,
		Avatar:   info.HeadImgURL,
		Unionid:  info.UnionID,
		Openid:   info.OpenID,
	}
	err := s.createOrUpdateUser(&account)
	if err != nil {
		return nil, "", err
	}

	token, err := user.SignToken(account.ID)
	if err != nil {
		return nil, "", err
	}

	return &account, token, err
}

func (s *userServiceOld) WechatMiniProgramLogin(command entities.LoginCommand) (*entities.User, string, string, error) {
	var info user.MiniProgramResponseDTO
	log.Infof("encryptPhone: %s, iv: %s", command.EncryptPhone, command.Iv)
	if err := s.Client.MiniProgramCode2Session(command.Code, command.EncryptPhone, command.Iv, &info); err != nil {
		return nil, "", "", err
	}
	account := entities.User{
		Nickname: command.Nickname,
		Avatar:   command.Avatar,
		Mobile:   info.DecryptedPhone,
		Unionid:  info.UnionID,
		Openid:   info.OpenID,
	}

	err := s.createOrUpdateUser(&account)
	if err != nil {
		return nil, "", "", err
	}

	token, err := user.SignToken(account.ID)
	if err != nil {
		return nil, "", "", err
	}

	return &account, token, info.SessionKey, err
}

func (s *userServiceOld) WechatMiniProgramLoginV2(command entities.LoginCommandV2) (*entities.User, string, string, error) {
	session, err := s.Client.CodeToSession(command.Code)
	if err != nil {
		return nil, "", "", err
	}

	account := entities.User{
		Nickname: command.Nickname,
		Avatar:   command.Avatar,
		Unionid:  session.UnionId,
		Openid:   session.Openid,
	}

	err = s.createOrUpdateUser(&account)
	if err != nil {
		return nil, "", "", err
	}

	token, err := user.SignToken(account.ID)
	if err != nil {
		return nil, "", "", err
	}

	return &account, token, session.SessionKey, err
}

func (s *userServiceOld) GetUserByOpenId(openid string) (*entities.User, bool, error) {
	var u entities.User
	exists, err := db.Table("account").Where("openid = ?", openid).Get(&u)
	return &u, exists, err
}

func (s *userServiceOld) createOrUpdateUser(account *entities.User) error {
	session := db.GetSession()
	defer session.Close()

	var current entities.User
	exists, err := session.Table("account").Where("openid = ?", account.Openid).Get(&current)
	if err != nil {
		return err
	}

	if exists {
		account.ID = current.ID
		account.Uid = current.Uid
		log.Info("account_exists: uid=", account.Uid)
		_, err := session.Table("account").ID(account.ID).Update(account)
		return err
	} else {
		account.Created = time.Now()
		account.Uid = s.generateUid()
		_, err := session.Table("account").Insert(account)
		if err != nil {
			return err
		}

		return emitter.Emit(&events.FirstLoginEvent{
			UserId:  account.ID,
			Openid:  account.Openid,
			LoginAt: account.Created,
		})
	}
}

func (s *userServiceOld) UpdatePosition(position *entities.Position) error {
	session := db.GetSession()
	defer session.Close()

	var current entities.Position
	exists, err := session.Table("account_position").
		Where("account_id = ?", position.AccountID).
		Get(&current)
	if err != nil {
		return err
	}

	if !exists {
		_, err := session.Table("account_position").Insert(position)
		return err
	} else {
		_, err := session.Table("account_position").ID(current.ID).Update(position)
		return err
	}
}

func (s *userServiceOld) ListPositionByUserIDs(accountIDs []int64) (map[int64]*entities.Position, error) {
	session := db.GetSession()
	defer session.Close()

	res := make(map[int64]*entities.Position, 0)
	arr := make([]*entities.Position, 0)
	cond := new(entities.Position)
	cond.AccountID = accountIDs[0]
	err := session.Table("account_position").In("account_id", accountIDs).Find(&arr)
	if err != nil {
		return nil, err
	}

	for i := range arr {
		res[arr[i].AccountID] = arr[i]
	}

	return res, nil
}

func (s *userServiceOld) ListAllUser() ([]*entities.User, error) {
	session := db.GetSession()
	defer session.Close()

	arr := make([]*entities.User, 0)
	err := session.Table("account").Find(&arr)
	if err != nil {
		return nil, err
	}

	return arr, nil
}

func (s *userServiceOld) ListUserByIDs(ids []int64) (r map[int64]*entities.User, err error) {
	return s.ListUserByIDsFromDB(ids)
}

func (s *userServiceOld) ListUserByIDsFromDB(ids []int64) (map[int64]*entities.User, error) {
	session := db.GetSession()
	defer session.Close()

	res := make(map[int64]*entities.User, 0)
	err := session.Table("account").In("id", ids).Find(&res)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (s *userServiceOld) GetUserByID(id int64) (account *entities.User, err error) {
	o, exists, err := cache.GetManager().Get(CacheName, getCacheKey(id))
	if err == nil && exists {
		return o.(*entities.User), nil
	}

	o, err = s.getUserByIDFromDB(id)
	if err != nil {
		return nil, err
	}

	_ = cache.GetManager().Put(CacheName, getCacheKey(id), o)
	return o.(*entities.User), err
}

func (s *userServiceOld) getUserByIDFromDB(id int64) (account *entities.User, err error) {
	session := db.GetSession()
	defer session.Close()

	account = &entities.User{}
	ok, err := session.Table("account").Where("id = ?", id).Get(account)
	if err != nil {
		return nil, err
	}

	if !ok {
		return nil, nil
	}

	return account, nil
}

func (s *userServiceOld) GetUserByIDAsync(id int64) func() (account *entities.User, err error) {
	resChan := make(chan interface{}, 1)

	utils.Go(func() {
		defer close(resChan)
		res, err := s.GetUserByID(id)
		if err == nil {
			resChan <- res
		} else {
			resChan <- err
		}
	})

	return func() (account *entities.User, err error) {
		res := <-resChan
		switch res.(type) {
		case *entities.User:
			return res.(*entities.User), nil
		case error:
			return nil, res.(error)
		default:
			return nil, errors.New("invalid type from chan")
		}
	}
}

func (s *userServiceOld) GetPositionByUserID(accountID int64) (*entities.Position, error) {
	session := db.GetSession()
	defer session.Close()

	position := &entities.Position{}
	ok, err := session.Table("account_position").Where("account_id = ?", accountID).Get(position)
	if err != nil {
		return nil, err
	}

	if !ok {
		return nil, nil
	}

	return position, nil
}

func (s *userServiceOld) UpdateMobile(accountID int64, command entities.MobileCommand) error {
	weSession, err := s.Client.CodeToSession(command.Code)
	if err != nil {
		return err
	}

	var phone crypto.WxUserPhone
	err = s.Client.Decrypt(command.EncryptPhone, command.Iv, weSession.SessionKey, &phone)
	if err != nil {
		return err
	}

	session := db.GetSession()
	defer session.Close()

	return db.WithTransaction(session, func() error {
		var current entities.User
		exists, err := session.Table("account").Where("id = ?", accountID).Get(&current)
		if err != nil {
			return err
		}

		if !exists {
			return errors.New("account does not exists")
		}

		current.Mobile = phone.PhoneNumber

		_, err = session.Table("account").ID(current.ID).Update(&current)
		if err != nil {
			return err
		}

		return cache.GetManager().Evict(CacheName, getCacheKey(accountID))
	})
}

func (s *userServiceOld) ListAllUserAsync() func() ([]*entities.User, error) {
	resChan := make(chan interface{}, 1)
	utils.Go(func() {
		defer close(resChan)
		accounts, err := s.ListAllUser()
		if err != nil {
			resChan <- err
		} else {
			resChan <- accounts
		}
	})

	return func() (accounts []*entities.User, err error) {
		res := <-resChan
		switch res.(type) {
		case []*entities.User:
			return res.([]*entities.User), nil
		case error:
			return nil, res.(error)
		default:
			return nil, errors.New("invalid type from chan")
		}
	}
}

func (s *userServiceOld) ListAllPositions() ([]*entities.Position, error) {
	session := db.GetSession()
	defer session.Close()

	res := make([]*entities.Position, 0)
	err := session.Table("account_position").Find(&res)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (s *userServiceOld) ListAllPositionsAsync() func() ([]*entities.Position, error) {
	resChan := make(chan interface{}, 1)
	utils.Go(func() {
		defer close(resChan)
		accounts, err := s.ListAllPositions()
		if err != nil {
			resChan <- err
		} else {
			resChan <- accounts
		}
	})

	return func() (accounts []*entities.Position, err error) {
		res := <-resChan
		switch res.(type) {
		case []*entities.Position:
			return res.([]*entities.Position), nil
		case error:
			return nil, res.(error)
		default:
			return nil, errors.New("invalid type from chan")
		}
	}
}

func (s *userServiceOld) ListPositionByUserIDsAsync(
	accountIDs []int64,
) func() (map[int64]*entities.Position, error) {
	resChan := make(chan interface{}, 1)
	utils.Go(func() {
		defer close(resChan)
		accounts, err := s.ListPositionByUserIDs(accountIDs)
		if err != nil {
			resChan <- err
		} else {
			resChan <- accounts
		}
	})

	return func() (accounts map[int64]*entities.Position, err error) {
		res := <-resChan
		switch res.(type) {
		case map[int64]*entities.Position:
			return res.(map[int64]*entities.Position), nil
		case error:
			return nil, res.(error)
		default:
			return nil, errors.New("invalid type from chan")
		}
	}
}

func (s *userServiceOld) UpdateUserInfo(account *entities.User) error {
	session := db.GetSession()
	defer session.Close()

	return db.WithTransaction(session, func() error {
		var current entities.User
		exists, err := session.Table("account").Where("id = ?", account.ID).Get(&current)
		if err != nil {
			return err
		}

		if !exists {
			return errors.New("account does not exists")
		}

		_, err = session.ID(account.ID).Update(account)

		return cache.GetManager().Evict(CacheName, getCacheKey(account.ID))
	})
}

func getCacheKey(id int64) string {
	return fmt.Sprintf("id:%d", id)
}

func (s *userServiceOld) generateUid() string {
	u := uuid.New()
	formated := strconv.FormatUint(uint64(u.ID()), 10)
	l := len(formated)
	if l < 10 {
		for i := 0; i < 10-i; i++ {
			formated += "0"
		}
	}
	return formated[:10]
}

func (s *userServiceOld) GenerateUidForExistsUser() error {
	users, err := s.ListAllUser()
	if err != nil {
		return err
	}

	for _, u := range users {
		if u.Uid == "" {
			err := s.UpdateUid(u)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func (s *userServiceOld) DealUserEvents(userId int64, key string, params []interface{}) {
	utils.Go(func() {
		interfaces.S.User.RecordUserEvent(userId, entities.GetUserEventTypeOfReport(key), params...)
	})

	utils.Go(func() {
		switch key {
		case entities.Report_enterAedMap:
			dealUserEnterAEDMap(userId)

		case entities.Report_readNews:
			dealUserReadNews(userId)

		case entities.Report_showSubscribeQrCode:
			dealShowSubscribeQrCode(userId)

		case entities.Report_openTreasureChest:
			if len(params) == 0 {
				return
			}
			arg := params[0].(map[string]interface{})
			id, ok := arg["treasureChestId"]
			if !ok {
				return
			}
			err := s.TreasureChest.OpenTreasureChest(userId, int(id.(float64)))
			if err != nil {
				log.Error("OpenTreasureChest err", err)
			}
		case entities.Report_scanPage:
		case entities.Report_scanVideo:
		}
	})

}

func dealShowSubscribeQrCode(userId int64) {
	log.Info("dealShowSubscribeQrCode", "userId=", userId)
	var isSet bool
	err := interfaces.S.UserConfig.GetConfigToValue(userId, entities.UserConfigKeySubscribeOfficialAccounts, &isSet)
	log.Info("GetConfigToValue", "userId=", userId, "key=", entities.UserConfigKeySubscribeOfficialAccounts, "value=", isSet)
	if err != nil {
		log.Info("dealUserEnterAEDMap error", err)
		return
	}

	if !isSet {
		log.Info("PutValueToConfig", "userId=", userId)
		updated, err := interfaces.S.UserConfig.PutValueToConfig(userId, entities.UserConfigKeySubscribeOfficialAccounts, true)
		if updated && err == nil {
			log.Info("Emit", "userId=", userId)
			err = emitter.Emit(&events.PointsEvent{
				PointsEventType: entities.PointsEventTypeSubscribe,
				UserId:          userId,
			})
		}

		if err != nil {
			log.Error("emit event error", err)
		}
	}
}

func genImgCommand(
	username string,
	userAvatar string,
	medals []string,
	donationPoints int,
	donationProject int,
	addStarDays int,
	qrContent string,
	saveAs string,
) (url string, err error) {
	config := interfaces.GetConfig()

	m := map[string]interface{}{
		"tplName": "user-charity-card",
		"args": map[string]interface{}{
			"username":        username,
			"userAvatar":      userAvatar,
			"medals":          medals,
			"donationPoints":  donationPoints,
			"donationProject": donationProject,
			"addStarDays":     addStarDays,
			"qrContent":       qrContent,
		},
		"save": saveAs,
	}

	jsonStr, _ := json.Marshal(m)
	payload := strings.NewReader(string(jsonStr))
	req, _ := http.NewRequest("POST", config.ImgBotService, payload)
	req.Header.Add("Content-Type", "application/json")

	res, err := http.DefaultClient.Do(req)

	if err != nil {
		return "", err
	}

	all, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return "", err
	}

	rst := map[string]interface{}{}

	err = json.Unmarshal(all, &rst)

	errCode, suc := rst["code"].(float64)
	var errMsg = "unknown error in img-bot"
	if !suc {
		return "", errors.New(errMsg)
	}

	if errCode != 200 {
		errMsg, suc = rst["message"].(string)
		if !suc {
			errMsg = "unknown error in img-bot"
		}
		return "", errors.New(errMsg)
	}
	return rst["url"].(string), nil
}

func (s *userServiceOld) GetUserCharityCard(userId int64) (string, error) {
	config := interfaces.GetConfig()

	rst, err := utils.PromiseAll(func() (interface{}, error) {
		return s.GetUserByID(userId)
	}, func() (interface{}, error) {
		return interfaces.S.UserMedal.GetUserMedalUrl(userId)
	}, func() (interface{}, error) {
		return interfaces.S.Donation.StatDonationByUserId(userId)
	})

	if err != nil {
		return "", err
	}

	user := rst[0].(*entities.User)
	modalUrls := rst[1].([]string)
	donationStat := rst[2].(entities.DonationStat)

	days := time.Now().Sub(user.Created).Hours() / 24

	url, err := genImgCommand(
		user.Nickname,
		user.Avatar,
		modalUrls,
		donationStat.DonationTotalPoints,
		donationStat.DonationProjectCount,
		int(days),
		fmt.Sprintf(`https://%s/share/cert?source=charity-card&sharer=%v`, config.Server.Host, userId),
		fmt.Sprintf("user-charity-card-%v.jpeg", userId),
	)
	url += fmt.Sprintf("?%d", time.Now().Unix())
	return url, err
}

func dealUserReadNews(userId int64) {
	has, err := interfaces.S.TaskBubble.HasReadNewsTask(userId)
	if err != nil {
		log.Error("TaskBubble.HasReadNewsTask error", err)
		return
	}

	if !has {
		return
	}

	err = emitter.Emit(&events.PointsEvent{
		PointsEventType: entities.PointsEventTypeReadNews,
		UserId:          userId,
	})
	if err != nil {
		log.Error("emit event error", err)
	}
}

func (s *userServiceOld) UpdateUid(user *entities.User) error {
	uid := s.generateUid()
	user.Uid = uid
	_, err := db.Table("account").ID(user.ID).Update(user)
	return err
}

func dealUserEnterAEDMap(userId int64) {
	var isSet bool
	err := interfaces.S.UserConfig.GetConfigToValue(userId, entities.UserConfigKeyFirstEnterAEDMap, &isSet)
	if err != nil {
		log.Info("dealUserEnterAEDMap error", err)
		return
	}

	if !isSet {
		updated, err := interfaces.S.UserConfig.PutValueToConfig(userId, entities.UserConfigKeyFirstEnterAEDMap, true)
		if updated && err == nil {
			err = emitter.Emit(&events.PointsEvent{
				PointsEventType: entities.PointsEventTypeFirstEnterAEDMap,
				UserId:          userId,
				Params: entities.PointsEventParams{
					RefTable:   "user_config#" + entities.UserConfigKeyFirstEnterAEDMap,
					RefTableId: userId,
				},
			})
		}
	}
}

func (s *userServiceOld) WechatMiniProgramLoginSimple(command entities.SimpleLoginCommand) (*entities.User, string, string, error) {
	session, err := s.Client.CodeToSession(command.Code)
	if err != nil {
		return nil, "", "", err
	}

	u, exists, err := s.GetUserByOpenId(session.Openid)
	if err != nil {
		return nil, "", "", err
	}

	if !exists {
		return nil, "", "", response.ErrorUserNotRegister
	}

	token, err := user.SignToken(u.ID)
	if err != nil {
		return nil, "", "", err
	}

	return u, token, session.SessionKey, nil
}

func (s *userServiceOld) CountUser() (int, error) {
	count, err := db.Table("account").Count()
	if err != nil {
		return 0, err
	}

	return int(count), nil
}
