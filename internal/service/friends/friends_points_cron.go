package friends

import (
	"aed-api-server/internal/interfaces"
	"aed-api-server/internal/interfaces/facility"
	"aed-api-server/internal/pkg/domain/emitter"

	log "github.com/sirupsen/logrus"
)

func (f *Service) Cron(run facility.RunFuncOnceAt) {
	run("0 20 * * *", facility.FriendAddPoints, handleFriendsPoints)
}

func handleFriendsPoints() {
	all, err := listAll()
	if err != nil {
		log.Error("get friends list error:", err)
		return
	}

	groupByParentId := groupingByParentId(all)

	for k, v := range groupByParentId {
		err = doPointAddition(k, v)
		if err != nil {
			log.Error("execute friend add points err:", err)
		}
	}
}

func doPointAddition(parentId int64, friends []*Relationship) error {
	for _, f := range friends {
		err := doFriendPointAddition(parentId, f)
		if err != nil {
			return err
		}
	}

	return nil
}

func doFriendPointAddition(parentId int64, friend *Relationship) error {
	pointsRecords, err := getYesterdayPointsRecord(friend.UserId)
	if err != nil {
		return err
	}

	if len(pointsRecords) == 0 {
		log.Info("[friends.cron]", "no flow for user,skip", friend.UserId)
		return nil
	}

	pointEvent := interfaces.S.PointsScheduler.BuildPointsEventTypeFriendsAddPoint(
		parentId,
		pointsRecords,
	)
	return emitter.Emit(pointEvent)
}
