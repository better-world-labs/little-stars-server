package task_bubble

import (
	"aed-api-server/internal/interfaces/entities"
	"aed-api-server/internal/interfaces/events"
	log "github.com/sirupsen/logrus"
	"time"
)

const (
	TriggerTimes = 5
)

type ReadNews struct {
}

func (*ReadNews) GetTriggerEvents() []events.UserBaseEvent {
	return []events.UserBaseEvent{
		//用户注册事件
		&events.FirstLoginEvent{},

		//用户获取积分
		&events.UserGetPoint{},
	}
}

func (r *ReadNews) ExecuteCondition(userId int64) (bool, *entities.TaskBubble) {
	//是否有5次
	if bubble, err := countTreeTaskBubble(userId, TaskReadNews); err != nil || bubble >= TriggerTimes {
		if err != nil {
			log.Error("countTreeTaskBubble err", userId, TaskReadNews, err)
		}
		return false, nil
	}

	//存在任务未完成，不发新任务
	has, err := isTodayHasTodoTaskBubble(userId, TaskReadNews)
	if err != nil {
		log.Error("isTodayHasTodoTaskBubble err", userId, TaskReadNews, err)
		return false, nil
	}
	if has {
		return false, nil
	}

	//今天有任务，任务发放到明天
	has, err = isTodayHasCompletedTaskBubble(userId, TaskReadNews)
	if err != nil {
		log.Error("isTodayHasCompletedTaskBubble err", userId, TaskReadNews, err)
		return false, nil
	}

	effectiveAt := time.Now()
	if has {
		t := time.Now()
		effectiveAt = time.Date(t.Year(), t.Month(), t.Day()+1, 0, 0, 0, 0, t.Location())
	}

	return true, &entities.TaskBubble{
		BubbleId:    TaskReadNews,
		Name:        TaskReadNewsName,
		Points:      20,
		EffectiveAt: effectiveAt,
	}
}
