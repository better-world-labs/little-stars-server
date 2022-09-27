package service

import (
	"aed-api-server/internal/interfaces/entities"
	page "aed-api-server/internal/pkg/query"
	"time"
)

type PointsService interface {

	// GetUserPointsRecords 读取积分流水
	GetUserPointsRecords(userId int64) ([]*entities.UserPointsRecord, error)

	// GetUserTotalPoints 获取用户积分余额
	GetUserTotalPoints(userId int64) (int, error)

	// GetUserIncomePoints 获取用户累收入的计积分
	GetUserIncomePoints(userId int64) (int, error)

	// GetUserPointsEventTimes 读取用户对某个行为的积分次数
	GetUserPointsEventTimes(userId int64, eventType entities.PointsEventType) (int64, error)

	// GetUserDonatedPoints 读取用户捐献的积分总数
	GetUserDonatedPoints(userId int64) (int, error)

	// ReceivePoints 领取积分
	ReceivePoints(userId int64, pointId int64) error

	// AddPoint 增加积分, points < 0 则减少积分, 会返回金额不足的error
	// 临时接口,用于交易,由于目前没有账户的概念,募捐增加的现在募捐出维护,此处只扣减人的
	AddPoint(userId int64, points int, description string, eventType entities.PointsEventType) error

	DealPoint(userId int64, points int, description string, eventType entities.PointsEventType, autoReceive bool) error

	// Pay 积分支付
	Pay(userId int64, points int, description string) error

	// GetUnReceivePoints 获取用户未领取的积分流水
	GetUnReceivePoints(userId int64) ([]*entities.UserPointsFlow, error)

	// GetUnReceivePointsCount 获取未领取积分记录个数
	GetUnReceivePointsCount(userId int64) (int, error)

	// GetUsersPeriodIncomePoints 读取多个用户在时间范围内的收入积分总数
	GetUsersPeriodIncomePoints(userIds []int64, beginTime time.Time, endTime time.Time) (map[int64]int, error)

	// GetUserPeriodIncomePointsRecords 读取用户在时间范围内的积分收入流水
	GetUserPeriodIncomePointsRecords(userId int64, begin, end time.Time) ([]*entities.UserPointsRecord, error)

	GetAllPointsExpiringUserIds(duration time.Duration) ([]int64, error)

	StatExpiringPoints(userId int64) (points int, minExpiredAt time.Time, err error)

	IsTodayHasPointFlowOfType(userId int64, pointsEventType entities.PointsEventType) (bool, error)

	IsTodayHasPointFlowOfTypeBatched(userIds []int64, pointsEventType entities.PointsEventType) ([]int64, error)

	Detail(accountID int64) ([]*entities.Point, error)

	TotalPoints(accountID int64) (float64, error)

	PageAwardPointFLows(query page.Query, filter entities.AwardFlowQueryCommand) (page.Result[*entities.AwardPointFlow], error)
}
