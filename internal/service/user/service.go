package user

import (
	"aed-api-server/internal/interfaces/entities"
	"aed-api-server/internal/interfaces/events"
	"aed-api-server/internal/interfaces/facility"
	service2 "aed-api-server/internal/interfaces/service"
	"aed-api-server/internal/pkg/db"
	page "aed-api-server/internal/pkg/query"
	"aed-api-server/internal/pkg/utils"
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
}

//go:inject-component
func NewService() *Service {
	return &Service{
		position: &positionService{},
	}
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

	_, err := db.Insert("user_event_record", event)

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

	err := db.SQL(fmt.Sprintf(`
		select * from (select * from user_event_record
        where user_id in %s
		and event_type = ?
        order by created_at desc) a
		group by user_id
    `, db.ParamPlaceHolder(len(userIds))), db.TupleOf(userIds, eventType)...).Find(&eventList)

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
