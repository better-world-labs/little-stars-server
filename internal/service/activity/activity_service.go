package activity

import (
	"aed-api-server/internal/interfaces"
	"aed-api-server/internal/interfaces/entities"
	"aed-api-server/internal/interfaces/events"
	"aed-api-server/internal/interfaces/service"
	"aed-api-server/internal/pkg"
	"aed-api-server/internal/pkg/async"
	"aed-api-server/internal/pkg/base"
	"aed-api-server/internal/pkg/db"
	"aed-api-server/internal/service/aid"
	"errors"
	"github.com/go-xorm/xorm"
	"github.com/google/uuid"
	log "github.com/sirupsen/logrus"
)

type activityService struct {
	Aid service.AidService `inject:"-"`
}

const (
	TableName = "aid_activity"
	Module    = "activity"
)

//go:inject-component
func NewActivityService() service.ActivityService {
	return &activityService{}
}

func (s *activityService) Create(a *entities.Activity) error {
	if a.Uuid == "" {
		a.Uuid = uuid.NewString()
	}

	session := db.GetSession()
	defer session.Close()

	return s.CreateWithSession(session, a)
}

func (s *activityService) CreateOrUpdateByUUID(a *entities.Activity) error {
	o, exists, err := s.GetOneByUUID(a.Uuid)
	if err != nil {
		return base.WrapError(Module, "CreateOrUpdateByUUID error", err)
	}

	if !exists {
		err = s.Create(a)
	} else {
		a.ID = o.ID
		err = s.UpdateByID(a)
	}

	return err
}
func (s *activityService) ListByAID(aid int64, limit int) ([]*entities.Activity, error) {
	session := db.GetSession()
	defer session.Close()

	activities := make([]*entities.Activity, 0)
	sess := session.Table(TableName).Desc("created").Where("help_info_id = ?", aid)

	if limit > 0 {
		sess = sess.Limit(limit, 0)
	}

	err := sess.Find(&activities)

	if err != nil {
		log.Errorf("list error: %v", err)
		return nil, err
	}

	return activities, nil
}

func (s *activityService) ListByAIDs(aids []int64) (map[int64][]*entities.Activity, error) {
	session := db.GetSession()
	defer session.Close()

	activities := make([]*entities.Activity, 0)
	err := session.Table(TableName).
		In("help_info_id", aids).
		Desc("created").
		Find(&activities)

	if err != nil {
		log.Errorf("list by id error: %v", err)
		return nil, err
	}

	m := make(map[int64][]*entities.Activity, 0)
	for i := range activities {
		arr := m[activities[i].HelpInfoID]
		arr = append(arr, activities[i])
		m[activities[i].HelpInfoID] = arr
	}

	return m, nil
}

func (s *activityService) ListMultiLatestCategorySortedAsync(aids []int64, latest int) func() (map[int64][]*entities.Activity, error) {
	return async.RunAsync(func() (map[int64][]*entities.Activity, error) {
		return s.ListMultiLatestCategorySorted(aids, latest)
	})
}

func (s *activityService) ListMultiLatestCategorySorted(aids []int64, latest int) (map[int64][]*entities.Activity, error) {
	list, err := s.ListByAIDs(aids)
	if err != nil {
		return nil, err
	}

	for k, v := range list {
		list[k] = s.ParseLatestCategorySorted(v, latest)
	}

	return list, nil
}

func (s *activityService) ListCategorySorted(aid int64) ([]*entities.Activity, error) {
	return s.ListLatestCategorySorted(aid, 0)
}

func (s *activityService) ListLatestCategorySorted(aid int64, latest int) ([]*entities.Activity, error) {
	list, err := s.ListByAID(aid, latest)
	if err != nil {
		return nil, err
	}

	return s.ParseLatestCategorySorted(list, latest), nil
}

func (s *activityService) ParseLatestCategorySorted(list []*entities.Activity, latest int) []*entities.Activity {
	if latest == 0 {
		latest = len(list)
	}
	for i := range list {
		a := list[i]
		if a.Category == "" {
			PutCategory(a)
		}
	}

	return list
}

