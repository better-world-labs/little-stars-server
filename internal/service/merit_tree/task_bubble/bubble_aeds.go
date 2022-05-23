package task_bubble

import (
	"aed-api-server/internal/interfaces/events"
	"aed-api-server/internal/interfaces/service"
	"gitlab.openviewtech.com/openview-pub/gopkg/log"
	"time"
)

//附近AED

type NearbyAED struct{}

func (*NearbyAED) GetTriggerEvents() []events.UserEvent {
	return []events.UserEvent{
		//用户注册事件
		&events.FirstLoginEvent{},
	}
}

func (d *NearbyAED) ExecuteCondition(userId int64) (bool, *service.TaskBubble) {
	count, err := countTreeTaskBubble(userId, TaskEnterMap)
	if err != nil {
		log.Error("countTreeTaskBubble err", userId, TaskEnterMap, err)
		return false, nil
	}

	if count >= 1 {
		return false, nil
	}

	return true, &service.TaskBubble{
		BubbleId:    TaskEnterMap,
		Name:        TaskEnterMapName,
		Points:      100,
		EffectiveAt: time.Now(),
	}
}
