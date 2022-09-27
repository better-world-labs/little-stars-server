package device

import (
	"aed-api-server/internal/interfaces"
	"aed-api-server/internal/interfaces/entities"
	"aed-api-server/internal/interfaces/events"
	service2 "aed-api-server/internal/interfaces/service"
	"aed-api-server/internal/pkg"
	"aed-api-server/internal/pkg/cache"
	"aed-api-server/internal/pkg/db"
	"aed-api-server/internal/pkg/domain/emitter"
	"aed-api-server/internal/pkg/location"
	page "aed-api-server/internal/pkg/query"
	"aed-api-server/internal/pkg/redis"
	"aed-api-server/internal/pkg/utils"
	"aed-api-server/internal/pkg/wx_crypto"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/go-xorm/xorm"
	log "github.com/sirupsen/logrus"
	"io"
	"sync"
	"time"

	"github.com/google/uuid"
	om "github.com/wk8/go-ordered-map"
)

type service struct {
	storage  Storage
	importer IDeviceImporter

	User    service2.UserServiceOld `inject:"-"`
	User2   service2.UserService    `inject:"-"`
	ClockIn service2.ClockInService `inject:"-"`
	Env     string                  `conf:"server.env"`
}

var GeoKey = "aed_device_geo"

var RiskAreaDefine = []int{500, 1000, 2000}

func (s service) getGeoKey() string {
	return fmt.Sprintf("%s_%s", s.Env, GeoKey)
}

//go:inject-component
func NewService() service2.DeviceService {
	return &service{storage: NewStorage(), importer: NewExcelDeviceImporter()}
}

func (s service) ListLatestDevices(latest int64) ([]*entities.Device, error) {
	baseDevices, err := s.storage.ListLatestUserAddedDevices(latest)
	if err != nil {
		return nil, err
	}

	devices := entities.EnhanceBaseDevice(baseDevices)

	var deviceIds []string

	for _, device := range devices {
		deviceIds = append(deviceIds, device.Id)
	}

	s.assembleDevices(devices)

	return devices, s.doBatchParseDeviceInspectors(devices)

}

func (s service) ListDevicesByIDs(from location.Coordinate, deviceIds []string) ([]*entities.Device, error) {
	defer utils.TimeStat("service.ListDevicesByIDs")()
	if len(deviceIds) == 0 {
		return nil, nil
	}

	all, err := utils.PromiseAll(func() (interface{}, error) {
		return s.storage.ListDevicesByIDs(deviceIds)
	}, func() (interface{}, error) {
		return interfaces.S.ClockIn.BatchGetDeviceClockInUserIds(deviceIds)
	})
	if err != nil {
		return nil, err
	}

	baseDevices := all[0].([]*entities.BaseDevice)
	clockInUsers := all[1].(map[string][]int64)

	if err != nil {
		return nil, err
	}

	devices := entities.EnhanceBaseDevice(baseDevices)
	group := sync.WaitGroup{}
	group.Add(1)
	go func() {
		computedDistance(from, devices)
		group.Done()
	}()
	s.assembleDevices(devices)

	err = s.batchParseDeviceInspectors(devices, clockInUsers)
	if err != nil {
		return nil, err
	}

	group.Wait()
	return devices, nil
}

func (s service) ListDevicesEncrypted(userId int64, from location.Coordinate, distance float64, keyVersion int) (string, error) {
	all, err := utils.PromiseAll(func() (interface{}, error) {
		return s.User2.GetUserEncryptKey(userId, keyVersion)
	}, func() (interface{}, error) {
		return s.ListDevices(from, distance)
	})
	if err != nil {
		return "", err
	}

	if all[0] == nil {
		return "", errors.New("encrypt key not found")
	}

	key := all[0].(*entities.WechatEncryptKey)
	devices := all[1].([]*entities.Device)

	jsonBytes, err := json.Marshal(&page.Result[*entities.Device]{List: devices, Total: len(devices)})
	if err != nil {
		return "", err
	}

	return wx_crypto.Encrypt(key.EncryptKey, key.Iv, jsonBytes)
}

