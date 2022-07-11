package service

type PicketConditionService interface {
	IsPicketNone(deviceId string) bool
	IsLastTwiceConflict(deviceId string) bool
	IsLastTwiceFalse(deviceId string) bool
}
