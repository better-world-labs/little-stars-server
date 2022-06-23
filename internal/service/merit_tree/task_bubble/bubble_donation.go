package task_bubble

import (
	"aed-api-server/internal/interfaces"
	"aed-api-server/internal/interfaces/entities"
	"aed-api-server/internal/interfaces/events"
	log "github.com/sirupsen/logrus"
	"time"
)

type DonationPoints struct{}

func (*DonationPoints) GetTriggerEvents() []events.UserBaseEvent {
	return []events.UserBaseEvent{
		//用户注册事件
		&events.FirstLoginEvent{},
	}
}

func (d *DonationPoints) ExecuteCondition(userId int64) (bool, *entities.TaskBubble) {
	count, err := countTreeTaskBubble(userId, TaskDonation)
	if err != nil {
		log.Error("countTreeTaskBubble err", userId, TaskDonation, err)
		return false, nil
	}
	if count >= 1 {
		return false, nil
	}

	recordCount, err := interfaces.S.Donation.CountUserRecord(userId)
	if err != nil {
		log.Error(" interfaces.S.Donation.CountUserRecord err", userId, TaskDonation, err)
		return false, nil
	}

	if recordCount > 1 {
		return false, nil
	}

	return true, &entities.TaskBubble{
		BubbleId:    TaskDonation,
		Name:        TaskDonationName,
		Points:      200,
		EffectiveAt: time.Now(),
	}
}