func (s service) ListDevices(from location.Coordinate, distance float64) ([]*entities.Device, error) {
	geoValues, err := redis.GeoRadius(s.getGeoKey(), from, distance)
	if err != nil {
		return nil, err
	}

	if len(geoValues) == 0 {
		return []*entities.Device{}, nil
	}

	baseDevices, err := s.storage.ListDevicesByIDs(entities.DistancedGeoValuesMapNames(geoValues))
	devices := make(map[string]*entities.Device, 0)

	for _, d := range baseDevices {
		devices[d.Id] = &entities.Device{BaseDevice: *d}
	}

	if err != nil {
		return nil, err
	}

	var res []*entities.Device
	for _, g := range geoValues {
		if device, ok := devices[g.Name]; ok {
			device.Distance = int64(g.Distance)
			s.assembleDevice(device)

			res = append(res, device)
		} else {
			s.handleDeviceNotExists(g)
		}
	}

	return res, s.doBatchParseDeviceInspectors(res)
}

func (s service) doBatchParseDeviceInspectors(devices []*entities.Device) error {
	var ids []string
	for _, d := range devices {
		ids = append(ids, d.Id)
	}

	clockInUsers, err := interfaces.S.ClockIn.BatchGetDeviceClockInUserIds(ids)
	if err != nil {
		return err
	}

	err = s.batchParseDeviceInspectors(devices, clockInUsers)
	if err != nil {
		return err
	}

	return nil
}

func (s service) handleDeviceNotExists(g *entities.DistancedGeoValue) {
	count, err := redis.GeoRemove(s.getGeoKey(), []string{g.Name})
	if err != nil {
		log.Errorf("handleDeviceNotExists: error: %v\n", err)
	} else {
		log.Infof("handleDeviceNotExists: remove %d", count)
	}
}

func computedDistance(from location.Coordinate, devices []*entities.Device) {
	positions := make([]location.Coordinate, len(devices))
	for i, e := range devices {
		positions[i] = location.Coordinate{Longitude: e.Longitude, Latitude: e.Latitude}
	}

	distances := location.DistanceFrom(from, positions)
	for i, e := range devices {
		e.Distance = distances[i]
	}
}

func (s service) assembleDevice(d *entities.Device) {
	//兼容
	d.Contract = d.Contact

	if d.Icon == "" {
		d.Icon = d.DeviceImage
	}

	if d.Icon == "" {
		d.Icon = d.EnvironmentImage
	}

	if d.Icon == "" {
		d.Icon = d.ClockInImage
	}
}

func (s service) assembleDevices(devices []*entities.Device) {
	for _, d := range devices {
		s.assembleDevice(d)
	}
}

func (s service) GetByID(deviceId string) (*entities.Device, bool, error) {
	baseDevice, exists, err := s.storage.GetDeviceByID(deviceId)
	if err != nil {
		return nil, exists, err
	}

	if !exists {
		return nil, exists, nil
	}

	return &entities.Device{BaseDevice: *baseDevice}, exists, err
}

func (s service) UpdateClockInImage(deviceId string, img string, timestamp int64) error {
	rows, err := db.Exec("update device set clock_in_image = ?, clock_in_image_timestamp  = ? where id = ? and clock_in_image_timestamp < ?", img, timestamp, deviceId, timestamp)
	affected, _ := rows.RowsAffected()
	log.Info("update", affected, "rows")
	return err
}

func (s service) UpdateCredibleState(deviceId string, credibleState int, timestamp int64) error {
	rows, err := db.Exec("update device set credible_state = ?, credible_state_timestamp = ? where id = ? and credible_state_timestamp < ?", credibleState, timestamp, deviceId, timestamp)
	affected, _ := rows.RowsAffected()
	log.Info("update", affected, "rows")
	return err
}
func (s service) GetUserAddDeviceTimes(userId int64) (int64, error) {
	count, err := db.Table("device").Where("create_by", userId).Count()
	return count, err
}

