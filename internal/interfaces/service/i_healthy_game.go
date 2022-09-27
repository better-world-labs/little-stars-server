package service

import "aed-api-server/internal/interfaces/entities"

type IHealthyGame interface {
	GetInfo(userId int64, fromUserUId string) (*entities.HealthyGameInfo, error)

	CommitAnswers(userId int64, answers []*entities.Answer) (*entities.Result, error)

	GetResult(userId int64) (*entities.Result, error)
}
