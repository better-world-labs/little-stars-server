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

// exam
var (
	ErrorExamUnSubmit     = &BusinessError{5000, "请继续作答"}
	ErrorExamOwnerError   = &BusinessError{5001, "参考人员不一致"}
	ErrorAlreadyCompleted = &BusinessError{5002, "考试已经提交,请勿重复提交"}
	ErrorUnknownQuestion  = &BusinessError{5003, "答案与试题不匹配"}
)

// donation
var (
	ErrorDonationNotStartYet = &BusinessError{Code: 1901, Message: "募集未开始"}
	ErrorDonationCompleted   = &BusinessError{Code: 1902, Message: "募集已完成"}
	ErrorDonationExpired     = &BusinessError{Code: 1903, Message: "募集已过期"}
)

// market
var (
	ErrorCommodityStockIsNotEnough = &BusinessError{Code: 5050, Message: "被抢完咯，看看其他商品呗"}
	ErrorCommodityNotReleased      = &BusinessError{Code: 5051, Message: "商品已下下架"}
)

// early
var (
	ErrorSignEarlyTimeNotAllowed        = &BusinessError{Code: 1801, Message: "5:00-10:00才能打卡 早睡早起身体好"}
	ErrorSignEarlyTodayAlreadySignedYet = &BusinessError{Code: 1802, Message: "今日已打卡"}
)

// vote
var (
	ErrorNoVoteChance  = &BusinessError{Code: 3030, Message: "投票次数已用完"}
	ErrorVoteCompleted = &BusinessError{Code: 3031, Message: "投票活动已结束"}
)

// walk
var (
	ErrorNotEnoughWalkToConvent = &BusinessError{Code: 1701, Message: "步数不足兑换"}
)

// user
var (
	ErrorUserNotRegister = &BusinessError{Code: 1600, Message: "用户未注册"}
)

// generic
var (
	ErrorNotFound            = &BusinessError{Code: 1404, Message: "not found"}
	ErrorInsufficientBalance = &BusinessError{Code: 1100, Message: "积分余额不足"}
	ErrorConcurrentOperation = &BusinessError{Code: 2222, Message: "重复操作"}
)

// clock_in
var (
	ErrorAlreadyClockIn = &BusinessError{6000, "今日已经打卡过了，明天再来吧"}
)

// other
var (
	ErrorTooFar = &BusinessError{2000, "操作失败，您与目的地距离过于远"}
)
