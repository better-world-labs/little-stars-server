package user

// Service 用户服务
type Service interface {

	// WechatAppLogin 登录
	// @Param code 授权码
	// @return *User & Token &错误
	WechatAppLogin(code string) (*User, string, error)

	// WechatMiniProgramLogin 登录
	// @Param command 登陆命令体
	// @return *User & Token & SessionKey & 错误
	WechatMiniProgramLogin(command LoginCommand) (*User, string, string, error)

	// WechatMiniProgramLoginSimple 免授权简单登录
	// @Param command 登陆命令体
	// @return *User & Token & SessionKey & 错误
	WechatMiniProgramLoginSimple(command SimpleLoginCommand) (*User, string, string, error)

	// UpdatePosition 更新用户所在位置
	// @param position 经纬度坐标
	// @return 错误
	UpdatePosition(position *Position) error

	// ListAllUser 读取所有用户位置
	// @return (map[用户ID]经纬度坐标, 错误)
	ListAllUser() ([]*User, error)

	// ListAllUserAsync 异步读取所有用户位置
	// @return func (map[用户ID]经纬度坐标, 错误)
	ListAllUserAsync() func() ([]*User, error)

	// ListPositionByUserIDs 读取多个用户的所在位置
	// @param accountIDs 用户 ID 集合
	// @return (map[用户ID]经纬度坐标, 错误)
	ListPositionByUserIDs(accountIDs []int64) (map[int64]*Position, error)

	// ListPositionByUserIDsAsync 读取多个账号信息
	// @param ids 用户 ID 集合
	// @return func (map[用户ID]经纬度坐标, 错误)
	ListPositionByUserIDsAsync(accountIDs []int64) func() (map[int64]*Position, error)

	// ListAllPositions 读取所有用户位置
	// @param ids 用户 ID 集合
	// @return func (map[用户ID]经纬度坐标, 错误)
	ListAllPositions() ([]*Position, error)

	// ListAllPositionsAsync 异步读取所有用户位置
	// @param ids 用户 ID 集合
	// @return func (map[用户ID]经纬度坐标, 错误)
	ListAllPositionsAsync() func() ([]*Position, error)
	// ListUserByIDs 读取多个账号信息
	// @param ids 用户 ID 集合
	// @return (map[用户ID]经纬度坐标, 错误)
	ListUserByIDs(ids []int64) (map[int64]*User, error)

	// GetUserByID 读取某个账号信息
	// @param id 账号ID
	// @return (map[用户ID]经纬度坐标, 错误)
	GetUserByID(id int64) (*User, error)

	// CountUser 读取注册用户数
	// @return (用户数, 错误)
	CountUser() (int, error)

	// GetPositionByUserID 读取某个账号信息
	// @param accountID 账号ID
	// @return (map[账号ID]经纬度坐标, 错误)
	GetPositionByUserID(accountID int64) (*Position, error)

	// UpdateMobile update mobile info
	// @param accountID account iD
	// @return error
	UpdateMobile(accountID int64, mobile string) error

	// UpdateUserInfo update account info
	// @param account entity
	// @return error
	UpdateUserInfo(account *User) error

	// GenerateUidForExistsUser 为没有UID的用户生成UID
	GenerateUidForExistsUser() error

	DealUserEvents(userId int64, userEventType string)

	GetUserCharityCard(userId int64) (string, error)
}
