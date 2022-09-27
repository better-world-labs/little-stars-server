package point

import (
	"aed-api-server/internal/interfaces"
	"aed-api-server/internal/interfaces/entities"
	"aed-api-server/internal/interfaces/events"
	"aed-api-server/internal/pkg/db"
	"aed-api-server/internal/pkg/domain/emitter"
	"aed-api-server/internal/service/merit_tree/task_bubble"
	"encoding/json"
	"errors"
	"github.com/go-xorm/xorm"
	log "github.com/sirupsen/logrus"
	"math"
	"reflect"
	"sort"
	"sync"
	"time"
)

type scheduleCondition struct {
	Operator string  `json:"operator"`
	Field    string  `json:"field"`
	Value    float64 `json:"value"`
}

// 临时，通用的需要改进
func (condition *scheduleCondition) matchEarly(data interface{}) bool {
	value := reflect.ValueOf(data).Elem().FieldByName(condition.Field)
	v := float64(value.Interface().(int))

	switch condition.Operator {
	case "eq":
		return v == condition.Value
	case "gte":
		return v >= condition.Value
	case "gt":
		return v > condition.Value
	case "lte":
		return v <= condition.Value
	case "lt":
		return v < condition.Value
	}

	return false
}

func (condition *scheduleCondition) match(data interface{}) bool {
	value := reflect.ValueOf(data).Elem().FieldByName(condition.Field)
	v, ok := value.Interface().(float64)
	if !ok {
		return false
	}

	switch condition.Operator {
	case "eq":
		return v == condition.Value
	case "gte":
		return v >= condition.Value
	case "gt":
		return v > condition.Value
	case "lte":
		return v <= condition.Value
	case "lt":
		return v < condition.Value
	}

	return false
}

type scheduleRule struct {
	Name       string               `json:"name"`
	Points     int                  `json:"points"`
	Conditions []*scheduleCondition `json:"conditions"`
}

type schedule struct {
	PointsEventType entities.PointsEventType `xorm:"points_event_type"`
	Name            string                   `xorm:"name"`
	Description     string                   `xorm:"description"`
	Points          int                      `xorm:"points"`
	PeakPeriod      int                      `xorm:"peak_period"`
	PointsRule      string                   `xorm:"points_rule"`
	Show            bool                     `xorm:"show"`
	Sort            int                      `xorm:"Sort"`
	Rules           []*scheduleRule          `xorm:"-"`
}

var schedulesMap map[entities.PointsEventType]*schedule = nil

type pointsScheduler struct{}

//go:inject-component
func NewPointScheduler() *pointsScheduler {
	return &pointsScheduler{}
}

func (*pointsScheduler) ReloadSchedule() error {
	var schedules = make([]*schedule, 0)
	err := db.Table("point_event_define").Find(&schedules)
	if err != nil {
		return err
	}
	_schedulesMap := make(map[entities.PointsEventType]*schedule)
	for i := range schedules {
		sch := schedules[i]
		if sch.PointsRule != "" {
			err = json.Unmarshal([]byte(sch.PointsRule), &sch.Rules)
			if err != nil {
				return err
			}
		}
		_schedulesMap[sch.PointsEventType] = sch
	}
	var mu sync.Mutex
	mu.Lock()
	schedulesMap = _schedulesMap
	mu.Unlock()
	return nil
}

func (s *pointsScheduler) BuildPointsEventTypeSignEarly(userId, signEarlyId int64, days int) *events.PointsEvent {
	return &events.PointsEvent{
		PointsEventType: entities.PointsEventTypeSignEarly,
		UserId:          userId,
		Params: &events.PointsEventTypeSignEarlyParams{
			SignEarlyId: signEarlyId,
			Days:        days,
		},
	}
}

