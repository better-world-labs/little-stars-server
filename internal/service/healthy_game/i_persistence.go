package healthy_game

import (
	"aed-api-server/internal/interfaces/entities"
	"time"
)

type iPersistence interface {
	userAnswerExisted(userId int64) (bool, error)
	saveAnswers(gameAnswer *healthyGameAnswer) error
	getLastAnswers(userId int64) (*healthyGameAnswer, error)
}

type healthyGameAnswer struct {
	UserId    int64
	Answers   []*entities.Answer
	Result    *entities.Result
	Score     int
	CreatedAt time.Time
}
