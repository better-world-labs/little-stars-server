package service

import (
	"aed-api-server/internal/interfaces/entities"
)

type TreasureChestService interface {
	GetUserTreasureChest(userId int64) (*entities.TreasureChest, error)
	OpenTreasureChest(userId int64, treasureChestId int) error
	CreateTreasureChest(request *entities.TreasureChestCreateRequest) error
}
