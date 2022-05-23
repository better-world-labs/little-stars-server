package friends

import (
	"aed-api-server/internal/interfaces/service"
)

type (
	PointsDto struct {
		FriendsPointsAcquiredToday int                     `json:"friendsPointsAcquiredToday"`
		PointsAcquireTomorrow      int                     `json:"pointsAcquireTomorrow"`
		Friends                    []*service.FriendPoints `json:"friends"`
	}
)

func NewPointsDto(users []*service.FriendPoints) *PointsDto {
	var p PointsDto

	p.Friends = users
	for _, up := range users {
		p.FriendsPointsAcquiredToday += up.PointsAcquiredToday
	}

	p.PointsAcquireTomorrow = p.FriendsPointsAcquiredToday / 10

	if p.Friends == nil {
		p.Friends = []*service.FriendPoints{}
	}

	return &p
}
