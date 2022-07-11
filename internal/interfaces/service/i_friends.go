package service

import "aed-api-server/internal/interfaces/entities"

type (
	FriendPoints struct {
		entities.SimpleUser

		PointsAcquiredYesterday int `json:"pointsAcquiredYesterday"`
		PointsAcquiredToday     int `json:"pointsAcquiredToday"`
		PointsAcquiredTotal     int `json:"pointsAcquiredTotal"`
	}

	FriendsService interface {

		// ListFriendsPoints 读取用户好友列表以及积分加成
		ListFriendsPoints(userId int64) ([]*FriendPoints, error)

		// HasFriend 用户是否有好友
		HasFriend(userId int64) (bool, error)
	}
)
