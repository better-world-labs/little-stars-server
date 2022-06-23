package friends

import (
	"aed-api-server/internal/interfaces"
	"aed-api-server/internal/pkg/cache"
	"aed-api-server/internal/pkg/domain/emitter"

	log "github.com/sirupsen/logrus"
)

func doCron() {
	lock, err := cache.GetDistributeLock("LOCK_FRIENDS_POINTS_CRON", 5000)
	if err != nil {
		log.Info("[friends.cron]", "doCron error: ", err)
	}

	defer lock.Release()

	if lock.Locked() {
		err := handleFriendsPoints()
		if err != nil {
			log.Error("doCron error: ", err)
		}
	}
}

func handleFriendsPoints() error {
	all, err := listAll()
	if err != nil {
		return err
	}

	//userIDs := mapUserIDs(all)
	//todayPointsMap, err := getYesterdayPoints(userIDs)
	//if err != nil {
	//	return err
	//}

	groupByParentId := groupingByParentId(all)

	for k, v := range groupByParentId {
		err := doPointAddition(k, v)
		if err != nil {
			return err
		}
	}

	return nil
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
