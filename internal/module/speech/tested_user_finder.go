package speech

import (
	"aed-api-server/internal/module/user"
	"aed-api-server/internal/pkg/location"
)

type testedUserFinder struct {
}

func NewTestedUserFinder() UserFinder {
	return &testedUserFinder{}
}

func (t *testedUserFinder) FindUser(position location.Coordinate) ([]*user.User, error) {
	return []*user.User{{
		Nickname: "F-S-W-J",
		Mobile:   "15548720906",
	}}, nil
}