func (s service) AddDevice(accountId int64, device *entities.AddDevice) ([]*entities.DealPointsEventRst, error) {
	var pointRst []*entities.DealPointsEventRst
	var err error

	d := new(entities.Device)
	d.Id = uuid.NewString()
	d.Address = device.Address
	d.Longitude = device.Longitude
	d.Latitude = device.Latitude
	d.DeviceImage = device.DeviceImage
	d.EnvironmentImage = device.EnvironmentImage
	d.State = device.State
	d.Title = device.Title
	d.Contact = device.Contract
	d.CreateBy = accountId
	d.OpenIn = *device.OpenIn

	var desc []string
	var remarks []string
	pics := [][]string{}
	for _, v := range device.GuideInfo {
		desc = append(desc, v.Desc)
		remarks = append(remarks, v.Remark)
		pics = append(pics, v.Pic)
	}

	uid := uuid.NewString()
	guides := make([]DeviceGuide, 0)
	for k := range desc {
		b, _ := json.Marshal(pics[k])
		guides = append(guides, DeviceGuide{Uid: uid, AccountId: accountId, DeviceId: d.Id, Desc: desc[k], Remark: remarks[k], Pic: string(b)})
	}

	return pointRst, db.Transaction(func(session *xorm.Session) error {
		_, err = session.Insert(guides)
		if err != nil {
			return err
		}
		_, err = session.Insert(d)
		if err != nil {
			return err
		}

		_, err := redis.GeoAdd(s.getGeoKey(), []*entities.GeoValue{entities.GeoValueFromDevice(&d.BaseDevice)})
		if err != nil {
			return err
		}

		err = s.addPointDeviceMarked(d.Id, accountId, &pointRst)
		if err != nil {
			return err
		}

		if len(device.GuideInfo) > 0 {
			err := s.addPointDeviceGuided(d.Id, accountId, &pointRst)
			if err != nil {
				return err
			}
		}

		err = emitter.Emit(events.NewDeviceMarkedEvent(*d))
		return err
	})
}

func (s service) addPointDeviceGuided(deviceId string, userId int64, pointRst *[]*entities.DealPointsEventRst) error {
	times, err := interfaces.S.Points.GetUserPointsEventTimes(userId, entities.PointsEventTypeDeviceGuide)
	if err != nil {
		return err
	}

	if times < pkg.UserPointsMaxTimesDeviceGuideMaxTimes {
		event, err := interfaces.S.PointsScheduler.DealPointsEvent(&events.PointsEvent{
			PointsEventType: entities.PointsEventTypeDeviceGuide,
			UserId:          userId,
			Params: map[string]interface{}{
				"deviceId": deviceId,
			},
		})

		if err != nil {
			return err
		}
		*pointRst = append(*pointRst, event)
	}

	return nil
}

func (s service) addPointDeviceMarked(deviceId string, userId int64, pointRst *[]*entities.DealPointsEventRst) error {
	added, err := s.HasTimeBeforeAddedDevice(time.Minute*30, userId)
	if err != nil {
		return err
	}
	if added {
		return nil
	}

	times, err := interfaces.S.Points.GetUserPointsEventTimes(userId, entities.PointsEventTypeAddDevice)
	if err != nil {
		return err
	}

	if times > pkg.UserPointsMaxTimesMarkDevice {
		return nil
	}

	event, err := interfaces.S.PointsScheduler.DealPointsEvent(&events.PointsEvent{
		PointsEventType: entities.PointsEventTypeAddDevice,
		UserId:          userId,
		Params: map[string]interface{}{
			"deviceId": deviceId,
		},
	})
	if err != nil {
		return err
	}

	*pointRst = append(*pointRst, event)
	return nil
}

