package user

import (
	"aed-api-server/internal/interfaces/entities"
	"aed-api-server/internal/interfaces/events"
	"aed-api-server/internal/interfaces/facility"
	service2 "aed-api-server/internal/interfaces/service"
	"aed-api-server/internal/pkg/db"
	page "aed-api-server/internal/pkg/query"
	"aed-api-server/internal/pkg/utils"
	"aed-api-server/internal/pkg/wx_crypto"
	"errors"
	"fmt"
	log "github.com/sirupsen/logrus"
	"strings"
	"time"
)

type Service struct {
	facility.Sender `inject:"sender"`
	position
	TreasureChest service2.TreasureChestService `inject:"-"`
	Donation      service2.DonationService      `inject:"-"`
	Aid           service2.AidService           `inject:"-"`
	Feed          service2.IFeed                `inject:"-"`
	Wechat        service2.IWechat              `inject:"-"`
	User          service2.UserServiceOld       `inject:"-"`

	Encrypt *CryptKeyCache `inject:"-"`
}

//go:inject-component
func NewService() *Service {
	return &Service{
		position: &positionService{},
	}
}

func (*Service) PageUsers(query page.Query, keyword string) (page.Result[*entities.SimpleUser], error) {
	where := "nickname like concat(?,'%') or mobile like concat(?,'%') or uid like concat(?,'%')"
	count, err := db.Table("account").Where(where, keyword, keyword, keyword).Count()
	if err != nil {
		return page.Result[*entities.SimpleUser]{}, err
	}

	var res []*entities.SimpleUser
	err = db.Table("account").Where(where, keyword, keyword, keyword).Limit(query.Size, (query.Page-1)*query.Size).Find(&res)
	if err != nil {
		return page.Result[*entities.SimpleUser]{}, err
	}

	return page.NewResult(res, int(count)), nil
}

func (*Service) ListAllUsers() ([]*entities.UserDTO, error) {
	accounts := make([]*entities.UserDTO, 0)

	err := db.Table("account").Find(&accounts)
	if err != nil {
		return nil, err
	}

	return accounts, nil
}

func (s *Service) GetMapUserByIDs(ids []int64) (map[int64]*entities.SimpleUser, error) {
	users, err := s.GetListUserByIDs(ids)
	if err != nil {
		return nil, err
	}

	m := make(map[int64]*entities.SimpleUser)
	for _, u := range users {
		m[u.ID] = u
	}

	return m, nil
}

func (*Service) GetListUserByIDs(ids []int64) ([]*entities.SimpleUser, error) {
	accounts := make([]*entities.SimpleUser, 0)
	err := db.Table("account").In("id", ids).Find(&accounts)
	if err != nil {
		return nil, err
	}
	return accounts, nil
}

func (*Service) StatUser() (dto entities.UserStat, err error) {
	exited, err := db.SQL(`
		select
			count(*) as total_count,
			count(if(length(mobile)>10,1,null)) as mobile_count
		from account
	`).Get(&dto)

	if err != nil || !exited {
		return entities.UserStat{}, err
	}
	return dto, nil
}

func (*Service) GetUserByPhone(phone string) (*entities.SimpleUser, bool, error) {
	var account entities.SimpleUser
	exists, err := db.Table("account").Where("mobile = ?", phone).Get(&account)
	if err != nil {
		return nil, exists, err
	}

	return &account, exists, nil
}
func (*Service) GetUserById(id int64) (*entities.SimpleUser, bool, error) {
	var account entities.SimpleUser
	exists, err := db.Table("account").Where("id = ?", id).Get(&account)
	if err != nil {
		return nil, exists, err
	}

	return &account, exists, nil
}

func (*Service) GetUserInfoByOpenid(openid string) (*entities.User, error) {
	var account entities.User
	exists, err := db.Table("account").Where("openid = ?", openid).Get(&account)
	if err != nil {
		return nil, err
	}

	if !exists {
		return nil, nil
	}

	return &account, nil
}

