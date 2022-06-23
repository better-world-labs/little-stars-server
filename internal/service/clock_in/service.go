package clock_in

import (
	"aed-api-server/internal/interfaces"
	"aed-api-server/internal/interfaces/entities"
	"aed-api-server/internal/interfaces/events"
	"aed-api-server/internal/interfaces/service"
	"aed-api-server/internal/pkg"
	"aed-api-server/internal/pkg/db"
	"aed-api-server/internal/pkg/domain/emitter"
	"aed-api-server/internal/pkg/global"
	"aed-api-server/internal/pkg/location"
	"aed-api-server/internal/pkg/utils"
	log "github.com/sirupsen/logrus"
	"time"
)

const clockInTableName = "device_clock_in"

func NewService() *Service {
	return &Service{}
}

type Service struct {
	service.ClockInService
}

func (Service) GetDeviceClockInUserIds(deviceId string) ([]int64, error) {
	var userIds []int64

	var results []struct {
		CreatedBy int64
	}

	err := db.Table(clockInTableName).Distinct("created_by").
		Cols("created_by").
		Where("device_id=?", deviceId).Find(&results)

	for _, r := range results {
		userIds = append(userIds, r.CreatedBy)
	}

	return userIds, err
}

func (Service) GetDeviceClockInLatest2(deviceId string) ([]*entities.ClockIn, error) {
	infos := make([]*entities.ClockInBaseInfo, 0)
	err := db.Table(clockInTableName).Where("device_id=?", deviceId).Desc("id").Limit(2, 0).Find(&infos)
	if err != nil {
		return nil, err
	}

	_m := make(map[int64]*entities.SimpleUser)
	ins := make([]*entities.ClockIn, len(infos), len(infos))
	for i := range infos {
		info := infos[i]
		account := entities.SimpleUser{
			ID: info.CreatedBy,
		}
		_m[account.ID] = &account
		ins[i] = &entities.ClockIn{
			ClockInBaseInfo: *info,
			CreatedBy:       &account,
			CreatedAt:       global.FormattedTime(info.CreatedAt),
		}
	}

	err = patchSimpleAccount(_m)
	if err != nil {
		return nil, err
	}

	return ins, nil
}

func (Service) GetDeviceClockInStat() (*entities.DeviceClockInStat, error) {
	total, err := db.Table("device").Count()
	if err != nil {
		return nil, err
	}

	clockInCount, err := db.Table(clockInTableName).Select("count(distinct device_id)").Count()
	if err != nil {
		return nil, err
	}

	return &entities.DeviceClockInStat{
		Total: total,
		Todo:  total - clockInCount,
	}, nil
}

func (s Service) DoDeviceClockIn(info *entities.ClockInBaseInfo, userId int64) (rst []*entities.DealPointsEventRst, err error) {
	info.CreatedBy = userId
	info.CreatedAt = time.Now()
	_, err = db.Insert(clockInTableName, info)
	if err != nil {
		return
	}

	//加积分
	job, err := interfaces.S.Task.FindUserTaskByUserIdAndDeviceId(userId, info.DeviceId)
	if err != nil {
		return
	}

	var event *events.PointsEvent
	if job != nil && job.Status != 10 {
		err = interfaces.S.Task.CompleteJob(userId, job.Id)
		if err != nil {
			return
		}
		event = interfaces.S.PointsScheduler.BuildPointsEventTypeClockInDevice(userId, info.Id, job)
	} else {
		event = interfaces.S.PointsScheduler.BuildPointsEventTypeClockInDevice(userId, info.Id, nil)
	}

	if info.OpenIn != nil {
		err := interfaces.S.Device.UpdateDeviceOpenIn(info.DeviceId, *info.OpenIn)
		if err != nil {
			return nil, err
		}
	}
	//加积分
	times, err := interfaces.S.Points.GetUserPointsEventTimes(userId, entities.PointsEventTypeClockInDevice)
	if err != nil {
		return
	}

	log.Info("DoDeviceClockIn", "user has clockIn", times, "times")
	if times < pkg.UserPointsMaxTimesDeviceClockIn {
		pointRst, err := interfaces.S.PointsScheduler.DealPointsEvent(event)
		if err != nil {
			return nil, err
		}

		rst = append(rst, pointRst)
	}

	err = emitter.Emit(&events.ClockInEvent{ClockInBaseInfo: info})
	return
}

func (Service) GetDeviceClockInList(deviceId string) ([]*entities.ClockIn, error) {
	infos := make([]*entities.ClockInBaseInfo, 0)
	err := db.Table(clockInTableName).Where("device_id=?", deviceId).Desc("id").Find(&infos)
	if err != nil {
		return nil, err
	}

	_m := make(map[int64]*entities.SimpleUser)
	ins := make([]*entities.ClockIn, len(infos), len(infos))
	for i := range infos {
		info := infos[i]

		account, ok := _m[info.CreatedBy]
		if !ok {
			account = &entities.SimpleUser{
				ID: info.CreatedBy,
			}
			_m[account.ID] = account
		}
		ins[i] = &entities.ClockIn{
			ClockInBaseInfo: *info,
			CreatedBy:       account,
			CreatedAt:       global.FormattedTime(info.CreatedAt),
		}
	}

	err = patchSimpleAccount(_m)
	if err != nil {
		return nil, err
	}

	return ins, nil
}

