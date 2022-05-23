package user

import (
	"aed-api-server/internal/interfaces"
	"aed-api-server/internal/interfaces/entities"
	"aed-api-server/internal/interfaces/events"
	"aed-api-server/internal/pkg/cache"
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

	"gitlab.openviewtech.com/openview-pub/gopkg/log"
)

const CacheName = "user"

type service struct {
	client WechatClient
}

func NewService(c WechatClient) Service {
	return &service{
		client: c,
	}
}
func (s *service) WechatAppLogin(code string) (*User, string, error) {
	var info OAuthInfoDTO
	if err := s.client.GetAccessToken(code, &info.OAuthAccessTokenDTO); err != nil {
		return nil, "", err
	}
	if err := s.client.GetUserInfo(info.AccessToken, info.OpenID, &info); err != nil {
		return nil, "", err
	}

	account := User{
		Nickname: info.Nickname,
		Avatar:   info.HeadImgURL,
		Unionid:  info.UnionID,
		Openid:   info.OpenID,
	}
	err := s.createOrUpdateUser(&account)
	if err != nil {
		return nil, "", err
	}

	token, err := SignToken(account.ID)
	if err != nil {
		return nil, "", err
	}

	return &account, token, err
}

func (s *service) WechatMiniProgramLogin(command LoginCommand) (*User, string, string, error) {
	var info MiniProgramResponseDTO
	log.DefaultLogger().Infof("encryptPhone: %s, iv: %s", command.EncryptPhone, command.Iv)
	if err := s.client.MiniProgramCode2Session(command.Code, command.EncryptPhone, command.Iv, &info); err != nil {
		return nil, "", "", err
	}
	account := User{
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

	token, err := SignToken(account.ID)
	if err != nil {
		return nil, "", "", err
	}

	return &account, token, info.SessionKey, err
}

func (s *service) GetUserByOpenId(openid string) (*User, bool, error) {
	var u User
	exists, err := db.Table("account").Where("openid = ?", openid).Get(&u)
	return &u, exists, err
}

func (s *service) createOrUpdateUser(account *User) error {
	session := db.GetSession()
	defer session.Close()

	var current User
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

func (s *service) UpdatePosition(position *Position) error {
	session := db.GetSession()
	defer session.Close()

	var current Position
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

func (s *service) ListPositionByUserIDs(accountIDs []int64) (map[int64]*Position, error) {
	session := db.GetSession()
	defer session.Close()

	res := make(map[int64]*Position, 0)
	arr := make([]*Position, 0)
	cond := new(Position)
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

func (s *service) ListAllUser() ([]*User, error) {
	session := db.GetSession()
	defer session.Close()

	arr := make([]*User, 0)
	err := session.Table("account").Find(&arr)
	if err != nil {
		return nil, err
	}

	return arr, nil
}

func (s *service) ListUserByIDs(ids []int64) (r map[int64]*User, err error) {
	return s.ListUserByIDsFromDB(ids)
}

func (s *service) ListUserByIDsFromDB(ids []int64) (map[int64]*User, error) {
	session := db.GetSession()
	defer session.Close()

	res := make(map[int64]*User, 0)
	err := session.Table("account").In("id", ids).Find(&res)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (s *service) GetUserByID(id int64) (account *User, err error) {
	o, exists, err := cache.GetManager().Get(CacheName, getCacheKey(id))
	if err == nil && exists {
		return o.(*User), nil
	}

	o, err = s.getUserByIDFromDB(id)
	if err != nil {
		return nil, err
	}

	_ = cache.GetManager().Put(CacheName, getCacheKey(id), o)
	return o.(*User), err
}

func (s *service) getUserByIDFromDB(id int64) (account *User, err error) {
	session := db.GetSession()
	defer session.Close()

	account = &User{}
	ok, err := session.Table("account").Where("id = ?", id).Get(account)
	if err != nil {
		return nil, err
	}

	if !ok {
		return nil, nil
	}

	return account, nil
}

func (s *service) GetUserByIDAsync(id int64) func() (account *User, err error) {
	resChan := make(chan interface{}, 1)

	go func() {
		defer close(resChan)
		res, err := s.GetUserByID(id)
		if err == nil {
			resChan <- res
		} else {
			resChan <- err
		}
	}()

	return func() (account *User, err error) {
		res := <-resChan
		switch res.(type) {
		case *User:
			return res.(*User), nil
		case error:
			return nil, res.(error)
		default:
			return nil, errors.New("invalid type from chan")
		}
	}
}

func (s *service) GetPositionByUserID(accountID int64) (*Position, error) {
	session := db.GetSession()
	defer session.Close()

	position := &Position{}
	ok, err := session.Table("account_position").Where("account_id = ?", accountID).Get(position)
	if err != nil {
		return nil, err
	}

	if !ok {
		return nil, nil
	}

	return position, nil
}

func (s *service) UpdateMobile(accountID int64, mobile string) error {
	session := db.GetSession()
	defer session.Close()

	return db.WithTransaction(session, func() error {
		var current User
		exists, err := session.Table("account").Where("id = ?", accountID).Get(&current)
		if err != nil {
			return err
		}

		if !exists {
			return errors.New("account does not exists")
		}

		current.Mobile = mobile

		_, err = session.ID(current.ID).Update(&current)

		return cache.GetManager().Evict(CacheName, getCacheKey(accountID))
	})
}

func (s *service) ListAllUserAsync() func() ([]*User, error) {
	resChan := make(chan interface{}, 1)
	go func() {
		defer close(resChan)
		accounts, err := s.ListAllUser()
		if err != nil {
			resChan <- err
		} else {
			resChan <- accounts
		}
	}()

	return func() (accounts []*User, err error) {
		res := <-resChan
		switch res.(type) {
		case []*User:
			return res.([]*User), nil
		case error:
			return nil, res.(error)
		default:
			return nil, errors.New("invalid type from chan")
		}
	}
}

func (s *service) ListAllPositions() ([]*Position, error) {
	session := db.GetSession()
	defer session.Close()

	res := make([]*Position, 0)
	err := session.Table("account_position").Find(&res)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (s *service) ListAllPositionsAsync() func() ([]*Position, error) {
	resChan := make(chan interface{}, 1)
	go func() {
		defer close(resChan)
		accounts, err := s.ListAllPositions()
		if err != nil {
			resChan <- err
		} else {
			resChan <- accounts
		}
	}()

	return func() (accounts []*Position, err error) {
		res := <-resChan
		switch res.(type) {
		case []*Position:
			return res.([]*Position), nil
		case error:
			return nil, res.(error)
		default:
			return nil, errors.New("invalid type from chan")
		}
	}
}

func (s *service) ListPositionByUserIDsAsync(
	accountIDs []int64,
) func() (map[int64]*Position, error) {
	resChan := make(chan interface{}, 1)
	go func() {
		defer close(resChan)
		accounts, err := s.ListPositionByUserIDs(accountIDs)
		if err != nil {
			resChan <- err
		} else {
			resChan <- accounts
		}
	}()

	return func() (accounts map[int64]*Position, err error) {
		res := <-resChan
		switch res.(type) {
		case map[int64]*Position:
			return res.(map[int64]*Position), nil
		case error:
			return nil, res.(error)
		default:
			return nil, errors.New("invalid type from chan")
		}
	}
}

func (s *service) UpdateUserInfo(account *User) error {
	session := db.GetSession()
	defer session.Close()

	return db.WithTransaction(session, func() error {
		var current User
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

func (s *service) generateUid() string {
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

func (s *service) GenerateUidForExistsUser() error {
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

func (s *service) DealUserEvents(userId int64, userEventType string) {
	utils.Go(func() {
		interfaces.S.User.RecordUserEvent(userId, "report:"+userEventType)
	})

	if userEventType == "enter-aed-map" {
		utils.Go(func() {
			dealUserEnterAEDMap(userId)
		})
	}

	if userEventType == "read-news" {
		utils.Go(func() {
			dealUserReadNews(userId)
		})
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

func (s *service) GetUserCharityCard(userId int64) (string, error) {
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

	user := rst[0].(*User)
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
		fmt.Sprintf(`https://%s/share/cert?source=charity-card&sharer=%v`, config.Host, userId),
		fmt.Sprintf("user-charity-card-%v.jpeg", userId),
	)
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

func (s *service) UpdateUid(user *User) error {
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

func (s *service) WechatMiniProgramLoginSimple(command SimpleLoginCommand) (*User, string, string, error) {
	session, err := s.client.CodeToSession(command.Code)
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

	token, err := SignToken(u.ID)
	if err != nil {
		return nil, "", "", err
	}

	return u, token, session.SessionKey, nil
}

func (s *service) CountUser() (int, error) {
	count, err := db.Table("account").Count()
	if err != nil {
		return 0, err
	}

	return int(count), nil
}