func (*Service) GetUserByOpenid(openid string) (*entities.SimpleUser, bool, error) {
	var account entities.SimpleUser
	exists, err := db.Table("account").Where("openid = ?", openid).Get(&account)
	if err != nil {
		return nil, exists, err
	}

	return &account, exists, nil
}

func (s *Service) RecordUserEvent(userId int64, eventType entities.UserEventType, eventParams ...interface{}) {
	event := events.UserEvent{
		UserId:      userId,
		EventType:   eventType,
		EventParams: eventParams,
		CreatedAt:   time.Now(),
	}

	_, err := db.Insert("user_event_record", &event)

	if err != nil {
		log.Error("RecordUserEvent err:", err)
	}

	err = s.Send(&event)
	if err != nil {
		log.Error("emitter.Emit(event) err:", err)
	}
}

func (*Service) BatchGetLastUserEventByType(userIds []int64, eventType entities.UserEventType) (map[int64]*events.UserEvent, error) {
	var eventList []*events.UserEvent
	sql := `
		select * from (
			select * 
			from user_event_record
        	where user_id in (?)
				and event_type = ?
        	order by created_at desc
		) a
		group by user_id
    `
	sql, args := db.MustIn(sql, userIds, eventType)

	err := db.SQL(sql, args...).Find(&eventList)

	r := make(map[int64]*events.UserEvent)
	for _, e := range eventList {
		r[e.UserId] = e
	}

	return r, err
}

func (*Service) GetLastUserEventByType(userId int64, eventType entities.UserEventType) (*events.UserEvent, error) {
	var event events.UserEvent
	existed, err := db.SQL(`
		select *
		from user_event_record
		where
			user_id = ?
			and event_type = ?
		order by created_at desc
		limit 1
	`, userId, eventType).Get(&event)
	if err != nil {
		return nil, err
	}
	if !existed {
		return nil, nil
	}
	return &event, nil
}

func (*Service) GetUserInfo(userId int64) (*entities.User, error) {
	var user entities.User
	existed, err := db.Table("account").Where("id = ?", userId).Get(&user)
	if err != nil {
		return nil, err
	}
	if !existed {
		return nil, errors.New("user not existed")
	}
	return &user, nil
}

func (*Service) GetUserOpenIdById(userId int64) (string, error) {
	var user entities.UserDTO
	existed, err := db.Table("account").Where("id = ?", userId).Select("openid").Get(&user)
	if err != nil {
		return "", err
	}
	if !existed {
		return "", errors.New("user not existed")
	}
	return user.Openid, nil
}

func (s *Service) Traverse(f func(dto entities.UserDTO)) {
	users, err := s.ListAllUsers()
	if err != nil {
		log.Error("ListAllUsers error", err.Error())
		return
	}

	for _, u := range users {
		f(*u)
	}
}

func (*Service) TraverseSubscribeMessageTicketUser(key entities.SubscribeMessageKey, f func(dto []*entities.UserDTO)) {
	for i := 1; ; i++ {
		p := page.Query{Page: i, Size: 5000}
		accounts := make([]*entities.UserDTO, 0)
		limit, offset := p.GetLimit()
		err := db.SQL(fmt.Sprintf(`
		select
			b.id, b.openid
		from user_subscribe_message as a
		inner join account as b on b.id = a.user_id
		where
			%s > 0
		limit %d, %d
	`, key, limit, offset)).Find(&accounts)
		if err != nil {
			log.Error("TraverseSubscribeMessageTicketUser error", err.Error())
			return
		}

		if len(accounts) == 0 {
			break
		}

		f(accounts)
	}
}

