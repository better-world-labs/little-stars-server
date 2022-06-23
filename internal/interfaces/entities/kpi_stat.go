package entities

type KpiStatItem struct {
	//小程序用户数
	MiniProgramUserCount int `json:"miniProgramUserCount"`

	//公众号粉丝数
	OffiaccountUserCount int `json:"offiaccountUserCount"`

	//用户数
	UserCount int `json:"userCount"`

	//完成培训用户数
	LearntUserCount int `json:"learntUserCount"`

	//用户日活跃数
	DailyUserCount int `json:"dailyUserCount"`

	//次日留存率
	UserNextDayUseRatio float64 `json:"userNextDayUseRatio"`

	//设备数
	DeviceCount int `json:"deviceCount"`

	//注册用户数
	RegUserCount int64 `json:"regUserCount"`

	//授权手机号的用户数
	MobileUserCount int64 `json:"mobileUserCount"`
}

type UserPointsRank struct {
	UserId       int64 `json:"userId"`
	PointsAmount int   `json:"pointsAmount"`
	PointsCount  int   `json:"pointsCount"`
}

type UserPointsTop struct {
	List []*UserPointsRank `json:"list"`
	Text string            `json:"text"`
}
