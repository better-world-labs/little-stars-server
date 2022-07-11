package speech

import (
	"aed-api-server/internal/interfaces/entities"
	"aed-api-server/internal/pkg/location"
)

type testedUserFinder struct {
}

func NewTestedUserFinder() UserFinder {
	return &testedUserFinder{}
}

func (t *testedUserFinder) FindUser(position location.Coordinate) ([]*entities.User, error) {
	return []*entities.User{{
		Nickname: "F-S-W-J",
		Mobile:   "15548720906",
	}}, nil
}
