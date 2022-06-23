package entities

import (
	"aed-api-server/internal/pkg/global"
	"time"
)

const (
	BubbleTodoTask         = 1
	BubblePointsNeedAccept = 2
)

type TreasureChestStatus int

const (
	TreasureChestStatusInit      TreasureChestStatus = 0  //待展示
	TreasureChestStatusAvailable TreasureChestStatus = 10 //展示中
	TreasureChestStatusExpired   TreasureChestStatus = 20 //已失效|已完成
)

type TreasureChestCreateRequest struct {
	Name   string `json:"name"`   //宝箱名字
	Tips   string `json:"tips"`   //宝箱引导语
	Link   string `json:"link"`   //宝箱链接
	Points int    `json:"points"` //任务对应积分
	Sort   int    `json:"sort"`   //排序
	TaskId int64  `json:"taskId"` //任务ID
}

type TreasureChest struct {
	Id         int                 `json:"id"`         //宝箱ID
	Name       string              `json:"name"`       //宝箱名字
	Tips       string              `json:"tips"`       //宝箱引导语
	Link       string              `json:"link"`       //宝箱链接
	Points     int                 `json:"points"`     //任务对应积分
	Status     TreasureChestStatus `json:"status"`     //宝箱状态：0-待展示；10-展示中
	ValidAt    time.Time           `json:"validAt"`    //生效时间
	ValidTtl   int                 `json:"validTtl"`   //生效倒计时TTL
	ExpiredAt  time.Time           `json:"expiredAt"`  //过期时间
	ExpiredTtl int                 `json:"expiredTtl"` //过期时间倒计时TTL
}

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

type TaskBubble struct {
	BubbleId    int64
	Name        string
	Points      int
	EffectiveAt time.Time
}
