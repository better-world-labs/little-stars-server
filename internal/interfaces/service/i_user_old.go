package service

import (
	"aed-api-server/internal/interfaces/entities"
)

type UserServiceOld interface {

	// WechatMiniProgramLogin 登录
	// @Param command 登陆命令体
	// @return *User & Token & SessionKey & 错误
	WechatMiniProgramLogin(command entities.LoginCommand) (*entities.User, string, string, error)

	WechatMiniProgramLoginV2(command entities.LoginCommandV2) (*entities.User, string, string, error)

	// WechatMiniProgramLoginSimple 免授权简单登录
	// @Param command 登陆命令体
	// @return *User & Token & SessionKey & 错误
	WechatMiniProgramLoginSimple(command entities.SimpleLoginCommand) (*entities.User, string, string, error)

	// ListAllPositions 读取所有用户位置
	// @param ids 用户 ID 集合
	// @return func (map[用户ID]经纬度坐标, 错误)
	ListAllPositions() ([]*entities.Position, error)

	ListUsers() (r map[int64]*entities.User, err error)

	UpdateUsersAvatar() (count int, err error)

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

	GetUserCharityCard(userId int64) (string, error)
}
