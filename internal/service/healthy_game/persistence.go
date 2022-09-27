package healthy_game

import "aed-api-server/internal/pkg/db"

//go:inject-component
func NewIPersistence() iPersistence {
	return &persistence{}
}

type persistence struct {
}

func (*persistence) userAnswerExisted(userId int64) (bool, error) {
	return db.Table("healthy_game_answer").Where("user_id=?", userId).Exist()
}

func (*persistence) saveAnswers(gameAnswer *healthyGameAnswer) error {
	_, err := db.Table("healthy_game_answer").Insert(gameAnswer)
	return err
}

func (*persistence) getLastAnswers(userId int64) (*healthyGameAnswer, error) {
	var a healthyGameAnswer
	has, err := db.Table("healthy_game_answer").Where("user_id = ?", userId).Limit(1).Desc("id").Get(&a)
	if err != nil {
		return nil, err
	}
	if has {
		return &a, nil
	}
	return nil, nil
}
