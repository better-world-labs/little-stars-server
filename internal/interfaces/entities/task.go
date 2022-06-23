package entities

import "time"

type (
	UserTask struct {
		Id          int64          `json:"id"`
		Name        string         `json:"name"`
		Point       float64        `json:"point"`
		IsTimeLimit int            `json:"isTimeLimit"`
		TimeLimit   int            `json:"timeLimit"`
		IsRead      int            `json:"isRead"`
		Status      int            `json:"status"`
		Image       string         `json:"image"`
		Description string         `json:"description"`
		IsExpired   bool           `json:"isExpired"   xorm:"-"`
		DeviceId    string         `json:"-"`
		DeviceInfo  *DeviceClockIn `json:"deviceInfo"  xorm:"-"`
	}

	Job struct {
		Id          int64
		TaskId      int64
		UserId      int64
		IsRead      int8
		Status      int8
		CreatedAt   time.Time
		CompletedAt time.Time
		IsTimeLimit int8
		BeginLimit  time.Time
		EndLimit    time.Time
		DeviceId    string
		Points      int
		KeyHash     string
		Param       string
	}

	UserTypePage struct {
		Records    []*UserTask `json:"records"`
		NextCursor string      `json:"nextCursor"`
	}

	UserTaskStat struct {
		Unread int `json:"unread"`
		Todo   int `json:"todo"`
	}
)
