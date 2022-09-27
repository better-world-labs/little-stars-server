package aid

import (
	"aed-api-server/internal/interfaces"
	"aed-api-server/internal/interfaces/entities"
	"aed-api-server/internal/interfaces/events"
	service2 "aed-api-server/internal/interfaces/service"
	"aed-api-server/internal/pkg"
	"aed-api-server/internal/pkg/base"
	"aed-api-server/internal/pkg/db"
	"aed-api-server/internal/pkg/domain/emitter"
	"aed-api-server/internal/pkg/location"
	page "aed-api-server/internal/pkg/query"
	"aed-api-server/internal/pkg/response"
	"errors"
	"fmt"
	log "github.com/sirupsen/logrus"
	"time"
)

type aidService struct {
	UserService       service2.UserServiceOld `inject:"-"`
	User              service2.UserService    `inject:"-"`
	PracticeOrganizer int64                   `conf:"practice-organizer"`

	npcSelector *NPCSelector
}

// ArrivedEffectiveRange 到达现场的有效距离
const ArrivedEffectiveRange = 1000

//go:inject-component
func NewAidService() service2.AidService {
	return &aidService{npcSelector: newNPCSelector()}
}

func (s *aidService) Action120Called(aidID int64) error {
	return emitter.Emit(events.NewAidCalled(aidID))
}

func (s *aidService) ActionNPCArrived(aidId int64) error {
	helpInfo, exists, err := s.GetHelpInfoByID(aidId)
	if err != nil {
		return err
	}

	if !exists {
		return errors.New("help info not found")
	}

	if helpInfo.NpcId == nil {
		return errors.New("no npc found")
	}

	event := events.NewSceneArrivedEvent(aidId, *helpInfo.NpcId, 0)
	return interfaces.S.Activity.SaveActivitySceneArrived(event)
}

func (s *aidService) ActionArrived(userId int64, aidID int64, coordinate *location.Coordinate) ([]*entities.DealPointsEventRst, error) {
	helpInfo, exists, err := s.GetHelpInfoByID(aidID)
	if err != nil {
		return nil, err
	}

	if !exists {
		return nil, base.NewError("AidService", "help info not found")
	}

	t, err := GetService().GetUserTrack(aidID, userId)
	if err != nil {
		return nil, err
	}

	if t.SceneArrived {
		return nil, nil
	}

	distances := location.DistanceFrom(helpInfo.Coordinate, []location.Coordinate{*coordinate})
	log.Infof("distance for (%f,%f) and (%f,%f) is %d", helpInfo.Longitude, helpInfo.Latitude, coordinate.Longitude, coordinate.Latitude, distances[0])
	if distances[0] > ArrivedEffectiveRange {
		return nil, response.ErrorTooFar
	}
	event := events.NewSceneArrivedEvent(aidID, userId, 0)

	times, err := interfaces.S.Points.GetUserPointsEventTimes(userId, entities.PointsEventTypeArrived)
	if err != nil {
		return nil, err
	}

	rst := make([]*entities.DealPointsEventRst, 0)

	if !helpInfo.Exercise && times < pkg.UserPointsMaxTimesSceneArrived {
		eventRst, err := interfaces.S.PointsScheduler.DealPointsEvent(&events.PointsEvent{
			PointsEventType: entities.PointsEventTypeArrived,
			UserId:          userId,
			Params: map[string]interface{}{
				"uuid": event.Id,
			},
		})

		if err != nil {
			return nil, err
		}

		event.Points = eventRst.PeckingPointsChange
		rst = append(rst, eventRst)
	}

	//TODO 事务
	err = interfaces.S.Activity.SaveActivitySceneArrived(event)
	return rst, err
}

func (s *aidService) ActionCalled(userId int64, aidID int64) error {
	return emitter.Emit(events.NewSceneCalledEvent(aidID, userId))
}

func (s *aidService) ActionGoingToScene(userId int64, aidID int64) error {
	return emitter.Emit(events.NewGoingToSceneEvent(aidID, userId))
}

func (s *aidService) randomNPC() (*entities.SimpleUser, error) {
	phone := s.npcSelector.RandomPhone()
	u, exists, err := s.User.GetUserByPhone(phone)
	if err != nil {
		return nil, err
	}

	if !exists {
		return s.randomNPC()
	}

	return u, nil
}

func (s *aidService) PublishHelpInfoExercise(userId int64, dto *entities.PublishDTO) (id int64, npc *entities.SimpleUser, err error) {
	session := db.GetSession()
	defer session.Close()

	npc, err = s.randomNPC()
	if err != nil {
		return 0, nil, err
	}

	err = db.WithTransaction(session, func() error {
		helpInfo := &entities.HelpInfo{
			Coordinate: location.Coordinate{
				Longitude: dto.Longitude,
				Latitude:  dto.Latitude,
			},
			Images:        dto.Images,
			Address:       dto.Address,
			DetailAddress: dto.DetailAddress,
			Exercise:      true,
			Publisher:     userId,
			NpcId:         &npc.ID,
			PublishTime:   time.Now(),
		}
		_, err := session.Table("aid_help_info").Insert(helpInfo)
		if err != nil {
			return err
		}

		err = emitter.Emit(&events.HelpInfoPublishedEvent{
			HelpInfo: *helpInfo,
		})

		id = helpInfo.ID
		return err
	})

	return
}