func (s service) batchParseDeviceInspectors(devices []*entities.Device, clockInUsers map[string][]int64) error {
	var userIds []int64
	for _, v := range clockInUsers {
		userIds = append(userIds, v...)
	}

	for _, d := range devices {
		userIds = append(userIds, d.CreateBy)
	}

	users, err := interfaces.S.User.GetMapUserByIDs(userIds)
	if err != nil {
		return err
	}

	for _, device := range devices {
		if createBy, ok := users[device.CreateBy]; ok {
			device.Inspector = append(device.Inspector, createBy)
		}

		if userIds, ok := clockInUsers[device.Id]; ok {
			for _, userId := range userIds {
				if user, ok := users[userId]; ok && device.CreateBy != userId {
					device.Inspector = append(device.Inspector, user)
				}
			}
		}
	}

	return nil
}

func (s service) parseDeviceInspectors(d *entities.Device) error {
	ids, err := interfaces.S.ClockIn.GetDeviceClockInUserIds(d.Id)
	if err != nil {
		return err
	}

	set := utils.NewSet[int64]()
	set.Add(d.CreateBy)
	set.AddAll(ids)

	slice := set.ToSlice()
	if len(slice) > 2 {
		slice = slice[:2]
	}

	users, err := interfaces.S.User.GetListUserByIDs(slice)
	if err != nil {
		return err
	}

	d.Inspector = users
	return nil
}

func (s service) InfoDevice(lnt, lat float64, udid string) (*entities.Device, error) {
	defer utils.TimeStat("InfoDevice")()
	sess := db.GetSession()
	defer sess.Close()

	d := new(entities.Device)
	has, err := sess.Where("id = ?", udid).Get(d)
	if err != nil {
		return d, err
	}
	if !has {
		return nil, fmt.Errorf("device id %v not found ", udid)
	}

	s.assembleDevice(d)

	if d.Icon == "" {
		latestPicketImage, err := interfaces.S.ClockIn.GetDeviceClockInPictures(udid, 1)
		if err != nil {
			return nil, err
		}
		if len(latestPicketImage) > 0 {
			d.Icon = latestPicketImage[0]
		}
	}

	distances := location.DistanceFrom(location.Coordinate{Longitude: lnt, Latitude: lat}, []location.Coordinate{{Longitude: d.Longitude, Latitude: d.Latitude}})
	d.Distance = distances[0]
	return d, s.parseDeviceInspectors(d)
}

func (s service) AddGuideInfo(accountId int64, deviceId string, desc []string, remark []string, pic [][]string) (pointRst []*entities.DealPointsEventRst, err error) {
	sess := db.GetSession()
	defer sess.Close()

	uid := uuid.NewString()
	guides := make([]DeviceGuide, 0)
	for k := range desc {
		b, _ := json.Marshal(pic[k])
		guides = append(guides, DeviceGuide{Uid: uid, AccountId: accountId, DeviceId: deviceId, Desc: desc[k], Remark: remark[k], Pic: string(b)})
	}

	_, err = sess.Insert(guides)
	if err != nil {
		return
	}

	err = s.addPointDeviceGuided(deviceId, accountId, &pointRst)
	return
}

func (s service) GetDeviceGuideInfo(deviceId string) (res entities.DeviceGuideList, err error) {
	sess := db.GetSession()
	defer sess.Close()

	guides := make([]DeviceGuide, 0)
	err = sess.Where("device_id = ?", deviceId).Find(&guides)
	if err != nil {
		return
	}

	var userIds []int64

	mInfos := om.New()
	for _, v := range guides {

		userIds = append(userIds, v.AccountId)
		if _, ok := mInfos.Get(v.Uid); !ok {
			mInfos.Set(v.Uid, make([]DeviceGuide, 0))
		}
		val, _ := mInfos.Get(v.Uid)
		list := val.([]DeviceGuide)
		list = append(list, v)
		mInfos.Set(v.Uid, list)
	}
	info := make([]entities.DeviceGuideListItem, 0)

	users, err := s.User.ListUserByIDs(userIds)
	if err != nil {
		return
	}

	for pair := mInfos.Newest(); pair != nil; pair = pair.Prev() {
		infos := pair.Value.([]DeviceGuide)
		userId := infos[0].AccountId
		//sess.Get(u)
		var listItem entities.DeviceGuideListItem

		if u, exists := users[userId]; exists {
			listItem.AccountId = infos[0].AccountId
			listItem.UserName = u.Nickname
			listItem.Avatar = u.Avatar
		}
		listItem.Uid = pair.Key.(string)
		listItem.Time = time.Time(infos[0].Created).Format("2006-01-02 15:04:05")
		for _, v := range infos {
			var pic []string
			err = json.Unmarshal([]byte(v.Pic), &pic)
			if err != nil {
				return
			}
			listItem.Info = append(listItem.Info, entities.GuideInfo{Desc: v.Desc, Remark: v.Remark, Pic: pic})
		}
		info = append(info, listItem)
	}

	return entities.DeviceGuideList{DeviceId: deviceId, Info: info}, nil
}

