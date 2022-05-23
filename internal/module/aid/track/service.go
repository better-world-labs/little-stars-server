package track

import (
	"aed-api-server/internal/pkg/db"
	"fmt"
	"github.com/go-xorm/xorm"
	"sync"
	"time"
)

var instance Service
var once sync.Once

func GetService() Service {
	once.Do(func() {
		instance = NewService()
	})

	return instance
}

func NewService() Service {
	return &service{}
}

type Service interface {
	GetUserTrackWithSession(session *xorm.Session, aidId int64, userId int64) (*Model, error)
	GetUserTrack(aidId int64, userId int64) (*Model, error)
	IsDeviceGot(session *xorm.Session, aidId int64, userId int64) (bool, error)
	IsSceneArrived(session *xorm.Session, aidId int64, userId int64) (bool, error)
	MarkDeviceGotWithSession(session *xorm.Session, aidId int64, userId int64) error
	MarkSceneArrivedWithSession(session *xorm.Session, aidId int64, userId int64) error
	GetAidDeviceGotTracksSorted(aidId int64) (models []*Model, err error)
}

const TableName = "aid_track"

type service struct{}

func (s service) GetUserTrackWithSession(session *xorm.Session, aidId int64, userId int64) (*Model, error) {
	var m Model
	_, err := session.Table(TableName).Where("help_info_id=? and user_id=?", aidId, userId).Get(&m)
	if err != nil {
		return nil, err
	}

	return &m, err
}

func (s service) IsDeviceGot(session *xorm.Session, aidId int64, userId int64) (bool, error) {
	track, err := s.GetUserTrackWithSession(session, aidId, userId)
	return track.DeviceGot, err
}

func (s service) IsSceneArrived(session *xorm.Session, aidId int64, userId int64) (bool, error) {
	track, err := s.GetUserTrackWithSession(session, aidId, userId)
	return track.SceneArrived, err
}

func (s service) GetUserTrack(aidId int64, userId int64) (*Model, error) {
	session := db.GetSession()
	defer session.Close()

	return s.GetUserTrackWithSession(session, aidId, userId)
}

func (s service) GetAidDeviceGotTracksSorted(aidId int64) (models []*Model, err error) {
	session := db.GetSession()
	defer session.Close()

	err = session.Table(TableName).Where("help_info_id = ? and device_got = 1", aidId).Asc("device_got_time").Find(&models)
	return
}

func (s service) CreateOrUpdateUserTrackWithSession(session *xorm.Session, model *Model) error {
	var exi Model
	exists, err := session.Table(TableName).Where("help_info_id=? and user_id=?", model.HelpInfoID, model.UserID).Get(&exi)
	if err != nil {
		return err
	}

	if exists {
		model.ID = exi.ID
		if _, err := session.Table(TableName).UseBool("device_got", "scene_arrived").ID(model.ID).Update(model); err != nil {
			return nil
		}
	} else {
		if _, err := session.Table(TableName).Insert(model); err != nil {
			return err
		}
	}

	return err
}

func (s service) MarkDeviceGotWithSession(session *xorm.Session, aidId int64, userId int64) error {
	now := time.Now()
	_, err := session.Exec(fmt.Sprintf("insert into `%s` (`help_info_id`,`user_id`,`device_got`,`device_got_time`) values(?, ?, ?, ?) on duplicate key update `device_got`=?, `device_got_time`=?", TableName), aidId, userId, true, now, true, now)
	return err
}

func (s service) MarkSceneArrivedWithSession(session *xorm.Session, aidId int64, userId int64) error {
	now := time.Now()
	_, err := session.Exec(fmt.Sprintf("insert into `%s` (`help_info_id`,`user_id`,`scene_arrived`,`scene_arrived_time`) values(?, ?, ?, ?) on duplicate key update `scene_arrived`=?, `scene_arrived_time`=?", TableName), aidId, userId, true, now, true, now)
	return err
}
