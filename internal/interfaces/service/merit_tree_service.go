package service

import (
	"aed-api-server/internal/interfaces/events"
	"aed-api-server/internal/pkg/global"
	"time"
)

const (
	BubbleTodoTask         = 1
	BubblePointsNeedAccept = 2
)

type Bubble struct {
	Id        int64                `json:"id"`
	Type      int8                 `json:"type"` //气泡类型：1-待办；2-待领
	Name      string               `json:"name"`
	Points    int                  `json:"points"`
	ExpiredAt global.FormattedTime `json:"expiredAt"`
}

type MeritTreeInfo struct {
	TotalPoints             int       `json:"totalPoints"`
	TreeLevel               int       `json:"treeLevel"`
	FriendsAddPointsPercent int       `json:"friendsAddPointsPercent"`
	Bubbles                 []*Bubble `json:"bubbles"`
}

type ReceiveBubblePointsRst struct {
	TotalPoints int `json:"totalPoints"`
	TreeLevel   int `json:"treeLevel"`
}

type MeritTreeService interface {
	GetTreeInfo(userId int64) (*MeritTreeInfo, error)
	GetTreeBubblesCount(userId int64) (int, error)
	ReadTreeBubblesCount(userId int64) error
	ReceiveBubblePoints(userId int64, bubbleId int64) (*ReceiveBubblePointsRst, error)
}

type MeritTreeTaskTaskBubble interface {
	//GetTreeBubblesCount 获取任务气泡数量
	GetTreeBubblesCount(userId int64) (int, error)

	//GetTreeBubbles 获取任务气泡
	GetTreeBubbles(userId int64) ([]*Bubble, error)

	//CompleteTaskBubble 完成气泡任务
	CompleteTaskBubble(userId int64, bubbleId int) error

	HasReadNewsTask(userId int64) (bool, error)
}

type TaskBubble struct {
	BubbleId    int64
	Name        string
	Points      int
	EffectiveAt time.Time
}

type MeritTreeTaskTaskBubbleDefine interface {
	GetTriggerEvents() []events.UserEvent
	ExecuteCondition(userId int64) (bool, *TaskBubble)
}