func (s *activityService) ListLatestCategorySortedAsync(aid int64, latest int) func() ([]*entities.Activity, error) {
	return async.RunAsync(func() ([]*entities.Activity, error) {
		return s.ListLatestCategorySorted(aid, latest)
	})
}

func (s *activityService) GetOneByID(id int64) (*entities.Activity, error) {
	session := db.GetSession()
	defer session.Close()

	var a entities.Activity
	_, err := session.Table(TableName).Where("id = ?", id).Get(&a)

	if err != nil {
		log.Errorf("get activity by id error: %v", err)
		return nil, err
	}

	PutCategory(&a)
	return &a, nil
}

func (s *activityService) GetOneByUUID(uuid string) (*entities.Activity, bool, error) {
	session := db.GetSession()
	defer session.Close()

	var a entities.Activity
	exists, err := session.Table(TableName).Where("uuid = ?", uuid).Get(&a)
	if err != nil {
		log.Errorf("get activity by id error: %v", err)
		return nil, exists, err
	}

	PutCategory(&a)
	return &a, exists, nil
}

func (s *activityService) GetManyByIDs(ids []int64) ([]*entities.Activity, error) {
	session := db.GetSession()
	defer session.Close()

	arr := make([]*entities.Activity, 0)
	err := session.Table(TableName).In("id", ids).Find(&arr)
	if err != nil {
		log.Errorf("list activity by ids error: %v", err)
		return nil, err
	}

	for i := range arr {
		PutCategory(arr[i])
	}

	return arr, nil
}

func (s *activityService) ListByAidAndClass(aid int64, class string) ([]*entities.Activity, error) {
	session := db.GetSession()
	defer session.Close()

	arr := make([]*entities.Activity, 0)
	err := session.Table(TableName).Where("help_info_id= ? and class = ?", aid, class).Find(&arr)
	if err != nil {
		log.Errorf("list activity by ids error: %v", err)
		return nil, err
	}

	return arr, nil

}
func (s *activityService) GetLastUpdated(aid int64) (*entities.Activity, error) {
	session := db.GetSession()
	defer session.Close()

	var a entities.Activity
	exists, err := session.Table(TableName).Where("help_info_id=?", aid).Desc("created").Limit(1, 0).Get(&a)
	if err != nil {
		log.Errorf("list activity by ids error: %v", err)
		return nil, err
	}

	if !exists {
		return nil, nil
	}

	return &a, nil
}

func (s *activityService) CreateWithSession(session *xorm.Session, a *entities.Activity) error {
	_, err := session.Table(TableName).Insert(a)
	if err != nil {
		log.Errorf("insert error: %v", err)
		return err
	}

	return nil
}

func (s *activityService) UpdateByID(a *entities.Activity) error {
	session := db.GetSession()
	defer session.Close()

	_, err := session.Table(TableName).ID(a.ID).Update(a)
	return err
}

func (s *activityService) SaveActivityGoingToScene(event *events.GoingToSceneEvent) error {
	return s.CreateOrUpdateByUUID(CreateActivityGoingToScene(event))
}

