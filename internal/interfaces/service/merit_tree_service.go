package service

import (
	"aed-api-server/internal/interfaces/entities"
	"aed-api-server/internal/interfaces/events"
)

type MeritTreeService interface {
	GetTreeInfo(userId int64) (*entities.MeritTreeInfo, error)
	GetTreeBubblesCount(userId int64) (int, error)
	ReadTreeBubblesCount(userId int64) error
	ReceiveBubblePoints(userId int64, bubbleId int64) (*entities.ReceiveBubblePointsRst, error)
}

type MeritTreeTaskTaskBubble interface {
	//GetTreeBubblesCount 获取任务气泡数量
	GetTreeBubblesCount(userId int64) (int, error)

	//GetTreeBubbles 获取任务气泡
	GetTreeBubbles(userId int64) ([]*entities.Bubble, error)

	//CompleteTaskBubble 完成气泡任务
	CompleteTaskBubble(userId int64, bubbleId int) error

	HasReadNewsTask(userId int64) (bool, error)
}

type MeritTreeTaskTaskBubbleDefine interface {
	GetTriggerEvents() []events.UserBaseEvent
	ExecuteCondition(userId int64) (bool, *entities.TaskBubble)
}
