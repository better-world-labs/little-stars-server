package task_bubble

import (
	"aed-api-server/internal/interfaces/entities"
	"aed-api-server/internal/interfaces/events"
	log "github.com/sirupsen/logrus"
	"time"
)

//急救学习

type AidLearning struct{}

func (*AidLearning) GetTriggerEvents() []events.UserBaseEvent {
	return []events.UserBaseEvent{
		//用户注册事件
		&events.FirstLoginEvent{},
	}
}

func (d *AidLearning) ExecuteCondition(userId int64) (bool, *entities.TaskBubble) {
	count, err := countTreeTaskBubble(userId, TaskAidLearn)
	if err != nil {
		log.Error("countTreeTaskBubble err", userId, TaskAidLearn, err)
		return false, nil
	}

	if count >= 1 {
		return false, nil
	}

	return true, &entities.TaskBubble{
		BubbleId:    TaskAidLearn,
		Name:        TaskAidLearnName,
		Points:      100,
		EffectiveAt: time.Now(),
	}
}
