package friends

import (
	"aed-api-server/internal/interfaces"
	"aed-api-server/internal/interfaces/entities"
	"aed-api-server/internal/interfaces/service"
	"aed-api-server/internal/pkg/utils"
	"errors"
	"time"
)

type Service struct {
}

//go:inject-component
func NewService() service.FriendsService {
	return &Service{}
}

func getYesterdayPoints(userIds []int64) (map[int64]int, error) {
	now := time.Now()
	h := time.Duration(now.Hour()) * time.Hour
	m := time.Duration(now.Minute()) * time.Minute
	s := time.Duration(now.Second()) * time.Second
	ns := time.Duration(now.Nanosecond()) * time.Nanosecond

	end := now.Add(-(h + m + s + ns))
	begin := end.Add(-24 * time.Hour)

	return interfaces.S.Points.GetUsersPeriodIncomePoints(userIds, begin, end)
}

func getYesterdayPointsRecord(userId int64) ([]*entities.UserPointsRecord, error) {
	now := time.Now()
	h := time.Duration(now.Hour()) * time.Hour
	m := time.Duration(now.Minute()) * time.Minute
	s := time.Duration(now.Second()) * time.Second
	ns := time.Duration(now.Nanosecond()) * time.Nanosecond

	end := now.Add(-(h + m + s + ns))
	begin := end.Add(-24 * time.Hour)

	return interfaces.S.Points.GetUserPeriodIncomePointsRecords(userId, begin, end)
}

func getTodayPoints(userIds []int64) (map[int64]int, error) {
	now := time.Now()
	h := time.Duration(now.Hour()) * time.Hour
	m := time.Duration(now.Minute()) * time.Minute
	s := time.Duration(now.Second()) * time.Second
	ns := time.Duration(now.Nanosecond()) * time.Nanosecond

	begin := now.Add(-(h + m + s + ns))

	return interfaces.S.Points.GetUsersPeriodIncomePoints(userIds, begin, now)
}

func getTotalPoints(userIds []int64) (map[int64]int, error) {
	start, err := time.Parse("2006-01-02", "1970-01-01")
	if err != nil {
		return nil, err
	}

	return interfaces.S.Points.GetUsersPeriodIncomePoints(userIds, start, time.Now())
}

func (f *Service) ListFriendsPoints(userId int64) ([]*service.FriendPoints, error) {
	relationships, err := listByParentID(userId)
	if err != nil {
		return nil, err
	}

	var ids []int64
	for _, r := range relationships {
		ids = append(ids, r.UserId)
	}

	all, err := utils.PromiseAll(func() (interface{}, error) {
		return getTotalPoints(ids)
	}, func() (interface{}, error) {
		return getTodayPoints(ids)
	}, func() (interface{}, error) {
		return getYesterdayPoints(ids)
	}, func() (interface{}, error) {
		return interfaces.S.User.GetListUserByIDs(ids)
	})

	if err != nil {
		return nil, err
	}

	total, ok := all[0].(map[int64]int)
	if !ok {
		return nil, errors.New("map[int64]int assert failed")
	}
	today, ok := all[1].(map[int64]int)
	if !ok {
		return nil, errors.New("map[int64]int assert failed")
	}
	yesterday, ok := all[2].(map[int64]int)
	if !ok {
		return nil, errors.New("map[int64]int assert failed")
	}
	users, ok := all[3].([]*entities.SimpleUser)

	var friends []*service.FriendPoints
	for _, user := range users {
		friends = append(friends, &service.FriendPoints{
			SimpleUser:              *user,
			PointsAcquiredYesterday: yesterday[user.ID],
			PointsAcquiredToday:     today[user.ID],
			PointsAcquiredTotal:     total[user.ID],
		})
	}

	return friends, nil
}

func (f *Service) HasFriend(userId int64) (bool, error) {
	return hasFriend(userId)
}
