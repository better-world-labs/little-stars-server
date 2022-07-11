package device

import (
	"aed-api-server/internal/interfaces"
	"aed-api-server/internal/interfaces/entities"
	"aed-api-server/internal/interfaces/events"
	service2 "aed-api-server/internal/interfaces/service"
	"aed-api-server/internal/pkg"
	"aed-api-server/internal/pkg/db"
	"aed-api-server/internal/pkg/domain/emitter"
	"aed-api-server/internal/pkg/location"
	page "aed-api-server/internal/pkg/query"
	"aed-api-server/internal/pkg/tencent"
	"aed-api-server/internal/pkg/utils"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/go-xorm/xorm"
	log "github.com/sirupsen/logrus"
	"io"
	"sort"
	"sync"
	"time"

	"github.com/google/uuid"
	om "github.com/wk8/go-ordered-map"
)

type service struct {
	storage  Storage
	importer IDeviceImporter

	User    service2.UserServiceOld `inject:"-"`
	ClockIn service2.ClockInService `inject:"-"`
}

//go:inject-component
func NewService() service2.DeviceService {
	return &service{storage: NewStorage(), importer: NewExcelDeviceImporter()}
}

func (s service) ListDevicesByIDs(from location.Coordinate, deviceIds []string) ([]*entities.Device, error) {
	defer utils.TimeStat("service.ListDevicesByIDs")()
	if len(deviceIds) == 0 {
		return nil, nil
	}

	devices, err := s.storage.ListDevicesByIDs(deviceIds)
	if err != nil {
		return nil, err
	}

	group := sync.WaitGroup{}
	group.Add(1)
	go func() {
		computedDistance(from, devices)
		group.Done()
	}()
	err = s.assembleDevices(devices)
	if err != nil {
		return nil, err
	}

	group.Wait()
	return devices, nil
}

func (s service) ListDevicesWithoutDistance(from location.Coordinate, distance float64, p page.Query) ([]*entities.Device, error) {
	ids, err := tencent.ListRangeDeviceIDs(from, distance, p)
	if err != nil {
		return nil, err
	}

	if len(ids) == 0 {
		return []*entities.Device{}, nil
	}

	devices, err := s.storage.ListDevicesByIDs(ids)
	if err != nil {
		return nil, err
	}

	return devices, err
}

func (s service) ListDevices(from location.Coordinate, distance float64, p page.Query) ([]*entities.Device, error) {
	ids, err := tencent.ListRangeDeviceIDs(from, distance, p)
	if err != nil {
		return nil, err
	}

	if len(ids) == 0 {
		return []*entities.Device{}, nil
	}

	devices, err := s.storage.ListDevicesByIDs(ids)
	if err != nil {
		return nil, err
	}
	group := sync.WaitGroup{}
	group.Add(1)
	go func() {
		computedDistance(from, devices)
		sort.SliceStable(devices, func(i, j int) bool {
			return devices[i].Distance < devices[j].Distance
		})
		group.Done()
	}()
	err = s.assembleDevices(devices)
	if err != nil {
		return nil, err
	}

	group.Wait()
	return devices, err
}

func computedDistance(from location.Coordinate, devices []*entities.Device) {
	positions := make([]location.Coordinate, len(devices))
	for i, e := range devices {
		positions[i] = location.Coordinate{Longitude: e.Longitude, Latitude: e.Latitude}
	}

	distances := tencent.DistanceFrom(from, positions)
	for i, e := range devices {
		e.Distance = distances[i]
	}
}

func (s service) assembleDevices(devices []*entities.Device) error {
	size := len(devices)
	ids := make([]string, size)
	for i, d := range devices {
		ids[i] = d.Id
	}

	for _, d := range devices {
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

	return nil
}

func (s service) GetByID(deviceId string) (*entities.Device, bool, error) {
	return s.storage.GetDeviceByID(deviceId)
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
	// 腾讯云添加
	var pointRst []*entities.DealPointsEventRst
	udid, err := tencent.AddDevice(device.Longitude, device.Latitude, device.Title)
	if err != nil {
		return nil, err
	}

	// msyql添加
	d := new(entities.Device)
	d.Id = udid
	d.Address = device.Address
	d.Longitude = device.Longitude
	d.Latitude = device.Latitude
	d.DeviceImage = device.DeviceImage
	d.EnvironmentImage = device.EnvironmentImage
	d.State = device.State
	d.Title = device.Title
	d.Tel = device.Contract
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
		guides = append(guides, DeviceGuide{Uid: uid, AccountId: accountId, DeviceId: udid, Desc: desc[k], Remark: remarks[k], Pic: string(b)})
	}

	return pointRst, db.Begin(func(session *xorm.Session) error {
		_, err = session.Insert(guides)
		if err != nil {
			return err
		}
		_, err = session.Insert(d)
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

		err := emitter.Emit(events.NewDeviceMarkedEvent(*d))
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

func (s service) parseDeviceInspectors(d *entities.Device) error {
	ids, err := interfaces.S.ClockIn.GetDeviceClockInUserIds(d.Id)
	if err != nil {
		return err
	}

	set := utils.NewInt64Set()
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

	if d.Icon == "" {
		d.Icon = d.DeviceImage
	}

	if d.Icon == "" {
		d.Icon = d.EnvironmentImage
	}

	if d.Icon == "" {
		latestPicketImage, err := interfaces.S.ClockIn.GetDeviceClockInPictures(udid, 1)
		if err != nil {
			return nil, err
		}
		if len(latestPicketImage) > 0 {
			d.Icon = latestPicketImage[0]
		}
	}

	distances := tencent.DistanceFrom(location.Coordinate{Longitude: lnt, Latitude: lat}, []location.Coordinate{{Longitude: d.Longitude, Latitude: d.Latitude}})
	d.Distance = distances[0]
	err = s.parseDeviceInspectors(d)

	return d, err
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

	return db.Transaction(func(session *xorm.Session) error {
		_, err := session.Table("device").Insert(devices)
		if err != nil {
			return err
		}

		return tencent.CreateDevice(devices)
	})
}
func (s service) SyncDevices() error {
	//TODO implements
	return nil
}
