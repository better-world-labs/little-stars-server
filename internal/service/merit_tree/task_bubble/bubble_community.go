package task_bubble

import (
	"aed-api-server/internal/interfaces/entities"
	"aed-api-server/internal/interfaces/events"
	log "github.com/sirupsen/logrus"
	"time"
)

//查看社区

type EnterCommunity struct{}

func (*EnterCommunity) GetTriggerEvents() []events.UserBaseEvent {
	return []events.UserBaseEvent{
		//用户注册事件
		&events.FirstLoginEvent{},
	}
}

func (d *EnterCommunity) ExecuteCondition(userId int64) (bool, *entities.TaskBubble) {
	count, err := countTreeTaskBubble(userId, TaskCommunity)
	if err != nil {
		log.Error("countTreeTaskBubble err", userId, TaskCommunity, err)
		return false, nil
	}

	if count >= 1 {
		return false, nil
	}

	return true, &entities.TaskBubble{
		BubbleId:    TaskCommunity,
		Name:        TaskCommunityName,
		Points:      80,
		EffectiveAt: time.Now(),
	}
}
