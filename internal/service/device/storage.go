package device

import (
	"aed-api-server/internal/interfaces/entities"
	"aed-api-server/internal/pkg/db"
	page "aed-api-server/internal/pkg/query"
	"aed-api-server/internal/pkg/utils"
	"errors"
	"fmt"
	"github.com/go-xorm/xorm"
)

type storage struct {
}

func NewStorage() Storage {
	return &storage{}
}

func (s storage) CreateOrUpdateDevice(device *entities.BaseDevice) (err error) {
	if device.Id == "" {
		return errors.New("invalid deviceId")
	}

	_, exists, err := s.GetDeviceByID(device.Id)
	if err != nil {
		return err
	}

	return db.Transaction(func(session *xorm.Session) error {
		if exists {
			_, err = session.Table("device").ID(device.Id).Update(device)
		} else {
			_, err = session.Table("device").Insert(device)
		}

		return err
	})
}

func (s storage) GetDeviceByID(deviceId string) (*entities.BaseDevice, bool, error) {
	session := db.GetSession()
	defer session.Close()

	var d entities.BaseDevice
	exists, err := session.Table("device").Where("id = ?", deviceId).Get(&d)
	return &d, exists, err
}

func (s storage) MapDevicesByIDs(deviceIds []string) (map[string]*entities.BaseDevice, error) {
	defer utils.TimeStat("storage.ListBaseDevicesByIDs")()
	session := db.GetSession()
	defer session.Close()

	m := make(map[string]*entities.BaseDevice, 0)
	err := session.Table("device").In("id", deviceIds).Find(&m)

	return m, err
}

func (s storage) ListDevicesByIDs(deviceIds []string) ([]*entities.BaseDevice, error) {
	defer utils.TimeStat("storage.ListBaseDevicesByIDs")()
	session := db.GetSession()
	defer session.Close()

	var arr []*entities.BaseDevice
	err := session.Table("device").In("id", deviceIds).Find(&arr)
	return arr, err
}

func (s storage) ListAllDevices() ([]*entities.BaseDevice, error) {
	session := db.GetSession()
	defer session.Close()

	var arr []*entities.BaseDevice
	err := session.Table("device").Find(&arr)
	return arr, err
}

func (s storage) PageDevices(query page.Query, keyword string) (page.Result[*entities.BaseDevice], error) {
	session := db.GetSession()
	defer session.Close()

	var arr []*entities.BaseDevice

	sess := session.Table("device")
	if keyword != "" {
		sess.Where("address like concat('%',?,'%') or title like concat('%',?,'%') or tel like concat('%',?,'%')", keyword, keyword, keyword)
	}

	count, err := sess.Limit(query.Size, (query.Page-1)*query.Size).FindAndCount(&arr)

	return page.NewResult[*entities.BaseDevice](arr, int(count)), err
}

func (s storage) ListLatestUserAddedDevices(latest int64) ([]*entities.BaseDevice, error) {
	session := db.GetSession()
	defer session.Close()

	var arr []*entities.BaseDevice
	err := session.Table("device").
		Where("create_by > 0").
		Desc("created").
		Limit(int(latest), 0).
		Find(&arr)

	return arr, err
}

func (s storage) Delete(ids []string) error {
	sql := fmt.Sprintf("delete from device where id in (%s)", db.ArrayPlaceholder(len(ids)))
	of := db.TupleOf(sql, ids)
	_, err := db.Exec(of...)
	return err
}