func (s service) GetGuideInfoById(uid string) (entities.DeviceGuideListItem, error) {
	sess := db.GetSession()
	defer sess.Close()

	guides := make([]DeviceGuide, 0)
	err := sess.Where("uid = ?", uid).Find(&guides)
	if err != nil {
		return entities.DeviceGuideListItem{}, err
	}

	mInfos := om.New()
	for _, v := range guides {
		if _, ok := mInfos.Get(v.Uid); !ok {
			mInfos.Set(v.Uid, make([]DeviceGuide, 0))
		}
		val, _ := mInfos.Get(v.Uid)
		list := val.([]DeviceGuide)
		list = append(list, v)
		mInfos.Set(v.Uid, list)
	}

	var listItem entities.DeviceGuideListItem

	for pair := mInfos.Oldest(); pair != nil; pair = pair.Next() {
		u := new(entities.User)
		infos := pair.Value.([]DeviceGuide)

		u.ID = infos[0].AccountId
		sess.Table("account").Get(u)

		uid = pair.Key.(string)
		listItem.Uid = uid
		listItem.AccountId = infos[0].AccountId
		listItem.UserName = u.Nickname
		listItem.Avatar = u.Avatar
		listItem.Time = time.Time(infos[0].Created).Format("2006-01-02 15:04:05")
		for _, v := range infos {
			var pic []string
			json.Unmarshal([]byte(v.Pic), &pic)
			listItem.Info = append(listItem.Info, entities.GuideInfo{Desc: v.Desc, Remark: v.Remark, Pic: pic})
		}
	}

	return listItem, nil
}

func (s service) GetDeviceGallery(deviceId string, latest int) ([]*entities.Gallery, error) {
	sql := `select id, 1 as type, origin as url
        from device
        where id = ? and origin != ''
        union all
        select id, 1 as type, env_origin as url
        from device
        where id = ? and env_origin != ''
    `
	if latest > 0 {
		sql += fmt.Sprintf("limit %d", latest)
	}

	res := make([]*entities.Gallery, 0)
	err := db.SQL(sql, deviceId, deviceId).Find(&res)
	if latest > 0 && len(res) >= latest {
		return res, nil
	}

	pictures, err := interfaces.S.ClockIn.GetDeviceClockInPictures(deviceId, latest-len(res))
	if err != nil {
		return nil, err
	}
	for _, url := range pictures {
		res = append(res, &entities.Gallery{
			Type: GalleryTypePicket,
			Url:  url,
		})
	}
	return res, err
}

func (s service) CountDeviceByCredibleState() ([]*entities.PicketedDeviceCount, error) {
	session := db.GetSession()
	session.Close()

	//res := map[int]int{}

	var arr = make([]*entities.PicketedDeviceCount, 0)
	sql := `select credible_state, count(id) count from device group by credible_state`
	err := session.SQL(sql).Find(&arr)
	return arr, err
}

func (s service) HasTodayAddDevice(userId int64) (bool, error) {
	exist, err := db.Table("device").Where("create_by = ? and created > UNIX_TIMESTAMP(CURRENT_DATE())", userId).Exist()
	if err != nil {
		return false, err
	}
	return exist, nil
}