func (s *aidService) PublishHelpInfo(userId int64, dto *entities.PublishDTO) (id int64, rst []*entities.DealPointsEventRst, err error) {
	session := db.GetSession()
	defer session.Close()

	practice := userId == s.PracticeOrganizer
	err = db.WithTransaction(session, func() error {
		helpInfo := &entities.HelpInfo{
			Coordinate: location.Coordinate{
				Longitude: dto.Longitude,
				Latitude:  dto.Latitude,
			},
			Images:        dto.Images,
			Address:       dto.Address,
			DetailAddress: dto.DetailAddress,
			Publisher:     userId,
			Practice:      practice,
			PublishTime:   time.Now(),
		}
		_, err := session.Table("aid_help_info").Insert(helpInfo)
		if err != nil {
			return err
		}

		times, err := interfaces.S.Points.GetUserPointsEventTimes(userId, entities.PointsEventTypePublishAid)
		if err != nil {
			return err
		}

		if times < pkg.UserPointsMaxTimesPublishHelpInfo {
			event, err := interfaces.S.PointsScheduler.DealPointsEvent(&events.PointsEvent{
				PointsEventType: entities.PointsEventTypePublishAid,
				UserId:          userId,
				Params: map[string]interface{}{
					"id": helpInfo.ID,
				},
			})
			if err != nil {
				return err
			}

			rst = append(rst, event)
		}

		err = emitter.Emit(&events.HelpInfoPublishedEvent{
			HelpInfo: *helpInfo,
		})

		id = helpInfo.ID
		return err
	})

	return
}

func (s *aidService) ComposeHelpInfoDTO(infos []*entities.HelpInfo, from *location.Coordinate) ([]*entities.HelpInfoComposedDTO, error) {
	helpInfoIDs := make([]int64, len(infos))
	locations := make([]location.Coordinate, len(infos))

	for i, e := range infos {
		helpInfoIDs[i] = e.ID
		locations[i] = e.Coordinate
	}

	latestActivityFuture := interfaces.S.Activity.ListMultiLatestCategorySortedAsync(helpInfoIDs, 1)

	latestActivitiesMap, err := latestActivityFuture()
	if err != nil {
		return nil, err
	}

	var distanceFrom []int64
	if from != nil && len(locations) > 0 {
		distanceFrom = location.DistanceFrom(*from, locations)
	}

	res := make([]*entities.HelpInfoComposedDTO, 0)
	for i, e := range infos {
		var dto entities.HelpInfoComposedDTO
		if from != nil {
			dto.Distance = distanceFrom[i]
		} else {
			dto.Distance = -1
		}

		s.FillHelpInfoSimpleDTO(e, &dto.HelpInfoDTO)
		activities := latestActivitiesMap[e.ID]
		s.FillDetailInfo(activities, &dto)

		res = append(res, &dto)
	}

	return res, nil
}

func (s aidService) FillTrackInfo(dto *entities.HelpInfoComposedDTO) {
	tracks, err := GetService().GetAidDeviceGotTracksSorted(dto.ID)
	if err != nil {
		return
	}

	if len(tracks) == 0 {
		return
	}

	firstAccount, err := s.UserService.GetUserByID(tracks[0].UserID)
	if err != nil {
		return
	}

	if firstAccount != nil {
		dto.FirstDeviceGetter = firstAccount.Nickname
	}

	dto.DeviceGetterCount = len(tracks)
}

func (s aidService) FillDetailInfo(activities []*entities.Activity, dto *entities.HelpInfoComposedDTO) {
	if len(activities) != 0 {
		latestActivity := activities[0]
		dto.NewestActivity = &entities.NewestActivityDTO{
			ID:      latestActivity.ID,
			Class:   latestActivity.Class,
			Record:  latestActivity.Record,
			Created: latestActivity.Created,
		}

		u := latestActivity.UserID
		if u != nil {
			us, err := s.UserService.GetUserByID(*u)
			if err == nil && us != nil {
				dto.NewestActivity.UserName = us.Nickname
			} else {
				log.Errorf("fill userName error: %v", err)
			}
		}
	}
}

func (s *aidService) ListHelpInfosInner24h() ([]*entities.HelpInfo, error) {
	session := db.GetSession()
	defer session.Close()

	arr := make([]*entities.HelpInfo, 0)
	t := time.Now().Add(-time.Hour * 24)
	fmt.Println(t)
	err := session.Table("aid_help_info").Where("publish_time >= ? and exercise = 0", t).Find(&arr)
	if err != nil {
		return nil, err
	}

	return arr, nil
}

