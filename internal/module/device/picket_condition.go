package device

import (
	"aed-api-server/internal/interfaces"
	"github.com/ethereum/go-ethereum/log"
)

func init() {
	interfaces.S.PicketCondition = PicketCondition{NewStorage()}
}

type PicketCondition struct {
	storage Storage
}

func (c PicketCondition) IsPicketNone(deviceId string) bool {
	d, b, err := c.storage.GetDeviceByID(deviceId)
	if err != nil {
		log.Info("IsPicketNone error", err)
	}
	if !b {
		return false
	}
	if d.CredibleState == CredibleStatusDeviceNotFound {
		return true
	}
	return false
}

func (c PicketCondition) IsLastTwiceConflict(deviceId string) bool {
	d, b, err := c.storage.GetDeviceByID(deviceId)
	if err != nil {
		log.Info("IsPicketNone error", err)
	}
	if !b {
		return false
	}
	if d.CredibleState == CredibleStatusDeviceNotFound {
		return true
	}
	return false
}

func (c PicketCondition) IsLastTwiceFalse(deviceId string) bool {
	d, b, err := c.storage.GetDeviceByID(deviceId)
	if err != nil {
		log.Info("IsPicketNone error", err)
	}
	if !b {
		return false
	}
	if d.CredibleState == CredibleStatusDeviceError {
		return true
	}
	return false
}