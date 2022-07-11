package service

import (
	"aed-api-server/internal/interfaces/entities"
	"aed-api-server/internal/pkg/location"
	page "aed-api-server/internal/pkg/query"
	"io"
)

type DeviceService interface {
	// ListDevices 设备列表
	// @Param lng 经度
	// @Param lat 纬度
	// @Param distance 周围多少米
	// @Param page 开始页
	// @Param page 数量
	// @return []Device 设备列表
	ListDevices(from location.Coordinate, distance float64, p page.Query) ([]*entities.Device, error)

	// ListDevicesWithoutDistance 不计算距离的设备列表
	// @Param lng 经度
	// @Param lat 纬度
	// @Param distance 周围多少米
	// @Param page 开始页
	// @Param page 数量
	// @return []Device 设备列表
	ListDevicesWithoutDistance(from location.Coordinate, distance float64, p page.Query) ([]*entities.Device, error)

	// ListDevicesByIDs  find list by multi deviceId
	// @Param from start point
	// @Param deviceIds device ids
	// @return []Device devices
	ListDevicesByIDs(from location.Coordinate, deviceIds []string) ([]*entities.Device, error)

	// UpdateCredibleState 更新纠查状态
	// @param deviceId 设备ID
	// @param credibleState 纠查状态
	// @param timestamp 纠查时间
	// @return error
	UpdateCredibleState(deviceId string, credibleState int, timestamp int64) error

	// AddDevice 添加设备
	// @Param device AddDevice
	// @return []Device 设备列表
	AddDevice(accountId int64, device *entities.AddDevice) ([]*entities.DealPointsEventRst, error)

	// InfoDevice 设备详情
	// @Param lng 经度
	// @Param lat 纬度
	// @Param udid 设备唯一ID
	// @return Device 设备
	InfoDevice(lnt, lat float64, udid string) (*entities.Device, error)

	// AddGuideInfo 添加指路
	// @Param accountId 用户ID
	// @Param deviceId 设备唯一ID
	// @Param desc 路径描述
	// @Param remark 关键路径点
	// @Param pic 图片
	// @return error
	AddGuideInfo(accountId int64, deviceId string, desc []string, remark []string, pic [][]string) ([]*entities.DealPointsEventRst, error)

	// GetDeviceGuideInfo 获取指路列表
	// @Param deviceId 设备唯一ID
	// @return DeviceGuideList,error
	GetDeviceGuideInfo(deviceId string) (entities.DeviceGuideList, error)

	// GetGuideInfoById 查询指路详情
	// @Param uid 指路信息id
	// @return DeviceGuideList,error
	GetGuideInfoById(uid string) (entities.DeviceGuideListItem, error)

	// GetDeviceGallery 设备图集
	// @param deviceId 设备ID
	// @param latest 最新 n 条
	// @return Gallery
	// @return Error
	GetDeviceGallery(deviceId string, latest int) ([]*entities.Gallery, error)

	// CountDeviceByCredibleState 设备计数
	// @return 每个 credibleStatus 的数量
	// @return Error
	CountDeviceByCredibleState() ([]*entities.PicketedDeviceCount, error)

	// UpdateClockInImage 更新一张设备打卡图片
	// @param deviceId 设备ID
	// @param clockInImg 一次打卡的第一张打卡图
	// @param timestamp 打卡图更新时间
	// @return Error
	UpdateClockInImage(deviceId string, clockInImg string, timestamp int64) error

	HasTodayAddDevice(userId int64) (bool, error)

	UpdateDeviceOpenIn(deviceId string, openIn entities.TimeRange) error

	ImportDevices(reader io.Reader) error

	SyncDevices() error
}
