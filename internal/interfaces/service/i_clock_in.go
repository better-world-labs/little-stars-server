package service

import (
	"aed-api-server/internal/interfaces/entities"
	"aed-api-server/internal/pkg/location"
)

type ClockInService interface {
	// GetDeviceClockInStat 获取待打卡设备统计
	GetDeviceClockInStat() (*entities.DeviceClockInStat, error)

	// DoDeviceClockIn 对设备打卡
	DoDeviceClockIn(info *entities.ClockInBaseInfo, userId int64) ([]*entities.DealPointsEventRst, error)

	// GetDeviceClockInList 获取某个设备的打卡信息
	GetDeviceClockInList(deviceId string) ([]*entities.ClockIn, error)

	// GetDeviceClockInLatest2 获取某个设备的打卡信息
	GetDeviceClockInLatest2(deviceId string) ([]*entities.ClockIn, error)

	// GetDeviceLastClockIn 获取某个设备最近打卡信息
	GetDeviceLastClockIn(from location.Coordinate, deviceId string) (*entities.DeviceClockIn, error)

	// GetBatchDeviceLastClockIn 批量获取设备最近打卡信息
	GetBatchDeviceLastClockIn(from location.Coordinate, deviceIds []string) (map[string]*entities.DeviceClockIn, error)

	//GetDeviceClockInPictures 获取打卡图集
	GetDeviceClockInPictures(deviceId string, sizeLimit int) ([]string, error)

	//GetDeviceClockInUserIds 获取对设备打卡的所有用户ID
	GetDeviceClockInUserIds(deviceId string) ([]int64, error)
}