func (s *Service) DealUserReportEvents(userId int64, key string, params []interface{}) {
	utils.Go(func() {
		s.RecordUserEvent(userId, entities.GetUserEventTypeOfReport(key), params...)
	})

	utils.Go(func() {
		switch key {
		case entities.Report_enterAedMap:
			dealUserEnterAEDMap(userId)

		case entities.Report_viewFeedPage:
			dealUserEnterCommunity(userId)

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

func (s *Service) UpdatePosition(position *entities.Position) error {
	utils.Go(func() {
		err := s.RecordPosition(position.AccountID, position.Latitude, position.Longitude)
		if err != nil {
			log.Error("s.RecordPosition error", err)
		}
	})
	return nil
}

func (*Service) ParseInfoFromJwtToken(token string) (*entities.User, error) {
	split := strings.Split(token, " ")

	if len(split) != 2 || split[0] != "Bearer" || split[1] == "" {
		return nil, errors.New("invalid token")
	}

	claims, err := ParseToken(split[1])
	if err != nil {
		return nil, errors.New("invalid token")
	}

	var acc entities.User
	exists, err := db.Table("account").Where("id=?", claims.ID).Get(&acc)
	if err != nil {
		log.Error("account query error:", err)
		return nil, errors.New("invalid token")
	}

	if !exists {
		log.Info("user not exited", claims.ID)
		return nil, errors.New("invalid token")
	}
	return &acc, nil
}

func (s *Service) GetWalks(req *entities.WechatDataDecryptReq) (*entities.WechatWalkData, error) {
	session, err := s.Code2Session(req.Code)
	if err != nil {
		return nil, err
	}

	var data entities.WechatWalkData
	b, err := wx_crypto.Decrypt(req.EncryptedData, req.Iv, session.SessionKey, &data)

	log.Infof("GetWalks: walkData=%s", string(b))
	return &data, err
}

func (s *Service) Code2Session(code string) (*entities.WechatCode2SessionRes, error) {
	session, err := s.Wechat.CodeToSession(code)
	if err != nil {
		return nil, err
	}

	user, err := s.GetUserInfoByOpenid(session.Openid)
	if err != nil {
		return nil, err
	}

	if user != nil {
		user.SessionKey = session.SessionKey
		if err := s.User.UpdateUserInfo(user); err != nil {
			return nil, err
		}
	}

	return session, nil
}
func (s *Service) GetUserEncryptKey(id int64, version int) (*entities.WechatEncryptKey, error) {
	key, err := s.Encrypt.GetKey(id, version)
	if err != nil {
		return nil, err
	}

	if key != nil {
		return key, nil
	}

	user, err := s.GetUserInfo(id)
	if err != nil {
		return nil, err
	}

	if user == nil {
		return nil, errors.New("user not found")
	}

	encryptKey, err := s.Wechat.GetUserEncryptKey(user.Openid, user.SessionKey)
	if err != nil {
		return nil, err
	}

	key = entities.GetEncryptKey(encryptKey, version)
	if key == nil {
		return nil, errors.New("encrypt key not found")
	}

	err = s.Encrypt.PutKeys(user.ID, encryptKey)
	return key, err
}

func (s *Service) GetUserAboutStat(id int64) (*entities.UserAboutStat, error) {
	all, err := utils.PromiseAll(func() (interface{}, error) {
		return s.Donation.StatDonationByUserId(id)
	}, func() (interface{}, error) {
		return s.Aid.CountHelpInfosAboutMe(id)
	}, func() (interface{}, error) {
		return s.Feed.GetMyFeedsCount(id)
	})

	if err != nil {
		return nil, err
	}

	return &entities.UserAboutStat{
		UserAboutDonationsCount: all[0].(entities.DonationStat).DonationProjectCount,
		UserAboutSosCount:       int(all[1].(int64)),
		UserAboutFeedsCount:     all[2].(int),
	}, nil
}

func (s *Service) GetUserByUid(uid string) (*entities.SimpleUser, error) {
	if uid == "" {
		return nil, nil
	}
	var user entities.SimpleUser
	has, err := db.Table("account").Where("uid = ?", uid).Get(&user)
	if err != nil {
		return nil, err
	}
	if has {
		return &user, nil
	}
	return nil, nil
}
