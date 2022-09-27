package healthy_game

import (
	"aed-api-server/internal/interfaces/entities"
	"aed-api-server/internal/interfaces/service"
	"aed-api-server/internal/pkg/utils"
	"time"
)

//go:inject-component
func NewIHealthyGame() service.IHealthyGame {
	return &svc{}
}

type svc struct {
	P    iPersistence        `inject:"-"`
	User service.UserService `inject:"-"`
}

func (s *svc) GetInfo(userId int64, fromUserUId string) (info *entities.HealthyGameInfo, err error) {
	todayDeadCount := time.Now().Sub(utils.TodayBegin()).Minutes()

	hadResult := false
	if userId > 0 {
		hadResult, err = s.P.userAnswerExisted(userId)
		if err != nil {
			return
		}
	}
	var shareUser *entities.ShareUser = nil
	if fromUserUId != "" {
		user, err := s.User.GetUserByUid(fromUserUId)
		if err != nil {
			return nil, err
		}
		if user != nil {
			shareUser = &entities.ShareUser{
				NickName: user.Nickname,
				Avatar:   user.Avatar,
			}
		}
	}

	return &entities.HealthyGameInfo{
		HadResult:      hadResult,
		TodayDeadCount: int(todayDeadCount),
		ShareUser:      shareUser,
		Questions:      questions,
	}, nil
}

func (s *svc) CommitAnswers(userId int64, answers []*entities.Answer) (*entities.Result, error) {
	result, score := runRule(answers)

	err := s.P.saveAnswers(&healthyGameAnswer{
		UserId:    userId,
		Answers:   answers,
		Result:    result,
		Score:     score,
		CreatedAt: time.Now(),
	})
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (s *svc) GetResult(userId int64) (*entities.Result, error) {
	a, err := s.P.getLastAnswers(userId)
	if err != nil {
		return nil, err
	}
	if a == nil {
		return nil, nil
	}
	return a.Result, nil
}
