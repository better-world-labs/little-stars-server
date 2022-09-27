package device

import (
	"aed-api-server/internal/interfaces/entities"
	page "aed-api-server/internal/pkg/query"
)

type (
	Storage interface {
		CreateOrUpdateDevice(device *entities.BaseDevice) error
		GetDeviceByID(deviceId string) (*entities.BaseDevice, bool, error)
		Delete(ids []string) error
		ListDevicesByIDs(deviceIds []string) ([]*entities.BaseDevice, error)
		MapDevicesByIDs(deviceIds []string) (map[string]*entities.BaseDevice, error)
		ListAllDevices() ([]*entities.BaseDevice, error)
		PageDevices(page page.Query, keyword string) (page.Result[*entities.BaseDevice], error)
		ListLatestUserAddedDevices(latest int64) ([]*entities.BaseDevice, error)
	}
)
