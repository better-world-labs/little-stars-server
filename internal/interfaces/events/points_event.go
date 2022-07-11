package events

import (
	"aed-api-server/internal/interfaces/entities"
	"aed-api-server/internal/pkg/domain/emitter"
	"encoding/json"
	"reflect"
)

// 事件参数
type (

	//PointsEventTypeGamePoints 游戏积分事件参数
	PointsEventTypeGamePoints struct {
		GameId      int64  `json:"gameId"`
		Points      int    `json:"points"`
		Description string `json:"description"`
	}

	//PointsEventTypeMockedExamParams 早起积分事件参数
	PointsEventTypeMockedExamParams struct {
		ExamId int64   `json:"examId"`
		Score  float64 `json:"score"`
	}

	//PointsEventTypeWalkParams 步行兑换积分事件参数
	PointsEventTypeWalkParams struct {
		TodayWalk        int //今日步行量
		ConvertWalk      int //今日已经兑换步行量
		ConvertedPoints  int //今日已经兑换积分
		Points           int //积分
		WalkConvertRatio int //当前步行兑换积分比率
	}

	//PointsEventTypeSignEarlyParams 早起积分事件参数
	PointsEventTypeSignEarlyParams struct {
		SignEarlyId int64 //早起ID
		Days        int   //持续早起天数
	}

	//PointsEventTypeFriendsAddPointParams  好友加成参数
	PointsEventTypeFriendsAddPointParams struct {
		FriendsPointsFlows    []*entities.UserPointsRecord //好友积分流水
		Points                int                          //获取积分
		FriendAddPointPercent int                          //好友加成百分比
	}

	// PointsEventTypeActivityGiveParams 活动赠送参数
	PointsEventTypeActivityGiveParams struct {
		Description string //活动描述
		Points      int    //积分
	}

	// PointsEventTypeClockInDeviceParams 设备打卡参数
	PointsEventTypeClockInDeviceParams struct {
		ClockInId int64              `json:"clockInId"`
		Job       *entities.UserTask `json:"job"`
	}

	PointsEventTypeRewardParams struct {
		JobId       int64  `json:"jobId"`
		Points      int    `json:"points"`
		Description string `json:"description"`
	}
)

// TODO 所有事件都要有Encode/Decode 的单元测试

// PointsEvent 积分事件；将获取积分的行为
type PointsEvent struct {
	PointsEventType entities.PointsEventType //积分事件类型
	UserId          int64                    //关联人
	Params          interface{}              //参数：存盘备查
}

func (p *PointsEvent) MarshalJSON() ([]byte, error) {
	o := make(map[string]interface{}, 0)
	o["PointsEventType"] = p.PointsEventType
	o["UserId"] = p.UserId
	o["Params"] = p.Params
	return json.Marshal(&o)
}

func (p *PointsEvent) UnmarshalJSON(b []byte) error {
	m := make(map[string]*json.RawMessage, 0)
	err := json.Unmarshal(b, &m)
	if err != nil {
		return err
	}

	userIdRaw := m["UserId"]
	eventTypeRaw := m["PointsEventType"]
	paramsRaw := m["Params"]

	err = json.Unmarshal(*userIdRaw, &p.UserId)
	if err != nil {
		return err
	}

	err = json.Unmarshal(*eventTypeRaw, &p.PointsEventType)
	if err != nil {
		return err
	}

	switch p.PointsEventType {
	case entities.PointsEventTypeFriendsAddPoint:
		p.Params = &PointsEventTypeFriendsAddPointParams{}

	case entities.PointsEventTypeExam:
		p.Params = &PointsEventTypeMockedExamParams{}

	case entities.PointsEventTypeWalk:
		p.Params = &PointsEventTypeWalkParams{}

	case entities.PointsEventTypeActivityGive:
		p.Params = &PointsEventTypeActivityGiveParams{}

	case entities.PointsEventTypeClockInDevice:
		p.Params = &PointsEventTypeClockInDeviceParams{}

	case entities.PointsEventTypeSignEarly:
		p.Params = &PointsEventTypeSignEarlyParams{}

	case entities.PointsEventTypeReward:
		p.Params = &PointsEventTypeRewardParams{}

	case entities.PointsEventTypeGameReward:
		p.Params = &PointsEventTypeGamePoints{}

	default:
		p.Params = make(map[string]interface{}, 0)
	}

	if paramsRaw != nil {
		if reflect.TypeOf(p.Params).Kind() == reflect.Ptr {
			return json.Unmarshal(*paramsRaw, p.Params)
		} else {
			return json.Unmarshal(*paramsRaw, &p.Params)
		}
	}

	return nil
}

func (p *PointsEvent) Decode(bt []byte) (emitter.DomainEvent, error) {
	var evt PointsEvent
	err := json.Unmarshal(bt, &evt)
	if err != nil {
		return nil, err
	}
	return &evt, nil
}

func (p *PointsEvent) Encode() ([]byte, error) {
	bt, err := json.Marshal(p)
	if err != nil {
		return nil, err
	}
	return bt, nil
}