func (s Service) GetDeviceLastClockIn(from location.Coordinate, deviceId string) (*entities.DeviceClockIn, error) {
	_map, err := s.GetBatchDeviceLastClockIn(from, []string{deviceId})
	if err != nil {
		return nil, err
	}
	return _map[deviceId], nil
}

func (Service) GetBatchDeviceLastClockIn(from location.Coordinate, deviceIds []string) (map[string]*entities.DeviceClockIn, error) {
	rstList, err := utils.PromiseAll(func() (interface{}, error) {
		return interfaces.S.Device.ListDevicesByIDs(from, deviceIds)
	}, func() (interface{}, error) {
		return findLastClockIn(deviceIds)
	})

	if err != nil {
		return nil, err
	}
	clockInMap := make(map[string]*entities.DeviceClockIn)
	ds := (rstList[0]).([]*entities.Device)
	if len(ds) == 0 {
		return clockInMap, nil
	}
	_map := rstList[1].(map[string]*findLastClockInItem)

	clockInsMap := make(map[int64]*entities.ClockIn, 0)
	clockInsIds := make([]int64, 0)
	accountsMap := make(map[int64]*entities.SimpleUser, 0)

	for i := range ds {
		d := ds[i]
		deviceClockIn := entities.DeviceClockIn{
			Device:                     d,
			IsLast2ClockInSame:         true,
			SupportExistedCount:        0,
			SupportNotExistedCount:     0,
			LastClockIns:               make([]*entities.ClockIn, 0),
			LastSupportExistedUsers:    make([]*entities.SimpleUser, 0),
			LastSupportNotExistedUsers: make([]*entities.SimpleUser, 0),
		}
		r := _map[d.Id]
		if r != nil {
			deviceClockIn.IsLast2ClockInSame = true
			deviceClockIn.SupportExistedCount = r.SupportExistedCount
			deviceClockIn.SupportNotExistedCount = r.SupportNotExistedCount
			deviceClockIn.LastSupportExistedUsers = userIdToSimpleAccountAndRecordMap(r.LastSupportExistedUserIds, accountsMap)
			deviceClockIn.LastSupportNotExistedUsers = userIdToSimpleAccountAndRecordMap(r.LastSupportNotExistedUserIds, accountsMap)
			deviceClockIn.LastClockIns, deviceClockIn.IsLast2ClockInSame, clockInsIds = covertFindLastClockInItemToLastClockInList(r, clockInsMap, clockInsIds)
		}
		clockInMap[d.Id] = &deviceClockIn
	}

	//补齐打卡信息
	err = patchClockIns(clockInsIds, clockInsMap, accountsMap)
	if err != nil {
		return nil, err
	}

	//补齐SimpleAccount
	err = patchSimpleAccount(accountsMap)
	if err != nil {
		return nil, err
	}
	return clockInMap, nil
}

func (Service) GetDeviceClockInPictures(deviceId string, sizeLimit int) ([]string, error) {
	infos := make([]*entities.ClockInBaseInfo, 0)
	session := db.Table(clockInTableName).
		Select("device_id, images").
		Where("device_id=? and images is not null", deviceId).
		Desc("id")
	if sizeLimit > 0 {
		session.Limit(sizeLimit)
	}
	err := session.Find(&infos)
	if err != nil {
		return nil, err
	}
	urls := make([]string, 0)
	for i := range infos {
		urls = append(urls, infos[i].Images...)
	}
	if sizeLimit > 0 {
		return urls[0:utils.IntMin(len(urls), sizeLimit)], nil
	}
	return urls, nil
}

func userIdToSimpleAccountAndRecordMap(userIds []int64, accountsMap map[int64]*entities.SimpleUser) []*entities.SimpleUser {
	//defer utils.TimeStat("userIdToSimpleAccountAndRecordMap")()
	accounts := make([]*entities.SimpleUser, 0, len(userIds))
	for j := range userIds {
		userId := userIds[j]
		account := accountsMap[userId]
		if account == nil {
			account = &entities.SimpleUser{
				ID: userId,
			}
			accountsMap[userId] = account
		}

		accounts = append(accounts, account)
	}
	return accounts
}

