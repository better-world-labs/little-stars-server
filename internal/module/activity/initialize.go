package activity

import (
	"aed-api-server/internal/interfaces/events"
	"aed-api-server/internal/pkg/domain/emitter"
)

func InitEventHandler() {
	emitter.On(&events.VolunteerNotifiedEvent{}, onVolunteerNotified)
	emitter.On(&events.AidCalledEvent{}, onAidCalled)
	emitter.On(&events.GoingToSceneEvent{}, onGoingToScene)
	emitter.On(&events.GoingToGetDeviceEvent{}, onGoingToGetDevice)
	emitter.On(&events.SceneCalledEvent{}, onSceneCalled)
	emitter.On(&events.DeviceGotEvent{}, onDeviceGot)
	emitter.On(&events.SceneArrivedEvent{}, onSceneArrived)
}
