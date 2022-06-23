package user

import (
	"aed-api-server/internal/interfaces/entities"
	"aed-api-server/internal/interfaces/events"
	"aed-api-server/internal/pkg/db"
	"aed-api-server/internal/pkg/domain/emitter"
	page "aed-api-server/internal/pkg/query"
	"errors"
	"fmt"
	log "github.com/sirupsen/logrus"
	"time"
)

type Service struct {
}

func NewService() *Service {
	return &Service{}
}

func (Service) ListAllUsers() ([]*entities.UserDTO, error) {
	accounts := make([]*entities.UserDTO, 0)

	err := db.Table("account").Find(&accounts)
	if err != nil {
		return nil, err
	}

	return accounts, nil
}

func (Service) GetListUserByIDs(ids []int64) ([]*entities.SimpleUser, error) {
	accounts := make([]*entities.SimpleUser, 0)
	err := db.Table("account").In("id", ids).Find(&accounts)
	if err != nil {
		return nil, err
	}
	return accounts, nil
}

func (Service) StatUser() (dto entities.UserStat, err error) {
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

func (Service) GetUserByPhone(phone string) (*entities.SimpleUser, bool, error) {
	var account entities.SimpleUser
	exists, err := db.Table("account").Where("mobile = ?", phone).Get(&account)
	if err != nil {
		return nil, exists, err
	}

	return &account, exists, nil
}
func (Service) GetUserById(id int64) (*entities.SimpleUser, bool, error) {
	var account entities.SimpleUser
	exists, err := db.Table("account").Where("id = ?", id).Get(&account)
	if err != nil {
		return nil, exists, err
	}

	return &account, exists, nil
}

func (Service) GetUserByOpenid(openid string) (*entities.SimpleUser, bool, error) {
	var account entities.SimpleUser
	exists, err := db.Table("account").Where("openid = ?", openid).Get(&account)
	if err != nil {
		return nil, exists, err
	}

	return &account, exists, nil
}

func (Service) RecordUserEvent(userId int64, eventType entities.UserEventType, eventParams ...interface{}) {
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

	err = emitter.Emit(&event)
	if err != nil {
		log.Error("emitter.Emit(event) err:", err)
	}
}

func (Service) BatchGetLastUserEventByType(userIds []int64, eventType entities.UserEventType) (map[int64]*events.UserEvent, error) {
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

func (Service) GetLastUserEventByType(userId int64, eventType entities.UserEventType) (*events.UserEvent, error) {
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

func (Service) GetUserOpenIdById(userId int64) (string, error) {
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

func (s Service) Traverse(f func(dto entities.UserDTO)) {
	users, err := s.ListAllUsers()
	if err != nil {
		log.Error("ListAllUsers error", err.Error())
		return
	}

	for _, u := range users {
		f(*u)
	}
}

func (Service) TraverseSubscribeMessageTicketUser(key entities.SubscribeMessageKey, f func(dto []*entities.UserDTO)) {
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
