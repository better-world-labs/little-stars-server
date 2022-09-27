package finder

import (
	"aed-api-server/internal/interfaces/entities"
	"aed-api-server/internal/pkg/location"
)

type UserFinder interface {
	FindUser(position location.Coordinate) ([]*entities.User, error)
}