func (s service) HasTimeBeforeAddedDevice(before time.Duration, userId int64) (bool, error) {
	seconds := before.Seconds()
	return db.Table("device").Where("create_by = ? and created > UNIX_TIMESTAMP() - ?", userId, seconds).Exist()
}

func (s service) UpdateDeviceOpenIn(deviceId string, openIn entities.TimeRange) error {
	d, exists, err := s.GetByID(deviceId)
	if err != nil {
		return err
	}

	if !exists {
		return errors.New("device not found")
	}

	d.OpenIn = openIn
	_, err = db.Table("device").ID(d.Id).Update(d)
	return err
}

func (s service) ImportDevices(reader io.Reader) error {
	devices, err := s.importer.ImportDevices(reader)
	if err != nil {
		return err
	}

	return s.CreateDevices(devices)
}

func (s service) CreateDevices(devices []*entities.BaseDevice) error {
	if len(devices) > 2000 {
		err := s.CreateDevices(devices[:2000])
		if err != nil {
			return err
		}

		err = s.CreateDevices(devices[2000:])
		if err != nil {
			return err
		}

		return nil
	}

	return db.Transaction(func(session *xorm.Session) error {
		_, err := session.Table("device").Insert(devices)
		if err != nil {
			return err
		}

		_, err = redis.GeoAdd(s.getGeoKey(), entities.GeoValuesFromDevices(devices))
		return err
	})
}

func (s service) SyncDevices() error {
	lock, err := cache.GetDistributeLock("device-sync", 1500000)
	if err != nil {
		return err
	}

	if !lock.Locked() {
		return errors.New("already doing")
	}

	cloudDevices, err := redis.GeoListNames(s.getGeoKey())
	if err != nil {
		return err
	}

	localDevices, err := s.storage.ListAllDevices()
	if err != nil {
		return err
	}

	log.Infof("SyncDevices: cloud: %d, local: %d\n", len(cloudDevices), len(localDevices))
	log.Infof("SyncDevices: len=%d", len(localDevices))
	cloudDeleting := checkCloudDeleting(cloudDevices, localDevices)
	cloudAdding := checkCloudAdding(cloudDevices, localDevices)

	log.Infof("diff: deleting: %d, adding: %d\n", len(cloudDeleting), len(cloudAdding))
	g := sync.WaitGroup{}
	defer func() {
		go func() {
			g.Wait()
			log.Infof("SyncDevices finished")
			err := lock.Release()
			if err != nil {
				log.Errorf("SyncDevices release lock error: %v", err)
			}
		}()
	}()

	if len(cloudDeleting) > 0 {
		g.Add(1)
		go func() {
			defer func() {
				g.Done()
				err := recover()
				if err != nil {
					log.Errorf("creating error: %v", err)
				} else {
					log.Infof("deleting finished\n")
				}
			}()
			_, err := redis.GeoRemove(s.getGeoKey(), cloudDeleting)
			if err != nil {
				log.Errorf("deleting error: %v", err)
			}
		}()
	}

	if len(cloudAdding) > 0 {
		g.Add(1)
		go func() {
			defer func() {
				g.Done()
				err := recover()
				if err != nil {
					log.Errorf("creating error: %v", err)
				} else {
					log.Infof("creating finished\n")
				}
			}()
			_, err := redis.GeoAdd(s.getGeoKey(), entities.GeoValuesFromDevices(cloudAdding))
			if err != nil {
				log.Errorf("creating error: %v", err)
			}
		}()
	}

	return nil
}

func checkCloudAdding(cloudDevices []string, localDevices []*entities.BaseDevice) []*entities.BaseDevice {
	cloudDeviceIds := utils.NewSet[string]()
	for _, cloudDevice := range cloudDevices {
		cloudDeviceIds.Add(cloudDevice)
	}

	var addingDevices []*entities.BaseDevice
	for _, localDevice := range localDevices {
		if _, ok := cloudDeviceIds[localDevice.Id]; !ok {
			addingDevices = append(addingDevices, localDevice)
		}
	}

	return addingDevices
}

