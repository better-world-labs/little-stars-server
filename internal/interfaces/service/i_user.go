package service

import (
	"aed-api-server/internal/interfaces/entities"
	"aed-api-server/internal/interfaces/events"
)

// UserService 用户服务
type UserService interface {

	// GetListUserByIDs 读取多个账号信息
	GetListUserByIDs(ids []int64) ([]*entities.SimpleUser, error)

	GetMapUserByIDs(ids []int64) (map[int64]*entities.SimpleUser, error)

	ListAllUsers() ([]*entities.UserDTO, error)

	StatUser() (entities.UserStat, error)

	// GetUserById 根据ID读取某个用户
	GetUserById(id int64) (*entities.SimpleUser, bool, error)

	// GetUserByOpenid 根据ID读取某个用户
	GetUserByOpenid(openid string) (*entities.SimpleUser, bool, error)

	// GetUserByPhone 根据手机号读取帐号
	GetUserByPhone(phone string) (*entities.SimpleUser, bool, error)

	RecordUserEvent(userId int64, eventType entities.UserEventType, eventParams ...interface{})

	GetLastUserEventByType(userId int64, eventType entities.UserEventType) (*events.UserEvent, error)

	BatchGetLastUserEventByType(userIds []int64, eventType entities.UserEventType) (map[int64]*events.UserEvent, error)

	GetUserOpenIdById(userId int64) (string, error)

	// Traverse 遍历所有用户
	Traverse(f func(dto entities.UserDTO))

	// TraverseSubscribeMessageTicketUser 遍历有发订阅消息权限的用户
	TraverseSubscribeMessageTicketUser(key entities.SubscribeMessageKey, f func(dto []*entities.UserDTO))

	DealUserReportEvents(userId int64, key string, params []interface{})

	// UpdatePosition 更新用户所在位置
	// @param position 经纬度坐标
	// @return 错误
	UpdatePosition(position *entities.Position) error

	ParseInfoFromJwtToken(token string) (*entities.User, error)
}
