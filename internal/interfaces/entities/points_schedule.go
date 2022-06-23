package entities

type PointsEventType int

const (
	PointsEventTypeWalk             PointsEventType = 51 //步行兑换
	PointsEventTypeSignEarly        PointsEventType = 52 //早起
	PointsEventTypeInvite           PointsEventType = 53 //邀请好友
	PointsEventTypeLearntVideo      PointsEventType = 54 //学习视频
	PointsEventTypeLearntCourse     PointsEventType = 55 //学习课程
	PointsEventTypeExam             PointsEventType = 56 //模拟考试
	PointsEventTypePublishAid       PointsEventType = 57 //发布求助信息
	PointsEventTypeArrived          PointsEventType = 58 //确认到达
	PointsEventTypeGotDevice        PointsEventType = 59 //获取设备
	PointsEventTypeUploadAidInfo    PointsEventType = 60 //上传现场信息
	PointsEventTypeAddDevice        PointsEventType = 61 //新增设备点
	PointsEventTypeClockInDevice    PointsEventType = 62 //设备点打卡
	PointsEventTypeCertificated     PointsEventType = 63 //认证成功
	PointsEventTypeDeviceGuide      PointsEventType = 64 //指路成功
	PointsEventTypeFriendsAddPoint  PointsEventType = 65 //好友加成
	PointsEventTypeActivityGive     PointsEventType = 66 //活动赠送
	PointsEventTypeBeInvited        PointsEventType = 67 //受邀奖励
	PointsEventTypeDonation         PointsEventType = 68 //项目捐献
	PointsEventTypeFirstEnterAEDMap PointsEventType = 69 //首次进入AED地图
	PointsEventTypeReadNews         PointsEventType = 70 //阅读咨询奖励
	PointsEventTypeDonationAward    PointsEventType = 71 //首次捐积分奖励
	PointsEventTypeTransaction      PointsEventType = 72 //积分交易
	PointsEventTypeSubscribe        PointsEventType = 73 //公众号
	PointsEventTypeReward           PointsEventType = 74 //积分奖励

	WalkConvertRatio      = 50 //步行兑换积分比率
	FriendAddPointPercent = 10 //好友加成比例
)

type DealPointsEventRst struct {
	UserId                 int64           `json:"userId"`                 //用户ID
	PointsEventType        PointsEventType `json:"pointsEventType"`        //积分类型
	PointsEventDescription string          `json:"pointsEventDescription"` //积分事件描述
	PeckingPoints          int             `json:"peckingPoints"`          //待领积分总数
	PeckingPointsChange    int             `json:"peckingPointsChange"`    //待领积分总数变化情况
	PeckedPoints           int             `json:"peckedPoints"`           //已领积分变化情况
	PeckedPointsChange     int             `json:"peckedPointsChange"`     //已领积分变化情况
}

type PointStrategy struct {
	Name   string `json:"name"`
	Points []int  `json:"points"`
	Sort   int    `json:"-"`
}

type PointsEventParams struct {
	RefTable   string
	RefTableId int64
}
