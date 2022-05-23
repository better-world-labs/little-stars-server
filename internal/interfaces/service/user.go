package service

import "aed-api-server/internal/interfaces/entities"

// UserService 用户服务
type UserService interface {

	// GetListUserByIDs 读取多个账号信息
	GetListUserByIDs(ids []int64) ([]*entities.SimpleUser, error)

	// GetUserById 根据ID读取某个用户
	GetUserById(id int64) (*entities.SimpleUser, bool, error)

	// GetUserByOpenid 根据ID读取某个用户
	GetUserByOpenid(openid string) (*entities.SimpleUser, bool, error)

	// GetUserByPhone 根据手机号读取帐号
	GetUserByPhone(phone string) (*entities.SimpleUser, bool, error)

	RecordUserEvent(userId int64, eventType string)
}