func (s *pointsScheduler) BuildPointsEventTypeGamePoints(userId, gameId int64, points int, description string) *events.PointsEvent {
	return &events.PointsEvent{
		PointsEventType: entities.PointsEventTypeGameReward,
		UserId:          userId,
		Params: &events.PointsEventTypeGamePoints{
			GameId:      gameId,
			Points:      points,
			Description: description,
		},
	}
}

func (s *pointsScheduler) BuildPointsEventTypeMockedExam(userId, examId int64, score int) *events.PointsEvent {
	return &events.PointsEvent{
		PointsEventType: entities.PointsEventTypeExam,
		UserId:          userId,
		Params: &events.PointsEventTypeMockedExamParams{
			ExamId: examId,
			Score:  float64(score),
		},
	}
}

//BuildPointsEventWalk 构造步行的积分事件
func (s *pointsScheduler) BuildPointsEventWalk(userId int64, todayWalk int, convertWalk int, convertedPoints int) *events.PointsEvent {
	points := int(math.Floor(float64(todayWalk-convertWalk) / entities.WalkConvertRatio))
	return &events.PointsEvent{
		UserId:          userId,
		PointsEventType: entities.PointsEventTypeWalk,
		Params: &events.PointsEventTypeWalkParams{
			TodayWalk:        todayWalk,
			ConvertWalk:      convertWalk,
			ConvertedPoints:  convertedPoints,
			Points:           points,
			WalkConvertRatio: entities.WalkConvertRatio,
		},
	}
}

//BuildPointsEventTypeFriendsAddPoint 构造好友加成的积分事件
func (s *pointsScheduler) BuildPointsEventTypeFriendsAddPoint(userId int64, friendsPointFlows []*entities.UserPointsRecord) *events.PointsEvent {
	points := 0
	for i := range friendsPointFlows {
		record := friendsPointFlows[i]
		points += record.Points
	}

	points = int(math.Ceil(float64(points*entities.FriendAddPointPercent) / 100))
	return &events.PointsEvent{
		UserId:          userId,
		PointsEventType: entities.PointsEventTypeFriendsAddPoint,
		Params: &events.PointsEventTypeFriendsAddPointParams{
			FriendsPointsFlows:    friendsPointFlows,
			Points:                points,
			FriendAddPointPercent: entities.FriendAddPointPercent,
		},
	}
}

//BuildPointsEventTypeActivityGive 构造活动赠与的积分事件
func (s *pointsScheduler) BuildPointsEventTypeActivityGive(userId int64, points int, description string) *events.PointsEvent {
	return &events.PointsEvent{
		UserId:          userId,
		PointsEventType: entities.PointsEventTypeWalk,
		Params: &events.PointsEventTypeActivityGiveParams{
			Points:      points,
			Description: description,
		},
	}
}

func (s *pointsScheduler) BuildPointsEventTypeClockInDevice(userId int64, clockInId int64, job *entities.UserTask) *events.PointsEvent {
	return &events.PointsEvent{
		UserId:          userId,
		PointsEventType: entities.PointsEventTypeClockInDevice,
		Params: &events.PointsEventTypeClockInDeviceParams{
			ClockInId: clockInId,
			Job:       job,
		},
	}
}

func (s *pointsScheduler) BuildPointsEventTypeReward(userId int64, jobId int64, points int, description string) *events.PointsEvent {
	return &events.PointsEvent{
		UserId:          userId,
		PointsEventType: entities.PointsEventTypeReward,
		Params: &events.PointsEventTypeRewardParams{
			JobId:       jobId,
			Points:      points,
			Description: description,
		},
	}
}