func checkCloudDeleting(cloudDevices []string, localDevices []*entities.BaseDevice) []string {
	localDeviceMap := make(map[string]*entities.BaseDevice)
	for _, localDevice := range localDevices {
		localDeviceMap[localDevice.Id] = localDevice
	}

	var deletingDeviceIds []string
	for _, cloudDevice := range cloudDevices {
		if _, ok := localDeviceMap[cloudDevice]; !ok {
			deletingDeviceIds = append(deletingDeviceIds, cloudDevice)
		}
	}

	return deletingDeviceIds
}

func (s service) RiskArea(center location.Coordinate) ([]*entities.RiskArea, error) {
	var futures []utils.PromiseProcessor

	for _, def := range RiskAreaDefine {
		def := def // 再次踩坑！！！
		futures = append(futures, func() (interface{}, error) {
			return redis.GeoRadiusCount(s.getGeoKey(), center, float64(def))
		})
	}

	res, err := utils.PromiseAllArr(futures)
	if err != nil {
		return nil, err
	}

	var areas []*entities.RiskArea
	for i, def := range RiskAreaDefine {
		devices := res[i].(int)

		if i > 0 {
			devices = devices - res[i-1].(int)
		}

		areas = append(areas, &entities.RiskArea{Radius: def, Level: computeRiskLevel(devices, def)})
	}

	return areas, nil
}

func computeRiskLevel(devices, radius int) int {
	if devices < 1 {
		return entities.RiskLevelHigh
	}

	switch radius {
	case 500:
		return entities.RiskLevelLow

	case 1000:
		if devices >= 1 && devices < 3 {
			return entities.RiskLevelMedium
		}

		return entities.RiskLevelLow

	case 2000:
		if devices >= 1 && devices < 5 {
			return entities.RiskLevelMedium
		}

		return entities.RiskLevelLow

	default:
		return entities.RiskLevelLow
	}
}

func (s service) PageDevices(page *page.Query, keyword string) (page.Result[*entities.BaseDevice], error) {
	return s.storage.PageDevices(*page, keyword)
}

func (s service) GetById(id string) (*entities.BaseDevice, error) {
	d, exists, err := s.storage.GetDeviceByID(id)
	if err != nil {
		return nil, err
	}

	if !exists {
		return nil, errors.New("device not found")
	}

	return d, nil
}

func (s service) UpdateDevice(device *entities.BaseDevice) error {
	d, exists, err := s.storage.GetDeviceByID(device.Id)
	if err != nil {
		return err
	}

	if !exists {
		return errors.New("device not found ")
	}

	d.OpenIn = device.OpenIn
	d.Address = device.Address
	d.DeviceImage = device.DeviceImage
	d.Address = device.Address
	d.Title = device.Title
	d.Contact = device.Contact
	d.EnvironmentImage = device.EnvironmentImage
	d.Longitude = device.Longitude
	d.Latitude = device.Latitude

	return db.Transaction(func(session *xorm.Session) error {
		err := s.storage.CreateOrUpdateDevice(d)
		if err != nil {
			return err
		}

		_, err = redis.GeoUpdate(s.getGeoKey(), entities.GeoValuesFromDevices([]*entities.BaseDevice{d}))
		return err
	})
}

//TODO 以此代所有添加单个设备

func (s service) CreateDevice(d *entities.BaseDevice) error {
	d.Id = uuid.NewString()
	d.Created = time.Now().UnixMilli()
	d.Source = 0
	return db.Transaction(func(session *xorm.Session) error {
		err := s.storage.CreateOrUpdateDevice(d)
		if err != nil {
			return err
		}

		_, err = redis.GeoAdd(s.getGeoKey(), []*entities.GeoValue{entities.GeoValueFromDevice(d)})
		return err
	})
}

func (s service) DeleteDevices(ids []string) error {
	return db.Transaction(func(session *xorm.Session) error {
		err := s.storage.Delete(ids)
		if err != nil {
			return err
		}

		_, err = redis.GeoRemove(s.getGeoKey(), ids)
		return err
	})
}
