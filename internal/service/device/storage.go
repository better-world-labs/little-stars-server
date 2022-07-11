package device

import (
	"aed-api-server/internal/interfaces/entities"
	"aed-api-server/internal/pkg/db"
	"aed-api-server/internal/pkg/utils"
	"errors"
)

type storage struct {
}

func NewStorage() Storage {
	return &storage{}
}

func (s storage) CreateDevice(device entities.Device) (err error) {
	if device.Id == "" {
		return errors.New("invalid deviceId")
	}

	session := db.GetSession()
	defer session.Close()

	_, exists, err := s.GetDeviceByID(device.Id)
	if err != nil {
		return err
	}

	if exists {
		_, err = session.Table("device").ID(device.Id).Update(device)
	} else {
		_, err = session.Table("device").Insert(&device)
	}

	return
}

func (s storage) GetDeviceByID(deviceId string) (*entities.Device, bool, error) {
	session := db.GetSession()
	defer session.Close()

	var d entities.Device
	exists, err := session.Table("device").Where("id = ?", deviceId).Get(&d)
	return &d, exists, err
}

func (s storage) ListDevicesByIDs(deviceIds []string) ([]*entities.Device, error) {
	defer utils.TimeStat("storage.ListDevicesByIDs")()
	session := db.GetSession()
	defer session.Close()

	var arr []*entities.Device
	err := session.Table("device").In("id", deviceIds).Find(&arr)
	return arr, err
}

func (s storage) ListAllDevices() ([]*entities.Device, error) {
	session := db.GetSession()
	defer session.Close()

	var arr []*entities.Device
	err := session.Table("device").Find(&arr)
	return arr, err
}
