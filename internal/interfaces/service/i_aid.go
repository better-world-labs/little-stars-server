package service

import (
	"aed-api-server/internal/interfaces/entities"
	"aed-api-server/internal/pkg/location"
	page "aed-api-server/internal/pkg/query"
)

type AidService interface {
	// PublishHelpInfo 发布求助信息
	// @param accountID 帐号ID
	// @param dto DTO
	// @return int64 ID
	// @return error 错误
	PublishHelpInfo(accountID int64, dto *entities.PublishDTO) (int64, []*entities.DealPointsEventRst, error)

	GetHelpInfoByID(id int64) (*entities.HelpInfo, bool, error)

	PublishHelpInfoExercise(userId int64, dto *entities.PublishDTO) (id int64, npc *entities.SimpleUser, err error)

	// ListHelpInfosPaged 分页读取求助信息列表
	// @param pageQuery 分页参数
	// @param position 当前经纬度坐标
	// @param condition 查询条件
	// @return page.Result 分页数据结果
	// @return error 错误
	ListHelpInfosPaged(pageQuery *page.Query, position *location.Coordinate, condition *entities.HelpInfo) (*page.Result[*entities.HelpInfoComposedDTO], error)

	// ListOneHoursInfos 查询一小时内的求助信息
	ListOneHoursInfos() ([]*entities.HelpInfo, error)

	// ListHelpInfosParticipatedPaged 分页读取我参与的求助信息
	// @param pageQuery 分页参数
	// @return page.Result 分页数据结果
	// @return error 错误
	ListHelpInfosParticipatedPaged(pageQuery *page.Query, userID int64) (*page.Result[*entities.HelpInfoComposedDTO], error)

	CountHelpInfosAboutMe(userId int64) (int64, error)

	// ListHelpInfosInner24h 查询24h内的求助信息
	// @return []*HelpInfo 分页数据结果
	// @return error 错误
	ListHelpInfosInner24h() ([]*entities.HelpInfo, error)

	// ActionArrived 到达现场行为触发
	// @param accountID 帐号 ID
	// @param aidID 求助信息 ID
	// @return 错误
	ActionArrived(accountID int64, aidID int64, coordinate *location.Coordinate) ([]*entities.DealPointsEventRst, error)

	ActionNPCArrived(aidId int64) error

	// ActionCalled 触发电话联系现场行为
	// @param accountID 帐号 ID
	// @param aidID 求助信息 ID
	// @return 错误
	ActionCalled(accountID int64, aidID int64) error

	// Action120Called 运营人员触发已经拨打120事件
	// @param aidID 求助信息 ID
	// @return 错误
	Action120Called(aidID int64) error

	// ActionGoingToScene 触发正在前往现场行为
	// @param accountID 帐号 ID
	// @param aidID 求助信息 ID
	// @return 错误
	ActionGoingToScene(accountID int64, aidID int64) error

	// GetHelpInfoComposedByID 获取HelpInfo组合对象
	// @param id 帐号 求助信息 ID
	// @param position 位置
	// @return 错误
	GetHelpInfoComposedByID(id int64, position *location.Coordinate) (*entities.HelpInfoComposedDTO, bool, error)
}