func (s *pointsScheduler) DealPointsEvent(evt *events.PointsEvent) (*entities.DealPointsEventRst, error) {
	if nil == schedulesMap {
		err := s.ReloadSchedule()
		if err != nil {
			return nil, err
		}
	}

	sch, ok := schedulesMap[evt.PointsEventType]
	if !ok {
		return nil, errors.New("不能匹配到积分策略")
	}

	var points = sch.Points
	var err error
	var expiredAt = time.Now().Add(time.Second * time.Duration(sch.PeakPeriod))
	var description = ""
	switch sch.PointsEventType {
	case entities.PointsEventTypeWalk:
		points, err = parseWalkParamsPoints(evt.Params, sch.Rules)

	case entities.PointsEventTypeSignEarly:
		points, err = parseSignEarlyPoints(evt.Params, sch.Rules)

	case entities.PointsEventTypeFriendsAddPoint:
		points, err = parseFriendAddPoints(evt.Params, sch.Rules)

	case entities.PointsEventTypeActivityGive:
		points, err = parseActivityGivePoints(evt.Params, sch.Rules)

	case entities.PointsEventTypeClockInDevice:
		points, err = parseClockInDevicePoint(evt.Params, sch.Rules)

	case entities.PointsEventTypeExam:
		points, err = parseMockExamPoint(evt.Params, sch.Rules)

	case entities.PointsEventTypeReadNews:
		var has bool
		has, err = interfaces.S.TaskBubble.HasReadNewsTask(evt.UserId)
		if !has {
			return nil, errors.New("not found task readNews")
		}
	case entities.PointsEventTypeReward:
		points, description, err = parsePointsEventTypeRewardPoints(evt.Params, sch.Rules)

	case entities.PointsEventTypeGameReward:
		points, description, err = parsePointsEventTypeGamePoints(evt.Params, sch.Rules)
	}

	if err != nil {
		log.Error("deal event params err:", err)
		return nil, err
	}

	var (
		totalPoints int
		unReceive   int
	)
	err = db.Transaction(func(session *xorm.Session) error {
		err = insertPoints(evt.UserId, points, sch.PointsEventType, description, evt.Params, expiredAt)
		if err != nil {
			return err
		}
		totalPoints, err = interfaces.S.Points.GetUserTotalPoints(evt.UserId)
		if err != nil {
			return err
		}
		unReceive, err = getUserTotalUnReceivedPoints(evt.UserId)
		return err
	})

	if err != nil {
		return nil, err
	}

	userGetPoint := events.UserGetPoint{
		PointsEventType: evt.PointsEventType,
		UserId:          evt.UserId,
		Points:          points,
	}
	err = emitter.Emit(&userGetPoint)
	if err != nil {
		log.Error("发送事件 events.UserGetPoint failed", userGetPoint)
	}

	completeTaskBubble(evt.PointsEventType, evt.UserId)

	return &entities.DealPointsEventRst{
		UserId:                 evt.UserId,
		PointsEventType:        evt.PointsEventType,
		PointsEventDescription: sch.Description,
		PeckingPoints:          unReceive,
		PeckingPointsChange:    points,
		PeckedPoints:           totalPoints,
		PeckedPointsChange:     0,
	}, nil
}

func parsePointsEventTypeRewardPoints(params interface{}, rules []*scheduleRule) (int, string, error) {
	args, ok := params.(*events.PointsEventTypeRewardParams)
	if !ok {
		return 0, "", errors.New("event params type error")
	}
	return args.Points, args.Description, nil
}

func parsePointsEventTypeGamePoints(params interface{}, rules []*scheduleRule) (int, string, error) {
	args, ok := params.(*events.PointsEventTypeGamePoints)
	if !ok {
		return 0, "", errors.New("event params type error")
	}
	return args.Points, args.Description, nil
}

func completeTaskBubble(eventType entities.PointsEventType, userId int64) {
	taskBubbleId := 0
	switch eventType {
	case entities.PointsEventTypeFirstEnterAEDMap:
		taskBubbleId = task_bubble.TaskEnterMap
		break
	case entities.PointsEventTypeLearntVideo, entities.PointsEventTypeLearntCourse:
		taskBubbleId = task_bubble.TaskAidLearn
		break
	case entities.PointsEventTypeReadNews:
		taskBubbleId = task_bubble.TaskReadNews
		break
	case entities.PointsEventTypeDonationAward:
		taskBubbleId = task_bubble.TaskDonation
		break
	case entities.PointsEventTypeFirstEnterCommunity:
		taskBubbleId = task_bubble.TaskCommunity
		break
	}

	if taskBubbleId != 0 {
		err := interfaces.S.TaskBubble.CompleteTaskBubble(userId, taskBubbleId)
		if err != nil {
			log.Error("CompleteTaskBubble err", userId, taskBubbleId)
		}
	}
}

