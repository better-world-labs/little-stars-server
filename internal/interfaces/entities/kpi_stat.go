package entities

type KpiStatItem struct {
	//小程序用户数
	MiniProgramUserCount int

	//公众号粉丝数
	OffiaccountUserCount int

	//用户数
	UserCount int

	//完成培训用户数
	LearntUserCount int

	//用户日活跃数
	DailyUserCount int

	//次日留存率
	UserNextDayUseRatio float64

	//设备数
	DeviceCount int
}
