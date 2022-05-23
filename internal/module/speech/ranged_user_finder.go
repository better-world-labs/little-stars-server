package speech

import (
	"aed-api-server/internal/interfaces/entities"
	"aed-api-server/internal/interfaces/service"
	"aed-api-server/internal/module/device"
	"aed-api-server/internal/module/user"
	"aed-api-server/internal/pkg/base"
	"aed-api-server/internal/pkg/location"
	page "aed-api-server/internal/pkg/query"
	"aed-api-server/internal/pkg/tencent"
	"gitlab.openviewtech.com/openview-pub/gopkg/log"
	"sort"
)

type rangedUserFinder struct {
	deviceService service.DeviceService
}

func NewRangedUserFinder() UserFinder {
	return &rangedUserFinder{
		deviceService: device.NewService(),
	}
}

func (u rangedUserFinder) FindUser(position location.Coordinate) ([]*user.User, error) {
	userRange := int64(5000)
	devices, err := u.deviceService.ListDevices(position, float64(userRange), page.Query{})
	if err != nil {
		return nil, base.WrapError(Module, "get device list error", err)
	}

	var userIDs []int64
	if len(devices) == 0 {
		userRange, err = u.FindNearestDistance(position)
		if err != nil {
			return nil, err
		}
	}

	userIDs, err = FindRangeUsers(position, userRange)
	if err != nil {
		return nil, err
	}

	log.DefaultLogger().Infof("notify range %dm users", userRange)

	if len(userIDs) == 0 {
		return nil, nil
	}

	userMap, err := userService.ListUserByIDs(userIDs)
	if err != nil {
		return nil, err
	}

	var userArr []*user.User
	for _, v := range userMap {
		userArr = append(userArr, v)
	}

	return userArr, nil
}

func FindRangeUsers(current location.Coordinate, distance int64) ([]int64, error) {
	accountPositions, err := userService.ListAllPositions()
	if err != nil {
		return nil, err
	}

	if len(accountPositions) == 0 {
		return nil, nil
	}

	userCoordinate := make([]location.Coordinate, len(accountPositions))
	if err != nil {
		return nil, err
	}

	accountIDs := make([]int64, len(accountPositions))
	for i, p := range accountPositions {
		userCoordinate[i] = *p.Coordinate
	}

	distances := tencent.DistanceFrom(current, userCoordinate)
	//if err != nil {
	//	return nil, err
	//}

	for i := 0; i < len(distances); i++ {
		if distances[i] <= distance {
			accountIDs[i] = accountPositions[i].AccountID
		}
	}

	return accountIDs, nil
}

func (u rangedUserFinder) FindNearestDistance(current location.Coordinate) (int64, error) {
	devices, err := u.deviceService.ListDevices(current, 10000, page.Query{})
	if err != nil {
		return 0, err
	}

	if len(devices) == 0 {
		log.DefaultLogger().Warnf("no device found in 10000m range")
		return 10000, nil
	}

	sort.Sort(DeviceListNearestSort(devices))
	return devices[0].Distance, nil
}

type DeviceListNearestSort []*entities.Device

func (d DeviceListNearestSort) Len() int {
	return len(d)
}

func (d DeviceListNearestSort) Less(i, j int) bool {
	return d[i].Distance-d[j].Distance < 0
}

func (d DeviceListNearestSort) Swap(i, j int) {
	d[i], d[j] = d[j], d[i]
}