func parseMockExamPoint(params interface{}, rules []*scheduleRule) (int, error) {
	examParam, ok := params.(*events.PointsEventTypeMockedExamParams)
	if !ok {
		return 0, errors.New("invalid param type")
	}

	for _, rule := range rules {
		for _, condition := range rule.Conditions {
			// 类型通用化，强约束，否则传值类型，使用其他实体会报错
			if condition.match(examParam) {
				return rule.Points, nil
			}
		}
	}

	return 0, errors.New("no rules hit")
}

func parseClockInDevicePoint(params interface{}, rules []*scheduleRule) (int, error) {
	deviceParams := params.(*events.PointsEventTypeClockInDeviceParams)
	if deviceParams.Job != nil {
		return int(deviceParams.Job.Point), nil
	}
	return schedulesMap[entities.PointsEventTypeClockInDevice].Points, nil
}

func parseActivityGivePoints(params interface{}, conditions []*scheduleRule) (int, error) {
	giveParams, ok := params.(*events.PointsEventTypeActivityGiveParams)
	if !ok {
		return 0, errors.New("活动赠送积分事件 参数 字段错误")
	}
	return giveParams.Points, nil
}

func parseFriendAddPoints(params interface{}, conditions []*scheduleRule) (int, error) {
	pointParams, ok := params.(*events.PointsEventTypeFriendsAddPointParams)
	if !ok {
		return 0, errors.New("好友加成积分事件 参数 字段错误")
	}
	return pointParams.Points, nil
}

func parseSignEarlyPoints(params interface{}, conditions []*scheduleRule) (int, error) {
	earlyParams, ok := params.(*events.PointsEventTypeSignEarlyParams)
	if !ok {
		return 0, errors.New("早起积分事件 参数 字段错误")
	}
	for i := range conditions {
		rule := conditions[i]
		var suc = true
		for j := range rule.Conditions {
			suc = rule.Conditions[j].matchEarly(earlyParams)
			if !suc {
				break
			}
		}
		if suc {
			return rule.Points, nil
		}
	}
	return 0, errors.New("不能匹配签到规则")
}

func parseWalkParamsPoints(params interface{}, conditions []*scheduleRule) (int, error) {
	walkParams, ok := params.(*events.PointsEventTypeWalkParams)
	if ok {
		return walkParams.Points, nil
	}
	return 0, errors.New("步行兑换积分事件 参数 字段错误")
}

type Strategies []*entities.PointStrategy

func (s Strategies) Len() int {
	return len(s)
}

func (s Strategies) Less(i, j int) bool {
	return s[i].Sort < s[j].Sort
}

func (s Strategies) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

func (s *pointsScheduler) GetPointStrategies() ([]*entities.PointStrategy, error) {
	if nil == schedulesMap {
		err := s.ReloadSchedule()
		if err != nil {
			return nil, err
		}
	}

	strategies := make([]*entities.PointStrategy, 0)
	for k := range schedulesMap {
		s := schedulesMap[k]
		if s.Show {
			points := make([]int, 0)
			if s.Points > 0 {
				points = append(points, s.Points)
			} else if len(s.Rules) > 0 {
				for i := range s.Rules {
					points = append(points, s.Rules[i].Points)
				}
			}
			strategies = append(strategies, &entities.PointStrategy{
				Name:   s.Name,
				Points: points,
				Sort:   s.Sort,
			})
		}
	}

	sort.Sort(Strategies(strategies))
	return strategies, nil
}