func (s *aidService) ListHelpInfosPaged(pageQuery *page.Query, position *location.Coordinate, condition *entities.HelpInfo) (*page.Result[*entities.HelpInfoComposedDTO], error) {
	session := db.GetSession()
	defer session.Close()

	joined := session.Table("aid_help_info")

	if pageQuery.Page > 0 {
		joined.Limit(pageQuery.Size, (pageQuery.Page-1)*pageQuery.Size)
	}

	infos := make([]*entities.HelpInfo, 0)
	count, err := joined.UseBool("exercise").Desc("publish_time").FindAndCount(&infos, condition)

	if err != nil {
		return nil, err
	}

	dtos, err := s.ComposeHelpInfoDTO(infos, position)
	var r page.Result[*entities.HelpInfoComposedDTO]
	r.List = dtos
	r.Total = int(count)
	return &r, err
}

func (s *aidService) ListOneHoursInfos() ([]*entities.HelpInfo, error) {
	oneHoursAgo := time.Now().Add(-1 * time.Hour)
	infos := make([]*entities.HelpInfo, 0)

	err := db.Table("aid_help_info").
		Where("publish_time > ? and exercise = 0", oneHoursAgo).Desc("publish_time").Find(&infos)

	return infos, err
}

func (s *aidService) CountHelpInfosAboutMe(userId int64) (int64, error) {
	var count struct {
		Count int64
	}
	_, err := db.SQL(` select count(1) count from aid_help_info   
      where publisher = ?
      or exists(
        select distinct help_info_id from aid_activity where user_id = ?
    )`, userId, userId).Get(&count)

	return count.Count, err
}

func (s *aidService) ListHelpInfosParticipatedPaged(pageQuery *page.Query, userID int64) (*page.Result[*entities.HelpInfoComposedDTO], error) {
	session := db.GetSession()
	defer session.Close()

	joined := session.Table("aid_help_info").
		Join("LEFT", "aid_activity", "aid_help_info.id=aid_activity.help_info_id").
		Cols("aid_help_info.id, aid_help_info.longitude, aid_help_info.latitude, aid_help_info.address, aid_help_info.detail_address, aid_help_info.publisher, aid_help_info.publish_time")

	if pageQuery.Page > 0 {
		joined.Limit(pageQuery.Size, (pageQuery.Page-1)*pageQuery.Size)
	}

	infos := make([]*entities.HelpInfo, 0)
	count, err := joined.Where("aid_activity.user_id=?", userID).
		Distinct("aid_help_info.id").
		FindAndCount(&infos)

	if err != nil {
		return nil, err
	}

	dtos, err := s.ComposeHelpInfoDTO(infos, nil)
	var r page.Result[*entities.HelpInfoComposedDTO]
	r.List = dtos
	r.Total = int(count)
	return &r, err
}

func (s *aidService) GetHelpInfoByID(id int64) (*entities.HelpInfo, bool, error) {
	session := db.GetSession()
	defer session.Close()

	var helpInfo entities.HelpInfo
	exists, err := session.Table("aid_help_info").Where("id = ?", id).Get(&helpInfo)
	if err != nil {
		log.Errorf("get helpInfo error: %v", err)
		return nil, false, err
	}

	return &helpInfo, exists, nil
}

func (s *aidService) GetHelpInfoComposedByID(id int64, position *location.Coordinate) (*entities.HelpInfoComposedDTO, bool, error) {
	info, exists, err := s.GetHelpInfoByID(id)

	if err != nil {
		return nil, false, err
	}

	dto, err := s.ComposeHelpInfoDTO([]*entities.HelpInfo{info}, position)
	if err != nil {
		return nil, false, base.WrapError("Aid", "get help info error", err)
	}
	s.FillTrackInfo(dto[0])
	return dto[0], exists, err
}

func (s aidService) FillHelpInfoSimpleDTO(model *entities.HelpInfo, dto *entities.HelpInfoDTO) {
	if model != nil {
		dto.ID = model.ID
		dto.Exercise = model.Exercise
		dto.Practice = model.Practice
		dto.NpcId = model.NpcId
		dto.Longitude = model.Longitude
		dto.Latitude = model.Latitude
		dto.Address = model.Address
		dto.DetailAddress = model.DetailAddress
		dto.Images = model.Images
		dto.PublishTime = model.PublishTime.Format("2006-01-02 15:04:05")
		u, err := s.UserService.GetUserByID(model.Publisher)

		if dto.Images == nil {
			dto.Images = make([]string, 0)
		}

		if err == nil && u != nil {
			dto.PublisherName = u.Nickname
			dto.PublisherMobile = u.Mobile
		}
	}
}
