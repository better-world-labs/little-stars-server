package task_bubble

import (
	"aed-api-server/internal/interfaces"
	"aed-api-server/internal/interfaces/events"
	"aed-api-server/internal/interfaces/service"
	"aed-api-server/internal/pkg/db"
	"aed-api-server/internal/pkg/domain/emitter"
	"encoding/json"
	"github.com/go-xorm/xorm"
	"gitlab.openviewtech.com/openview-pub/gopkg/log"
	"time"
)

const (
	AddAEDTask = -1

	ClockInAEDTask     = -2
	ClockInAEDTaskName = "打卡设备"

	TaskEnterMap     = -3
	TaskEnterMapName = "附近AED"

	TaskAidLearn     = -4
	TaskAidLearnName = "急救学习"

	TaskReadNews     = -5
	TaskReadNewsName = "资讯头条"

	TaskDonation     = -6
	TaskDonationName = "捐献积分"
)

func Init() {
	interfaces.S.TaskBubble = &MeritTreeTaskTaskBubbleService{}
}

type TreeTaskBubble struct {
	Id           int64
	UserId       int64
	BubbleId     int64
	Name         string
	Points       int
	IsCompleted  bool
	CompletedAt  time.Time
	CreatedAt    time.Time
	CreatedEvent string
	EffectiveAt  time.Time
}

func countTreeTaskBubble(userId int64, bubbleId int64) (int, error) {
	count, err := db.GetEngine().SQL(`select count(1) from tree_task_bubble where user_id = ? and bubble_id=?`,
		userId, bubbleId).Count()
	if err != nil {
		return 0, err
	}
	return int(count), nil
}

func countEffectTaskBubbleByUserId(userId int64) (int, error) {
	count, err := db.GetEngine().SQL(`
		select count(1) 
		from tree_task_bubble 
		where 
			user_id = ?
			and now() > effective_at
			and is_completed = 0
		`,
		userId).Count()
	if err != nil {
		return 0, err
	}
	return int(count), nil
}

//是否今天有未完成的任务
func isTodayHasTodoTaskBubble(userId int64, bubbleId int64) (bool, error) {
	return db.Exist(`
	select 1 
	from tree_task_bubble 
	where 
		user_id = ? 
		and bubble_id= ?
		and effective_at > CURRENT_DATE
		and is_completed = 0
	`, userId, bubbleId)
}

//isTodayHasCompletedTaskBubble 是否有今天完成的任务
func isTodayHasCompletedTaskBubble(userId int64, bubbleId int64) (bool, error) {
	return db.Exist(`
	select 1 
	from tree_task_bubble 
	where 
		user_id = ? 
		and bubble_id= ?
		and completed_at > CURRENT_DATE
		and is_completed = 1
	`, userId, bubbleId)
}

func findBubblesByUserId(userId int64) (list []*TreeTaskBubble, err error) {
	err = db.SQL(`
		select *
		from tree_task_bubble
		where
			user_id = ?
			and now() > effective_at
			and is_completed = 0
	`, userId).Find(&list)
	if err != nil {
		return nil, err
	}
	return list, nil
}

func InitEventHandler() {
	defines := make([]service.MeritTreeTaskTaskBubbleDefine, 0)
	defines = append(defines,
		&DonationPoints{},
		&ReadNews{},
		&NearbyAED{},
		&AidLearning{},
	)

	for _, def := range defines {
		fn := defGenFn(def)
		for _, event := range def.GetTriggerEvents() {
			emitter.On(event.(emitter.DomainEvent), fn)
		}
	}
}

func defGenFn(def service.MeritTreeTaskTaskBubbleDefine) emitter.DomainEventHandler {
	return func(event emitter.DomainEvent) error {
		evt := event.(events.UserEvent)

		if condition, b := def.ExecuteCondition(evt.GetUserId()); condition {
			eventJson, _ := json.Marshal(event)
			taskBubble := TreeTaskBubble{
				UserId:       evt.GetUserId(),
				BubbleId:     b.BubbleId,
				Name:         b.Name,
				Points:       b.Points,
				EffectiveAt:  b.EffectiveAt,
				IsCompleted:  false,
				CreatedAt:    time.Now(),
				CreatedEvent: string(eventJson),
			}
			_, err := db.Insert("tree_task_bubble", taskBubble)
			if err != nil {
				log.Error("taskBubble insert failed:", taskBubble)
			}
		}
		return nil
	}
}

type MeritTreeTaskTaskBubbleService struct{}

//GetTreeBubblesCount 获取任务气泡数量
func (*MeritTreeTaskTaskBubbleService) GetTreeBubblesCount(userId int64) (int, error) {
	return countEffectTaskBubbleByUserId(userId)
}

//GetTreeBubbles 获取任务气泡
func (*MeritTreeTaskTaskBubbleService) GetTreeBubbles(userId int64) (bubbles []*service.Bubble, err error) {
	list, err := findBubblesByUserId(userId)
	if err != nil {
		return nil, err
	}

	for i := range list {
		item := list[i]
		bubbles = append(bubbles, &service.Bubble{
			Id:     item.BubbleId,
			Type:   service.BubbleTodoTask,
			Name:   item.Name,
			Points: item.Points,
		})
	}
	return bubbles, err
}

//CompleteTaskBubble 完成气泡任务
func (*MeritTreeTaskTaskBubbleService) CompleteTaskBubble(userId int64, bubbleId int) error {
	log.Info("CompleteTaskBubble", userId, bubbleId)

	err := db.Transaction(func(session *xorm.Session) error {
		_, err := session.Exec(`
			update tree_task_bubble 
			set	
				is_completed= 1,
				completed_at = now()
			where
				user_id = ?
				and bubble_id = ?
				and now() > effective_at
				and is_completed = 0
			`, userId, bubbleId)
		return err
	})

	return err
}

func (*MeritTreeTaskTaskBubbleService) HasReadNewsTask(userId int64) (bool, error) {
	count, err := db.GetEngine().SQL(`
		select count(1) 
		from tree_task_bubble 
		where 
			user_id = ?
			and bubble_id = ?
			and now() > effective_at
			and is_completed = 0
		`, userId, TaskReadNews).Count()

	if err != nil {
		return false, err
	}
	return count > 0, nil
}
