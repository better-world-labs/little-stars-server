package activity

import (
	"aed-api-server/internal/interfaces/entities"
	"aed-api-server/internal/interfaces/events"
	"aed-api-server/internal/pkg/global"
	"encoding/json"
	"time"
)

func CreateActivityVolunteerNotified(event *events.VolunteerNotifiedEvent) *entities.Activity {
	return &entities.Activity{
		Uuid:       event.Id,
		HelpInfoID: event.Aid,
		Class:      ClassVolunteerNotified,
		Record: map[string]interface{}{
			"aid":   event.Aid,
			"count": event.Count,
		},
		Created: global.FormattedTime(time.Now()),
	}
}

func CreateActivitySceneCalled(event *events.SceneCalledEvent) *entities.Activity {
	return &entities.Activity{
		Uuid:       event.Id,
		HelpInfoID: event.Aid,
		Class:      ClassSceneCalled,
		UserID:     &event.UserId,
		Created:    global.FormattedTime(event.Time),
	}
}

func CreateActivityDeviceGot(event *events.DeviceGotEvent) *entities.Activity {
	return &entities.Activity{
		Uuid:       event.Id,
		HelpInfoID: event.Aid,
		Class:      ClassDeviceGot,
		UserID:     &event.UserId,
		Created:    global.FormattedTime(event.Time),
	}
}

func CreateActivitySceneArrived(event *events.SceneArrivedEvent) *entities.Activity {
	return &entities.Activity{
		Uuid:       event.Id,
		HelpInfoID: event.Aid,
		Class:      ClassSceneArrived,
		Points:     event.Points,
		UserID:     &event.UserId,
		Created:    global.FormattedTime(event.Time),
	}
}

func CreateActivityGoingToScene(event *events.GoingToSceneEvent) *entities.Activity {
	return &entities.Activity{
		Uuid:       event.Id,
		HelpInfoID: event.Aid,
		Class:      ClassGoingToScene,
		UserID:     &event.UserId,
		Created:    global.FormattedTime(event.Time),
	}
}

func CreateActivityAidCalled(event *events.AidCalledEvent) *entities.Activity {
	return &entities.Activity{
		Uuid:       event.Id,
		HelpInfoID: event.Aid,
		Class:      ClassAidCalled,
		Created:    global.FormattedTime(event.Time),
	}
}

func CreateActivityGoingToGetDevice(event *events.GoingToGetDeviceEvent) *entities.Activity {
	return &entities.Activity{
		Uuid:       event.Id,
		HelpInfoID: event.Aid,
		Class:      ClassGoingToGetDevice,
		UserID:     &event.UserId,
		Created:    global.FormattedTime(time.Now()),
	}
}

func CreateActivitySceneReport(event *events.SceneReportEvent) (*entities.Activity, error) {
	var m map[string]interface{}
	bytes, err := json.Marshal(event)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(bytes, &m)
	if err != nil {
		return nil, err
	}

	return &entities.Activity{
		HelpInfoID: event.Aid,
		Class:      ClassSceneReport,
		Uuid:       event.Id,
		UserID:     &event.UserId,
		Record:     m,
		Created:    global.FormattedTime(time.Now()),
	}, nil
}
