package entities

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
)
