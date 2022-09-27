package service

import (
	"aed-api-server/internal/interfaces/entities"
	"aed-api-server/internal/interfaces/events"
	"github.com/go-xorm/xorm"
)

type ActivityService interface {
	// UpdateByID 根据 ID 更新动态
	// @param activity 实体
	// @return error 错误
	UpdateByID(activity *entities.Activity) error

	// Create 创建动态
	// @param activity 实体
	// @return error 错误
	Create(activity *entities.Activity) error

	// CreateOrUpdateByUUID 创建或者更新
	// @param activity 实体
	// @return error 错误
	CreateOrUpdateByUUID(a *entities.Activity) error

	// CreateWithSession 创建动态
	// @param session xorm session
	// @param activity 实体
	// @return error 错误
	CreateWithSession(session *xorm.Session, activity *entities.Activity) error

	// GetOneByID 根据 ID 读取动态信息
	// @param id ID
	// @return error 错误
	GetOneByID(id int64) (*entities.Activity, error)

	// GetOneByUUID 根据 ID 读取动态信息
	// @param id ID
	// @return error 错误
	GetOneByUUID(uuid string) (*entities.Activity, bool, error)

	// GetManyByIDs 根据多个 ID 读取多个动态信息
	// @param ids 多个 ID
	// @return error 错误
	GetManyByIDs(ids []int64) ([]*entities.Activity, error)

	// ListByAID 读取某个求助信息的动态信息
	// @param helpInfoIDs 求助信息 ID
	// @param limit 条数
	// @return []*entities.Activity 数据切片
	// @return error 错误
	ListByAID(aid int64, limit int) ([]*entities.Activity, error)

	// ListByAIDs 根据多个求助信息 ID 读取某个求助信息的动态信息
	// @param helpInfoIDs 求助信息 ID
	// @return []*Activity 数据切片
	// @return error 错误
	ListByAIDs(aids []int64) (map[int64][]*entities.Activity, error)

	// 	ListCategorySorted(aid int64) ([]*Activity, error) 读取排序后的动态信息，填充了 Category字段，并按照Category优先级 - 日期 进行排序
	// @param helpInfoIDs 求助信息 ID
	// @return []*Activity 数据切片
	// @return error 错误
	ListCategorySorted(aid int64) ([]*entities.Activity, error)

	// ListMultiLatestCategorySorted 读取多个求助信息的排序后的动态信息，填充了 Category字段，并按照Category优先级 - 日期 进行排序，使用limit限制数据条数
	// @param aid 求助信息 ID
	// @param latest 最近 n 条数据
	// @return 数据 map
	// @return error 错误
	ListMultiLatestCategorySorted(aids []int64, latest int) (map[int64][]*entities.Activity, error)

	// ListMultiLatestCategorySortedAsync 异步读取多个求助信息的排序后的动态信息，填充了 Category字段，并按照Category优先级 - 日期 进行排序，使用limit限制数据条数
	// @param aid 求助信息 ID
	// @param latest 最近 n 条数据
	// @return 带有数据 map 的闭包
	// @return error 错误
	ListMultiLatestCategorySortedAsync(aids []int64, latest int) func() (map[int64][]*entities.Activity, error)

	// ListLatestCategorySorted 读取排序后的动态信息，填充了 Category字段，并按照Category优先级 - 日期 进行排序，使用limit限制数据条数
	// @param aid 求助信息 ID
	// @param latest 最近 n 条数据
	// @return []*Activity  数据切片
	// @return error 错误
	ListLatestCategorySorted(aid int64, latest int) ([]*entities.Activity, error)

	// ListLatestCategorySortedAsync ListLatestCategorySorted 的异步版，返回包含相关返回值的 channel
	// @param aid 求助信息 ID
	// @param latest 最近 n 条数据
	// @return 调用阻塞直到有结果或者错误返回
	ListLatestCategorySortedAsync(aid int64, latest int) func() ([]*entities.Activity, error)

	// ListByAidAndClass 根据 Aid 和 Class 查询 Activity 列表
	// @param aid 求助信息 ID
	// @param class 动态类型
	// @return 调用阻塞直到有结果或者错误返回
	ListByAidAndClass(aid int64, class string) ([]*entities.Activity, error)

	// GetLastUpdated 读取最新动态
	// @param aid 求助信息 ID
	// @param latest 最近 n 条数据
	// @return []*Activity  数据切片
	// @return error 错误
	GetLastUpdated(aid int64) (*entities.Activity, error)

	// SaveActivityAidCalled 保存已拨打急救电话动态
	SaveActivityAidCalled(event *events.AidCalledEvent) error

	// SaveActivityVolunteerNotified 保存已通知志愿者动态
	SaveActivityVolunteerNotified(event *events.VolunteerNotifiedEvent) error

	// SaveActivitySceneArrived 保存确认到达现场动态
	SaveActivitySceneArrived(event *events.SceneArrivedEvent) error

	// SaveActivityDeviceGot 保存已获取设备动态
	SaveActivityDeviceGot(event *events.DeviceGotEvent) ([]*entities.DealPointsEventRst, error)

	// SaveActivityNPCDeviceGot 保存已获取设备动态(NPC)
	SaveActivityNPCDeviceGot(event *events.DeviceGotEvent) error

	// SaveActivityGoingToScene 保存正在前往现场动态
	SaveActivityGoingToScene(event *events.GoingToSceneEvent) error

	// SaveActivitySceneCalled  保存已联系现场动态
	SaveActivitySceneCalled(event *events.SceneCalledEvent) error

	// SaveActivityGoingToGetDevice 保存正在前往取设备动态
	SaveActivityGoingToGetDevice(event *events.GoingToGetDeviceEvent) error

	// SaveActivitySceneReport  保存已上传场播报动态
	SaveActivitySceneReport(event *events.SceneReportEvent) ([]*entities.DealPointsEventRst, error)
}
