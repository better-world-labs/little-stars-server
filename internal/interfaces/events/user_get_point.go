package events

import (
	"aed-api-server/internal/interfaces/entities"
	"aed-api-server/internal/pkg/domain/emitter"
	"encoding/json"
)

type UserGetPoint struct {
	PointsEventType entities.PointsEventType //积分事件类型
	UserId          int64                    //关联人
	Points          int                      //获得积分
}

func (e *UserGetPoint) GetUserId() int64 {
	return e.UserId
}

func (e *UserGetPoint) Decode(bt []byte) (emitter.DomainEvent, error) {
	var evt UserGetPoint
	err := json.Unmarshal(bt, &evt)
	if err != nil {
		return nil, err
	}
	return &evt, nil
}

func (e *UserGetPoint) Encode() ([]byte, error) {
	bt, err := json.Marshal(e)
	if err != nil {
		return nil, err
	}
	return bt, nil
}
