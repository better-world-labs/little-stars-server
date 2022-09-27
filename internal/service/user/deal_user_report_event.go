package user

import (
	"aed-api-server/internal/interfaces"
	"aed-api-server/internal/interfaces/entities"
	"aed-api-server/internal/interfaces/events"
	"aed-api-server/internal/pkg/domain/emitter"
	log "github.com/sirupsen/logrus"
)

func dealUserEnterCommunity(userId int64) {
	var isSet bool
	err := interfaces.S.UserConfig.GetConfigToValue(userId, entities.UserConfigKeyFirstEnterCommunity, &isSet)
	if err != nil {
		log.Info("dealUserEnterAEDMap error", err)
		return
	}

	if !isSet {
		updated, err := interfaces.S.UserConfig.PutValueToConfig(userId, entities.UserConfigKeyFirstEnterCommunity, true)
		if updated && err == nil {
			err = emitter.Emit(&events.PointsEvent{
				PointsEventType: entities.PointsEventTypeFirstEnterCommunity,
				UserId:          userId,
				Params: entities.PointsEventParams{
					RefTable:   "user_config#" + entities.UserConfigKeyFirstEnterCommunity,
					RefTableId: userId,
				},
			})
		}
	}
}

func dealUserEnterAEDMap(userId int64) {
	var isSet bool
	err := interfaces.S.UserConfig.GetConfigToValue(userId, entities.UserConfigKeyFirstEnterAEDMap, &isSet)
	if err != nil {
		log.Info("dealUserEnterAEDMap error", err)
		return
	}

	if !isSet {
		updated, err := interfaces.S.UserConfig.PutValueToConfig(userId, entities.UserConfigKeyFirstEnterAEDMap, true)
		if updated && err == nil {
			err = emitter.Emit(&events.PointsEvent{
				PointsEventType: entities.PointsEventTypeFirstEnterAEDMap,
				UserId:          userId,
				Params: entities.PointsEventParams{
					RefTable:   "user_config#" + entities.UserConfigKeyFirstEnterAEDMap,
					RefTableId: userId,
				},
			})
		}
	}
}

func dealUserReadNews(userId int64) {
	has, err := interfaces.S.TaskBubble.HasReadNewsTask(userId)
	if err != nil {
		log.Error("TaskBubble.HasReadNewsTask error", err)
		return
	}

	if !has {
		return
	}

	err = emitter.Emit(&events.PointsEvent{
		PointsEventType: entities.PointsEventTypeReadNews,
		UserId:          userId,
	})
	if err != nil {
		log.Error("emit event error", err)
	}
}

func dealShowSubscribeQrCode(userId int64) {
	log.Info("dealShowSubscribeQrCode", "userId=", userId)
	var isSet bool
	err := interfaces.S.UserConfig.GetConfigToValue(userId, entities.UserConfigKeySubscribeOfficialAccounts, &isSet)
	log.Info("GetConfigToValue", "userId=", userId, "key=", entities.UserConfigKeySubscribeOfficialAccounts, "value=", isSet)
	if err != nil {
		log.Info("dealUserEnterAEDMap error", err)
		return
	}

	if !isSet {
		log.Info("PutValueToConfig", "userId=", userId)
		updated, err := interfaces.S.UserConfig.PutValueToConfig(userId, entities.UserConfigKeySubscribeOfficialAccounts, true)
		if updated && err == nil {
			log.Info("Emit", "userId=", userId)
			err = emitter.Emit(&events.PointsEvent{
				PointsEventType: entities.PointsEventTypeSubscribe,
				UserId:          userId,
			})
		}

		if err != nil {
			log.Error("emit event error", err)
		}
	}
}
