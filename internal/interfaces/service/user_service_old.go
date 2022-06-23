package service

import (
	"aed-api-server/internal/interfaces/entities"
)

type UserServiceOld interface {

	// WechatAppLogin 登录
	// @Param code 授权码
	// @return *User & Token &错误
	WechatAppLogin(code string) (*entities.User, string, error)

	// WechatMiniProgramLogin 登录
	// @Param command 登陆命令体
	// @return *User & Token & SessionKey & 错误
	WechatMiniProgramLogin(command entities.LoginCommand) (*entities.User, string, string, error)

	WechatMiniProgramLoginV2(command entities.LoginCommandV2) (*entities.User, string, string, error)

	// WechatMiniProgramLoginSimple 免授权简单登录
	// @Param command 登陆命令体
	// @return *User & Token & SessionKey & 错误
	WechatMiniProgramLoginSimple(command entities.SimpleLoginCommand) (*entities.User, string, string, error)

	// UpdatePosition 更新用户所在位置
	// @param position 经纬度坐标
	// @return 错误
	UpdatePosition(position *entities.Position) error

	// ListAllUser 读取所有用户位置
	// @return (map[用户ID]经纬度坐标, 错误)
	ListAllUser() ([]*entities.User, error)

	// ListAllUserAsync 异步读取所有用户位置
	// @return func (map[用户ID]经纬度坐标, 错误)
	ListAllUserAsync() func() ([]*entities.User, error)

	// ListPositionByUserIDs 读取多个用户的所在位置
	// @param accountIDs 用户 ID 集合
	// @return (map[用户ID]经纬度坐标, 错误)
	ListPositionByUserIDs(accountIDs []int64) (map[int64]*entities.Position, error)

	// ListPositionByUserIDsAsync 读取多个账号信息
	// @param ids 用户 ID 集合
	// @return func (map[用户ID]经纬度坐标, 错误)
	ListPositionByUserIDsAsync(accountIDs []int64) func() (map[int64]*entities.Position, error)

	// ListAllPositions 读取所有用户位置
	// @param ids 用户 ID 集合
	// @return func (map[用户ID]经纬度坐标, 错误)
	ListAllPositions() ([]*entities.Position, error)

	// ListAllPositionsAsync 异步读取所有用户位置
	// @param ids 用户 ID 集合
	// @return func (map[用户ID]经纬度坐标, 错误)
	ListAllPositionsAsync() func() ([]*entities.Position, error)
	// ListUserByIDs 读取多个账号信息
	// @param ids 用户 ID 集合
	// @return (map[用户ID]经纬度坐标, 错误)
	ListUserByIDs(ids []int64) (map[int64]*entities.User, error)

	// GetUserByID 读取某个账号信息
	// @param id 账号ID
	// @return (map[用户ID]经纬度坐标, 错误)
	GetUserByID(id int64) (*entities.User, error)

	// CountUser 读取注册用户数
	// @return (用户数, 错误)
	CountUser() (int, error)

	// GetPositionByUserID 读取某个账号信息
	// @param accountID 账号ID
	// @return (map[账号ID]经纬度坐标, 错误)
	GetPositionByUserID(accountID int64) (*entities.Position, error)

	// UpdateMobile update mobile info
	// @param accountID account iD
	// @return error
	UpdateMobile(accountID int64, command entities.MobileCommand) error

	// UpdateUserInfo update account info
	// @param account entity
	// @return error
	UpdateUserInfo(account *entities.User) error

	// GenerateUidForExistsUser 为没有UID的用户生成UID
	GenerateUidForExistsUser() error

	DealUserEvents(userId int64, key string, params []interface{})

	GetUserCharityCard(userId int64) (string, error)
}
