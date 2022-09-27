package user

import (
	"aed-api-server/internal/interfaces"
	"aed-api-server/internal/interfaces/entities"
	"aed-api-server/internal/interfaces/events"
	"aed-api-server/internal/interfaces/service"
	"aed-api-server/internal/pkg/cache"
	"aed-api-server/internal/pkg/db"
	"aed-api-server/internal/pkg/domain/emitter"
	"aed-api-server/internal/pkg/utils"
	"aed-api-server/internal/pkg/wx_crypto"
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
	Wechat        service.IWechat              `inject:"-"`
	TreasureChest service.TreasureChestService `inject:"-"`
	Oss           service.OssService           `inject:"-"`
}

//go:inject-component
func NewUserServiceOld() service.UserServiceOld {
	return &userServiceOld{}
}

func (s *userServiceOld) uploadAvatar(userId int64, originAvatar string) (string, error) {
	resp, err := http.Get(originAvatar)
	if err != nil {
		return "", err
	}
	if resp.StatusCode != 200 {
		return "", errors.New(resp.Status)
	}

	url, err := s.Oss.Upload(fmt.Sprintf("avatar/%d.png", userId), resp.Body)
	if err != nil {
		return "", err
	}

	return url, nil
}

func (s *userServiceOld) WechatMiniProgramLogin(command entities.LoginCommand) (*entities.User, string, string, error) {
	var info entities.WechatMiniProgramRes
	log.Infof("encryptPhone: %s, iv: %s", command.EncryptPhone, command.Iv)
	if err := s.Wechat.MiniProgramCode2Session(command.Code, command.EncryptPhone, command.Iv, &info); err != nil {
		return nil, "", "", err
	}
	account := entities.User{
		Nickname:   command.Nickname,
		Avatar:     command.Avatar,
		Mobile:     info.DecryptedPhone,
		Unionid:    info.UnionID,
		Openid:     info.OpenID,
		SessionKey: info.SessionKey,
	}

	err := s.createOrUpdateUser(&account)
	if err != nil {
		return nil, "", "", err
	}

	err = s.updateAvatar(&account)
	if err != nil {
		return nil, "", "", err
	}

	token, err := SignToken(account.ID)
	if err != nil {
		return nil, "", "", err
	}

	return &account, token, info.SessionKey, err
}

func (s *userServiceOld) WechatMiniProgramLoginV2(command entities.LoginCommandV2) (*entities.User, string, string, error) {
	session, err := s.Wechat.CodeToSession(command.Code)
	if err != nil {
		return nil, "", "", err
	}

	account := entities.User{
		Nickname:   command.Nickname,
		Avatar:     command.Avatar,
		Unionid:    session.UnionId,
		Openid:     session.Openid,
		SessionKey: session.SessionKey,
	}

	err = s.createOrUpdateUser(&account)
	if err != nil {
		return nil, "", "", err
	}

	err = s.updateAvatar(&account)
	if err != nil {
		return nil, "", "", err
	}

	token, err := SignToken(account.ID)
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

func (s *userServiceOld) updateAvatar(account *entities.User) error {
	avatar, err := s.uploadAvatar(account.ID, account.Avatar)
	if err != nil {
		return err
	}

	account.Avatar = avatar
	return s.createOrUpdateUser(account)
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
		if err != nil {
			return err
		}

		if account.Available() && !current.Available() {
			err = emitter.Emit(&events.FirstLoginEvent{
				UserId:  account.ID,
				Openid:  account.Openid,
				LoginAt: account.Created,
			})
			if err != nil {
				log.Errorf("FirstLoginEvent emit error: %v", err)
			}
		}

		return cache.GetManager().Evict(CacheName, getCacheKey(account.ID))
	} else {
		account.Created = time.Now()
		account.Uid = s.generateUid()
		_, err := session.Table("account").Insert(account)
		if err != nil {
			return err
		}

		return nil
	}

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

func (s *userServiceOld) UpdateUsersAvatar() (count int, err error) {
	users, err := s.ListUsers()
	if err != nil {
		return 0, err
	}

	for _, u := range users {
		if !strings.HasPrefix(u.Avatar, "https://openview-oss") {
			err := s.updateAvatar(u)
			if err != nil {
				log.Infof("updateAvatoar for user %d error: %v\n", u.ID, err)
				continue
			}

			count++
		}
	}

	return
}
func (s *userServiceOld) ListUsers() (r map[int64]*entities.User, err error) {
	session := db.GetSession()
	defer session.Close()

	res := make(map[int64]*entities.User, 0)
	err = session.Table("account").Find(&res)
	if err != nil {
		return nil, err
	}

	return res, nil
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
	weSession, err := s.Wechat.CodeToSession(command.Code)
	if err != nil {
		return err
	}

	var phone wx_crypto.WxUserPhone
	_, err = wx_crypto.Decrypt(command.EncryptPhone, command.Iv, weSession.SessionKey, &phone)
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
		current.SessionKey = weSession.SessionKey

		_, err = session.Table("account").ID(current.ID).Update(&current)
		if err != nil {
			return err
		}

		return cache.GetManager().Evict(CacheName, getCacheKey(accountID))
	})
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

		_, err = session.Table("account").ID(account.ID).Update(account)
		if err != nil {
			return err
		}

		return cache.GetManager().Evict(CacheName, getCacheKey(account.ID))
	})
}

func getCacheKey(id int64) string {
	return fmt.Sprintf("id:%d", id)
}

func (s *userServiceOld) generateUid() string {
	u := uuid.New()
	formatted := strconv.FormatUint(uint64(u.ID()), 10)
	l := len(formatted)
	if l < 10 {
		for i := 0; i < 10-i; i++ {
			formatted += "0"
		}
	}
	return formatted[:10]
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

func (s *userServiceOld) UpdateUid(user *entities.User) error {
	uid := s.generateUid()
	user.Uid = uid
	_, err := db.Table("account").ID(user.ID).Update(user)
	return err
}

func (s *userServiceOld) WechatMiniProgramLoginSimple(command entities.SimpleLoginCommand) (*entities.User, string, string, error) {
	session, err := s.Wechat.CodeToSession(command.Code)
	if err != nil {
		return nil, "", "", err
	}

	u, exists, err := s.GetUserByOpenId(session.Openid)
	if err != nil {
		return nil, "", "", err
	}

	if !exists {
		u = &entities.User{
			Openid:     session.Openid,
			Unionid:    session.UnionId,
			SessionKey: session.SessionKey,
		}
	} else {
		u.SessionKey = session.SessionKey
	}

	err = s.createOrUpdateUser(u)
	if err != nil {
		return nil, "", "", err
	}

	token, err := SignToken(u.ID)
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
