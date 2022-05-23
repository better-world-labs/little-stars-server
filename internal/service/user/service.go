package user

import (
	"aed-api-server/internal/interfaces"
	"aed-api-server/internal/interfaces/entities"
	"aed-api-server/internal/interfaces/service"
	"aed-api-server/internal/pkg/db"
	"gitlab.openviewtech.com/openview-pub/gopkg/log"
	"time"
)

type Service struct {
	service.UserService
}

func init() {
	if interfaces.S.User == nil {
		interfaces.S.User = &Service{}
	}
}

func (Service) GetListUserByIDs(ids []int64) ([]*entities.SimpleUser, error) {
	accounts := make([]*entities.SimpleUser, 0)
	err := db.Table("account").In("id", ids).Find(&accounts)
	if err != nil {
		return nil, err
	}
	return accounts, nil
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

type UserEvent struct {
	Id        int64
	UserId    int64
	EventType string
	CreatedAt time.Time
}

func (Service) RecordUserEvent(userId int64, eventType string) {
	event := UserEvent{
		UserId:    userId,
		EventType: eventType,
		CreatedAt: time.Now(),
	}
	_, err := db.Insert("user_event_record", event)

	if err != nil {
		log.Error("RecordUserEvent err:", err)
	}
}
