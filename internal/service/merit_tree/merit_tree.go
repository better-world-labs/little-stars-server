package merit_tree

import (
	"aed-api-server/internal/interfaces"
	"aed-api-server/internal/interfaces/entities"
	"aed-api-server/internal/pkg/db"
	"aed-api-server/internal/pkg/utils"
	"aed-api-server/internal/service/merit_tree/task_bubble"
	"time"
)

type meritTree struct{}

//go:inject-component
func NewMeritTreeService() *meritTree {
	return &meritTree{}
}

func getBubbles(userId int64) ([]*entities.Bubble, error) {
	results, err := utils.PromiseAll(func() (interface{}, error) {
		return interfaces.S.Points.GetUnReceivePoints(userId)
	}, func() (interface{}, error) {
		return interfaces.S.TaskBubble.GetTreeBubbles(userId)
	}, func() (interface{}, error) {
		return interfaces.S.Task.IsTodayHasBubble(userId)
	})

	if err != nil {
		return nil, err
	}

	bubbles := results[1].([]*entities.Bubble)

	//待办任务泡泡
	if results[2].(bool) {
		bubbles = append(bubbles, &entities.Bubble{
			Id:     task_bubble.ClockInAEDTask,
			Type:   entities.BubbleTodoTask,
			Name:   task_bubble.ClockInAEDTaskName,
			Points: 200,
		})
	}

	//带领积分泡泡
	points := results[0].([]*entities.UserPointsFlow)
	for i := range points {
		flow := points[i]
		bubbles = append(bubbles, &entities.Bubble{
			Id:        flow.Id,
			Type:      entities.BubblePointsNeedAccept,
			Name:      flow.Name,
			Points:    flow.Points,
			ExpiredAt: flow.ExpiredAt,
		})
	}

	return bubbles, nil
}

func (s *meritTree) GetTreeInfo(userId int64) (*entities.MeritTreeInfo, error) {
	utils.Go(func() {
		interfaces.S.User.RecordUserEvent(userId, entities.UserEventTypeGetTreeInfo)
	})

	rsts, err := utils.PromiseAll(func() (interface{}, error) {
		return getBubbles(userId)
	}, func() (interface{}, error) {
		return interfaces.S.Points.GetUserIncomePoints(userId)
	}, func() (interface{}, error) {
		return interfaces.S.Points.GetUserTotalPoints(userId)
	})

	if err != nil {
		return nil, err
	}

	bubbles := rsts[0].([]*entities.Bubble)
	incomePoints := rsts[1].(int)
	totalPoints := rsts[2].(int)

	if bubbles == nil {
		bubbles = make([]*entities.Bubble, 0)
	}

	info := entities.MeritTreeInfo{
		TotalPoints:             totalPoints,
		TreeLevel:               getTreeLevel(incomePoints),
		FriendsAddPointsPercent: getFriendsAddPointsPercent(userId),
		Bubbles:                 bubbles,
	}

	return &info, nil
}

func getFriendsAddPointsPercent(userId int64) int {
	has, err := interfaces.S.Friends.HasFriend(userId)
	if err != nil {
		return 0
	}
	if !has {
		return 0
	}
	return entities.FriendAddPointPercent
}

func getTreeLevel(points int) int {
	level := 0
	if 500 <= points && points < 2000 {
		level = 1
	} else if 2000 <= points && points < 5000 {
		level = 2
	} else if 5000 <= points && points < 10000 {
		level = 3
	} else if 10000 <= points && points < 50000 {
		level = 4
	} else if 50000 <= points {
		level = 5
	}
	return level
}

func (s *meritTree) GetTreeBubblesCount(userId int64) (int, error) {
	isRead, err := isReadBubble(userId)
	if err != nil {
		return 0, err
	}
	if isRead {
		return 0, nil
	}

	results, err := utils.PromiseAll(func() (interface{}, error) {
		return interfaces.S.Points.GetUnReceivePointsCount(userId)
	}, func() (interface{}, error) {
		return interfaces.S.TaskBubble.GetTreeBubblesCount(userId)
	}, func() (interface{}, error) {
		return interfaces.S.Task.IsTodayHasBubble(userId)
	})
	if err != nil {
		return 0, err
	}

	count := results[0].(int)
	count += results[1].(int)

	if results[2].(bool) {
		count++
	}
	return count, nil
}

func isReadBubble(userId int64) (bool, error) {
	exist, err := db.Table("tree_bubble_read_record").Where("user_id=? and created_at >= CURRENT_DATE()", userId).Exist()
	if err != nil {
		return false, err
	}
	return exist, nil
}

func (s *meritTree) ReadTreeBubblesCount(userId int64) error {
	type ReadRecord struct {
		UserId    int64     `xorm:"user_id"`
		CreatedAt time.Time `xorm:"created_at"`
	}

	_, err := db.Insert("tree_bubble_read_record", &ReadRecord{
		UserId:    userId,
		CreatedAt: time.Now(),
	})
	if err != nil {
		return err
	}
	return nil
}

func (s *meritTree) ReceiveBubblePoints(userId int64, bubbleId int64) (*entities.ReceiveBubblePointsRst, error) {
	err := interfaces.S.Points.ReceivePoints(userId, bubbleId)
	if err != nil {
		return nil, err
	}
	points, err := interfaces.S.Points.GetUserTotalPoints(userId)
	if err != nil {
		return nil, err
	}

	return &entities.ReceiveBubblePointsRst{
		TotalPoints: points,
		TreeLevel:   getTreeLevel(points),
	}, nil
}
