package device

import (
	"aed-api-server/internal/interfaces/entities"
	"io"
)

type IDeviceImporter interface {
	ImportDevices(reader io.Reader) ([]*entities.Device, error)
}
