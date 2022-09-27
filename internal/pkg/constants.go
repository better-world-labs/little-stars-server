package pkg

const TraceHeaderKey = "X-Request-ID"
const AuthorizationHeaderKey = "Authorization"
const AccountIDKey = "AccountID"
const AccountKey = "Account"

// 价值行为积分获取次数限制
const (
	UserPointsMaxTimesDeviceGuideMaxTimes = 5
	UserPointsMaxTimesMockExam            = 5
	UserPointsMaxTimesPublishHelpInfo     = 2
	UserPointsMaxTimesMarkDevice          = 3
	UserPointsMaxTimesDeviceClockIn       = 5
	UserPointsMaxTimesUploadScene         = 2
	UserPointsMaxTimesGetDevice           = 2
	UserPointsMaxTimesSceneArrived        = 2
	UserPointsMaxTimesInvite              = 10
)

const (
	VoteCostPoints            = 200 //投票积分花费
	VoteCostPointsDescription = "兑换投票特权"
)