func (s *activityService) SaveActivityNPCDeviceGot(event *events.DeviceGotEvent) error {
	return db.Begin(func(session *xorm.Session) error {
		a := CreateActivityDeviceGot(event)

		err := s.CreateWithSession(session, a)
		if err != nil {
			log.Errorf("onDeviceGot error: %v", err)
			return err
		}

		err = aid.GetService().MarkDeviceGotWithSession(session, event.Aid, event.UserId)
		if err != nil {
			return err
		}

		return nil
	})
}
func (s *activityService) SaveActivityDeviceGot(event *events.DeviceGotEvent) ([]*entities.DealPointsEventRst, error) {
	var rst []*entities.DealPointsEventRst
	helpInfo, exists, err := s.Aid.GetHelpInfoByID(event.Aid)
	if err != nil {
		return nil, err
	}

	if !exists {
		return nil, errors.New("help info not found")
	}

	return rst, db.Begin(func(session *xorm.Session) error {
		a := CreateActivityDeviceGot(event)
		isDeviceGot, err := aid.GetService().IsDeviceGot(session, event.Aid, event.UserId)
		if err != nil {
			return err
		}

		if err != nil {
			log.Errorf("onDeviceGot error: %v", err)
			return err
		}

		err = s.CreateWithSession(session, a)
		if err != nil {
			log.Errorf("onDeviceGot error: %v", err)
			return err
		}

		err = aid.GetService().MarkDeviceGotWithSession(session, event.Aid, event.UserId)
		if err != nil {
			return err
		}

		times, err := interfaces.S.Points.GetUserPointsEventTimes(event.UserId, entities.PointsEventTypeGotDevice)
		if err != nil {
			return err
		}

		if !helpInfo.Exercise && !isDeviceGot && times < pkg.UserPointsMaxTimesGetDevice { //领取过不再加积分
			eventRst, err := interfaces.S.PointsScheduler.DealPointsEvent(&events.PointsEvent{
				PointsEventType: entities.PointsEventTypeGotDevice,
				UserId:          event.UserId,
				Params: map[string]interface{}{
					"uuid": event.Id,
				},
			})

			if err != nil {
				return err
			}

			event.Points = eventRst.PeckingPointsChange
			rst = append(rst, eventRst)
		}

		return nil
	})
}

func (s *activityService) SaveActivitySceneArrived(event *events.SceneArrivedEvent) error {
	return db.Begin(func(session *xorm.Session) error {
		a := CreateActivitySceneArrived(event)
		if err := s.CreateWithSession(session, a); err != nil {
			return err
		}

		if err := aid.GetService().MarkSceneArrivedWithSession(session, event.Aid, event.UserId); err != nil {
			return err
		}

		return nil
	})
}

func (s *activityService) SaveActivityGoingToGetDevice(event *events.GoingToGetDeviceEvent) error {
	a := CreateActivityGoingToGetDevice(event)
	err := s.CreateOrUpdateByUUID(a)
	if err != nil {
		return err
	}

	return nil
}

func (s *activityService) SaveActivityVolunteerNotified(event *events.VolunteerNotifiedEvent) error {
	a := CreateActivityVolunteerNotified(event)
	err := s.CreateOrUpdateByUUID(a)
	if err != nil {
		return err
	}

	return nil
}

func (s *activityService) SaveActivityAidCalled(event *events.AidCalledEvent) error {
	a := CreateActivityAidCalled(event)
	activities, err := s.ListByAidAndClass(a.HelpInfoID, a.Class)
	if err != nil {
		return err
	}

	if len(activities) > 0 {
		return nil
	}

	err = s.CreateOrUpdateByUUID(a)
	if err != nil {
		return err
	}

	return nil
}

func (s *activityService) SaveActivitySceneCalled(event *events.SceneCalledEvent) error {
	a := CreateActivitySceneCalled(event)
	err := s.CreateOrUpdateByUUID(a)
	if err != nil {
		return err
	}

	return nil
}

func (s *activityService) SaveActivitySceneReport(event *events.SceneReportEvent) ([]*entities.DealPointsEventRst, error) {
	var rst []*entities.DealPointsEventRst

	helpInfo, exists, err := s.Aid.GetHelpInfoByID(event.Aid)
	if err != nil {
		return nil, err
	}

	if !exists {
		return nil, errors.New("not found")
	}

	times, err := interfaces.S.Points.GetUserPointsEventTimes(event.UserId, entities.PointsEventTypeUploadAidInfo)
	if err != nil {
		return nil, err
	}

	log.Info("SaveActivitySceneReport: user has do ", times, "times")
	if !helpInfo.Exercise && times < pkg.UserPointsMaxTimesUploadScene {
		eventRst, err := interfaces.S.PointsScheduler.DealPointsEvent(&events.PointsEvent{
			PointsEventType: entities.PointsEventTypeUploadAidInfo,
			UserId:          event.UserId,
			Params: map[string]interface{}{
				"uuid": event.Id,
			},
		})
		if err != nil {
			return nil, err
		}

		rst = append(rst, eventRst)
		event.Points = eventRst.PeckingPointsChange
	}

	activity, err := CreateActivitySceneReport(event)
	if err != nil {
		return nil, err
	}

	return rst, s.Create(activity)
}
