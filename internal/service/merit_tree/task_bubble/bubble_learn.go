package task_bubble

import (
	"aed-api-server/internal/interfaces/events"
	"aed-api-server/internal/interfaces/service"
	"gitlab.openviewtech.com/openview-pub/gopkg/log"
	"time"
)

//急救学习

type AidLearning struct{}

func (*AidLearning) GetTriggerEvents() []events.UserEvent {
	return []events.UserEvent{
		//用户注册事件
		&events.FirstLoginEvent{},
	}
}

func (d *AidLearning) ExecuteCondition(userId int64) (bool, *service.TaskBubble) {
	count, err := countTreeTaskBubble(userId, TaskAidLearn)
	if err != nil {
		log.Error("countTreeTaskBubble err", userId, TaskAidLearn, err)
		return false, nil
	}

	if count >= 1 {
		return false, nil
	}

	return true, &service.TaskBubble{
		BubbleId:    TaskAidLearn,
		Name:        TaskAidLearnName,
		Points:      100,
		EffectiveAt: time.Now(),
	}
}
