package device

import (
	"aed-api-server/internal/interfaces"
	"aed-api-server/internal/interfaces/events"
	"aed-api-server/internal/interfaces/facility"
	"aed-api-server/internal/pkg/domain/emitter"
	"aed-api-server/internal/pkg/utils"
	"errors"
	log "github.com/sirupsen/logrus"
	"time"
)

func (*service) Listen(on facility.OnEvent) {
	on(&events.ClockInEvent{}, clockInEventHandler)
}

func clockInEventHandler(event emitter.DomainEvent) error {
	evt, ok := event.(*events.ClockInEvent)
	if !ok {
		return errors.New("get event failed")
	}

	log.Info("[device.EventHandler]: Handle ClockInEvent", evt.Id)
	return updateDeviceInfoForClockIn(evt.DeviceId, evt.Images, evt.CreatedAt)
}

func updateDeviceInfoForClockIn(deviceId string, images []string, at time.Time) error {
	_, err := utils.PromiseAll(func() (interface{}, error) {
		err := UpdateCredibleStatus(deviceId, at.UnixMilli())
		return nil, err
	}, func() (interface{}, error) {
		err := UpdateImage(deviceId, images, at.UnixMilli())
		return nil, err
	})

	return err
}

func UpdateImage(deviceId string, images []string, timestamp int64) error {
	if len(images) < 1 {
		log.Info("[device.EventHandler]: no clockIn image updated")
		return nil
	}
	img := images[0]
	err := interfaces.S.Device.UpdateClockInImage(deviceId, img, timestamp)
	return err
}

func UpdateCredibleStatus(deviceId string, timestamp int64) error {
	status, err := computeCredibleStatus(deviceId)
	if err != nil {
		return err
	}

	err = interfaces.S.Device.UpdateCredibleState(deviceId, status, timestamp)
	return err
}

func computeCredibleStatus(deviceId string) (int, error) {
	pickets, err := interfaces.S.ClockIn.GetDeviceClockInLatest2(deviceId)
	if err != nil {
		return 0, err
	}

	i := len(pickets)
	if i == 0 {
		return CredibleStatusDeviceNotClockIn, nil
	}

	if i == 1 {
		p := pickets[0]
		return GetCredibleStatus(p.IsDeviceExisted), nil
	}

	if pickets[0].IsDeviceExisted == pickets[1].IsDeviceExisted {
		return GetCredibleStatus(pickets[0].IsDeviceExisted), nil
	}

	return CredibleStatusDeviceError, nil
}