func covertFindLastClockInItemToLastClockInList(
	r *findLastClockInItem,
	clockInsMap map[int64]*entities.ClockIn,
	clockInsIds []int64,
) ([]*entities.ClockIn, bool, []int64) {
	//defer utils.TimeStat("covertFindLastClockInItemToLastClockInList")()
	lastClockIns := make([]*entities.ClockIn, 0, 2)
	isLast2ClockInSame := true

	if r.Last1Id > 0 {
		in := entities.ClockIn{
			ClockInBaseInfo: entities.ClockInBaseInfo{
				Id: r.Last1Id,
			},
		}
		clockInsIds = append(clockInsIds, r.Last1Id)
		clockInsMap[r.Last1Id] = &in
		lastClockIns = append(lastClockIns, &in)

		isLast2ClockInSame = !(r.Last2Result != r.Last1Result && r.Last2Id > 0)
		if !isLast2ClockInSame {
			in = entities.ClockIn{
				ClockInBaseInfo: entities.ClockInBaseInfo{
					Id: r.Last2Id,
				},
			}
			clockInsIds = append(clockInsIds, r.Last2Id)
			clockInsMap[r.Last2Id] = &in
			lastClockIns = append(lastClockIns, &in)
		}
	}
	return lastClockIns, isLast2ClockInSame, clockInsIds
}

func patchClockIns(clockInsIds []int64, clockInsMap map[int64]*entities.ClockIn, accountsMap map[int64]*entities.SimpleUser) error {
	defer utils.TimeStat("patchClockIns")()
	ins, err := findClockIns(clockInsIds)
	if err != nil {
		return err
	}
	for i := range ins {
		info := ins[i]
		in := clockInsMap[info.Id]
		in.ClockInBaseInfo = *info
		in.CreatedAt = global.FormattedTime(info.CreatedAt)

		account := accountsMap[info.CreatedBy]
		if account == nil {
			account = &entities.SimpleUser{
				ID: info.CreatedBy,
			}
			accountsMap[info.CreatedBy] = account
		}
		in.CreatedBy = account
	}
	return nil
}

func patchSimpleAccount(accountsMap map[int64]*entities.SimpleUser) error {
	defer utils.TimeStat("patchSimpleAccount")()
	userIds := make([]int64, 0, len(accountsMap))
	for k := range accountsMap {
		userIds = append(userIds, k)
	}
	ds, err := interfaces.S.User.GetListUserByIDs(userIds)
	if err != nil {
		return err
	}

	for i := range ds {
		src := ds[i]
		dst := accountsMap[src.ID]
		dst.Nickname = src.Nickname
		dst.Avatar = src.Avatar
	}
	return nil
}

func findClockIns(ids []int64) ([]*entities.ClockInBaseInfo, error) {
	infos := make([]*entities.ClockInBaseInfo, 0)
	err := db.Table(clockInTableName).In("id", ids).Find(&infos)
	if err != nil {
		return nil, err
	}
	return infos, err
}

type findLastClockInItem struct {
	DeviceId string
	Last1Id  int64
	Last2Id  int64

	Last1Result bool
	Last2Result bool

	SupportExistedCount    int //设备存在的打卡数量
	SupportNotExistedCount int //设备不存在的打卡数量

	LastSupportExistedUserIds    []int64
	LastSupportNotExistedUserIds []int64
}

func findLastClockIn(deviceIds []string) (map[string]*findLastClockInItem, error) {
	defer utils.TimeStat("findLastClockIn")()
	m := make(map[string]*findLastClockInItem)
	if len(deviceIds) == 0 {
		return m, nil
	}

	results := make([]*findLastClockInItem, 0)
	err := db.Table(clockInTableName).Select(`
			device_id,
			substring_index(GROUP_CONCAT(id order by id desc),',',1) as last1_id,
			substring_index(substring_index(GROUP_CONCAT(id order by id desc),',',2),',',-1) as last2_id,
	
			substring_index(GROUP_CONCAT(is_device_existed order by id desc),',',1) as last1_result,
			substring_index(substring_index(GROUP_CONCAT(is_device_existed order by id desc),',',2),',',-1) as last2_result,
			
			count(if(is_device_existed =1, 1,null)) as support_existed_count,
			count(if(is_device_existed =0, 1,null)) as support_not_existed_count,
	
			concat('[',substring_index(GROUP_CONCAT(distinct created_by order by is_device_existed desc,id desc),',',2),']') as last_support_existed_user_ids,
			concat('[',substring_index(GROUP_CONCAT(distinct created_by order by is_device_existed asc,id desc),',',2),']') as last_support_not_existed_user_ids
	`).In("device_id", deviceIds).GroupBy("device_id").Find(&results)

	if err != nil {
		return nil, err
	}

	for i := range results {
		r := results[i]
		supportExistedCount := utils.IntMin(r.SupportExistedCount, len(r.LastSupportExistedUserIds))
		supportNotExistedCount := utils.IntMin(r.SupportNotExistedCount, len(r.LastSupportNotExistedUserIds))
		r.LastSupportExistedUserIds = r.LastSupportExistedUserIds[0:supportExistedCount]
		r.LastSupportNotExistedUserIds = r.LastSupportNotExistedUserIds[0:supportNotExistedCount]
		m[r.DeviceId] = r
	}
	return m, nil
}
