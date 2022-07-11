package service

import (
	"aed-api-server/internal/interfaces/entities"
	"aed-api-server/internal/interfaces/events"
)

type PointsScheduler interface {
	//ReloadSchedule 重新载入策略信息
	ReloadSchedule() error

	//GetPointStrategies 获取积分策略
	GetPointStrategies() ([]*entities.PointStrategy, error)

	//DealPointsEvent 处理积分事件
	DealPointsEvent(*events.PointsEvent) (*entities.DealPointsEventRst, error)

	// BuildPointsEventTypeGamePoints 构建游戏积分参数
	BuildPointsEventTypeGamePoints(userId, gameId int64, points int, description string) *events.PointsEvent
	//
	////BuildPointsEventTypeSignEarly 构造早起行为的积分事件
	//BuildPointsEventTypeSignEarly(userId int64, evtParams *PointsEventTypeSignEarlyParams) *PointsEvent

	BuildPointsEventTypeFriendsAddPoint(userId int64, friendsPointFlows []*entities.UserPointsRecord) *events.PointsEvent

	//BuildPointsEventTypeActivityGive 构造活动赠与的积分事件
	BuildPointsEventTypeActivityGive(userId int64, points int, description string) *events.PointsEvent

	BuildPointsEventTypeClockInDevice(userId int64, clockInId int64, job *entities.UserTask) *events.PointsEvent

	BuildPointsEventTypeMockedExam(userId, examId int64, score int) *events.PointsEvent

	BuildPointsEventTypeSignEarly(userId, signEarlyId int64, days int) *events.PointsEvent

	//BuildPointsEventWalk 构造步行的积分事件
	BuildPointsEventWalk(userId int64, todayWalk int, convertWalk int, convertedPoints int) *events.PointsEvent

	BuildPointsEventTypeReward(userId int64, jobId int64, points int, description string) *events.PointsEvent
}
