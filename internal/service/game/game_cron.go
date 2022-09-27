package game

import (
	"aed-api-server/internal/interfaces"
	"aed-api-server/internal/interfaces/domains"
	"aed-api-server/internal/interfaces/entities"
	"aed-api-server/internal/interfaces/facility"
	"aed-api-server/internal/pkg/domain/emitter"
	"fmt"
	"github.com/sirupsen/logrus"
	"math"
)

func (g *GameImpl) Cron(run facility.RunFuncOnceAt) {
	run("0 * * * *", facility.GameProcess, func() {
		games, err := g.listEndUnsettledGames()
		if err != nil {
			logrus.Errorf("listEndedGames error: %v", err)
		}

		for _, game := range games {
			if !game.Settled {
				err := g.doSettle(game)
				if err != nil {
					logrus.Errorf("doSettle for game %d error: %v", game.Id, err)
				}
			}
		}
	})
}

func (g *GameImpl) doSettle(game *domains.Game) error {
	completed, err := g.playerPersistence.ListUserProcessesCompleted(game.Id)
	if err != nil {
		return err
	}

	for _, process := range completed {
		points := int(math.Round(float64(game.Points) / float64(len(completed))))
		err := g.pointsAward(game.Id, process.UserId, points)
		if err != nil {
			logrus.Errorf("doSettleFor process %d error: %v", process.Id, process)
			continue
		}

		err = g.ProgramNotifyWithPoints(process.UserId, points)
		if err != nil {
			logrus.Errorf("ProgramNotifyWithPoints process %d error: %v", process.Id, process)
		}
	}

	return g.updateSettled(game.Id, true)
}

func (g *GameImpl) ProgramNotifyWithoutPoints(userId int64) error {
	openid, err := g.User.GetUserOpenIdById(userId)
	if err != nil {
		return err
	}

	_, err = interfaces.S.SubscribeMsg.Send(userId, openid, entities.SMkGamePoints, map[string]interface{}{
		"thing1": map[string]interface{}{
			"value": "你的AED寻宝已完成，请等待瓜分积分",
		},
		"thing2": map[string]interface{}{
			"value": "请到首页及时收集积分，过期失效~",
		},
	}, "")

	return err

}
func (g *GameImpl) ProgramNotifyWithPoints(userId int64, points int) error {
	openid, err := g.User.GetUserOpenIdById(userId)
	if err != nil {
		return err
	}

	_, err = interfaces.S.SubscribeMsg.Send(userId, openid, entities.SMkGamePoints, map[string]interface{}{
		"thing1": map[string]interface{}{
			"value": fmt.Sprintf("你的AED寻宝已完成，获得%d积分", points),
		},
		"thing2": map[string]interface{}{
			"value": "请到首页及时收集积分，过期失效~",
		},
	}, "")

	return err

}

func (g *GameImpl) pointsAward(gameId, userId int64, points int) error {
	evt := interfaces.S.PointsScheduler.BuildPointsEventTypeGamePoints(userId, gameId, points, "游戏-瓜分积分奖励")
	return emitter.Emit(evt)
}
