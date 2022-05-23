package activity

import (
	"aed-api-server/internal/interfaces"
	"aed-api-server/internal/interfaces/events"
	"aed-api-server/internal/interfaces/service"
	"aed-api-server/internal/pkg/domain/emitter"
	"errors"
	"gitlab.openviewtech.com/openview-pub/gopkg/log"
)

var (
	ErrorInvalidEventType = errors.New("event assert failed, invalid event type")
)

func GetService() service.ActivityService {
	return interfaces.S.Activity
}

func onVolunteerNotified(event emitter.DomainEvent) error {
	if evt, ok := event.(*events.VolunteerNotifiedEvent); ok {
		log.Info("[activity.EventHandler]", "onVolunteerNotified", event)
		return GetService().SaveActivityVolunteerNotified(evt)
	}

	return ErrorInvalidEventType
}

func onAidCalled(event emitter.DomainEvent) error {
	if evt, ok := event.(*events.AidCalledEvent); ok {
		log.Info("[activity.EventHandler]", "onAidCalled", event)
		return GetService().SaveActivityAidCalled(evt)
	}

	return ErrorInvalidEventType
}

func onGoingToScene(event emitter.DomainEvent) error {
	if evt, ok := event.(*events.GoingToSceneEvent); ok {
		log.Info("[activity.EventHandler]", "onGoingToScene", event)
		return GetService().SaveActivityGoingToScene(evt)
	}

	return ErrorInvalidEventType
}

func onGoingToGetDevice(event emitter.DomainEvent) error {
	if evt, ok := event.(*events.GoingToGetDeviceEvent); ok {
		log.Info("[activity.EventHandler]", "onGoingToGetDevice", event)
		return GetService().SaveActivityGoingToGetDevice(evt)

	}

	return ErrorInvalidEventType
}

func onSceneCalled(event emitter.DomainEvent) error {
	if evt, ok := event.(*events.SceneCalledEvent); ok {
		log.Info("[activity.EventHandler]", "onSceneCalled", event)
		return GetService().SaveActivitySceneCalled(evt)
	}

	return ErrorInvalidEventType
}

func onDeviceGot(event emitter.DomainEvent) error {
	if evt, ok := event.(*events.DeviceGotEvent); ok {
		log.Info("[activity.EventHandler]", "onDeviceGot", event)
		_, err := GetService().SaveActivityDeviceGot(evt)
		return err
	}

	return ErrorInvalidEventType
}

func onSceneArrived(event emitter.DomainEvent) error {
	if evt, ok := event.(*events.SceneArrivedEvent); ok {
		log.Info("[activity.EventHandler]", "onSceneArrived", event)
		return GetService().SaveActivitySceneArrived(evt)
	}

	return ErrorInvalidEventType
}