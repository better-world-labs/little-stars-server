package response

type BusinessError struct {
	Code    int
	Message string
}

func (e *BusinessError) Error() string {
	return e.Message
}

//业务异常，全球变量定义
//======================

var (
	ErrorUserNotRegister = &BusinessError{Code: 1600, Message: "用户未注册"}

	ErrorNotEnoughWalkToConvent         = &BusinessError{Code: 1701, Message: "步数不足兑换"}
	ErrorSignEarlyTimeNotAllowed        = &BusinessError{Code: 1801, Message: "5:00-10:00才能打卡 早睡早起身体好"}
	ErrorSignEarlyTodayAlreadySignedYet = &BusinessError{Code: 1802, Message: "今日已打卡"}

	ErrorDonationNotStartYet = &BusinessError{Code: 1901, Message: "募集未开始"}
	ErrorDonationCompleted   = &BusinessError{Code: 1902, Message: "募集已完成"}
	ErrorDonationExpired     = &BusinessError{Code: 1903, Message: "募集已过期"}

	ErrorInsufficientBalance = &BusinessError{Code: 1100, Message: "积分余额不足"}

	ErrorConcurrentOperation = &BusinessError{Code: 2222, Message: "重复操作"}
)
