package task

import (
	"aed-api-server/internal/interfaces/entities"
	"aed-api-server/internal/interfaces/events"
	"aed-api-server/internal/interfaces/facility"
	"aed-api-server/internal/pkg/domain/emitter"
	log "github.com/sirupsen/logrus"
)

func (t *Task) Listen(on facility.OnEvent) {
	//监听用户事件，匹配到浏览文章和浏览视频，触发积分事件
	on(&events.UserEvent{}, func(event emitter.DomainEvent) error {
		userEvent, ok := event.(*events.UserEvent)
		if !ok {
			return nil
		}

		if userEvent.EventType == entities.GetUserEventTypeOfReport(entities.Report_scanPage) {
			if len(userEvent.EventParams) > 0 {
				arg, ok := userEvent.EventParams[0].(map[string]interface{})
				if ok {
					pageUrl, ok := arg["pageUrl"].(string)
					if ok {
						scanPage(userEvent.UserId, pageUrl)
					}
				}

				return nil
			}
			log.Error("userEvent.EventParams err", userEvent)
		}

		if userEvent.EventType == entities.GetUserEventTypeOfReport(entities.Report_scanVideo) {
			if len(userEvent.EventParams) > 0 {
				arg, ok := userEvent.EventParams[0].(map[string]interface{})
				if ok {
					pageUrl, ok1 := arg["videoPageUrl"].(string)
					process, ok2 := arg["videoProgress"].(float64)

					if ok1 && ok2 {
						scanVideo(userEvent.UserId, pageUrl, int(process))
					}
				}

				return nil
			}
			log.Error("userEvent.EventParams err", userEvent)
		}
		return nil
	})

	//监听开宝箱事件，触发任务
	on(&events.UserOpenTreasureChest{}, func(event emitter.DomainEvent) error {
		chest, ok := event.(*events.UserOpenTreasureChest)
		if !ok {
			log.Error("convert event error")
			return nil
		}
		if chest.TaskId == 0 {
			log.Error("task =0，不创建任务")
			return nil
		}
		task := findTaskById(chest.TaskId)
		task.genJobByTreasureChest(chest)
		return nil
	})
}
