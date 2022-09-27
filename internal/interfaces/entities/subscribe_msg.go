package entities

//SubscribeMessageKey 订阅消息Key
// 定义的值要和user_subscribe_message中的字段保持一致
type SubscribeMessageKey string

type NotifiedSubscription struct {
	Openid         string
	UserId         int64
	PointsExpiring int
	WalkExpiring   int
}

const (
	// SMkPointsExpiring 积分即将过期
	SMkPointsExpiring SubscribeMessageKey = "points_expiring"

	//SMkWalkExpiring 步行即将过期
	SMkWalkExpiring SubscribeMessageKey = "walk_expiring"

	//SMkGamePoints 游戏通知
	SMkGamePoints SubscribeMessageKey = "game_points"

	//SMKFeedComment 帖子被评论
	SMKFeedComment SubscribeMessageKey = "feed_comment"
)
