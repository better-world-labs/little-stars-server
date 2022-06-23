package task_bubble

import (
	"aed-api-server/internal/interfaces/entities"
	"aed-api-server/internal/interfaces/events"
	log "github.com/sirupsen/logrus"
	"time"
)

//附近AED

type NearbyAED struct{}

func (*NearbyAED) GetTriggerEvents() []events.UserBaseEvent {
	return []events.UserBaseEvent{
		//用户注册事件
		&events.FirstLoginEvent{},
	}
}

func (d *NearbyAED) ExecuteCondition(userId int64) (bool, *entities.TaskBubble) {
	count, err := countTreeTaskBubble(userId, TaskEnterMap)
	if err != nil {
		log.Error("countTreeTaskBubble err", userId, TaskEnterMap, err)
		return false, nil
	}

	if count >= 1 {
		return false, nil
	}

	return true, &entities.TaskBubble{
		BubbleId:    TaskEnterMap,
		Name:        TaskEnterMapName,
		Points:      100,
		EffectiveAt: time.Now(),
	}
}
