package device

import (
	"aed-api-server/internal/interfaces/entities"
)

type (
	Storage interface {
		CreateDevice(device entities.Device) error
		GetDeviceByID(deviceId string) (*entities.Device, bool, error)
		ListDevicesByIDs(deviceIds []string) ([]*entities.Device, error)
		ListAllDevices() ([]*entities.Device, error)
	}
)
