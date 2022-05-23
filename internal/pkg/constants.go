package pkg

const TraceHeaderKey = "X-Request-ID"
const AuthorizationHeaderKey = "Authorization"
const AccountIDKey = "AccountID"
const AccountKey = "Account"

// 价值行为积分获取次数限制
const (
	UserPointsMaxTimesDeviceGuideMaxTimes = 5
	UserPointsMaxTimesPublishHelpInfo     = 2
	UserPointsMaxTimesMarkDevice          = 3
	UserPointsMaxTimesDeviceClockIn       = 5
	UserPointsMaxTimesUploadScene         = 2
	UserPointsMaxTimesGetDevice           = 2
	UserPointsMaxTimesSceneArrived        = 2
)
