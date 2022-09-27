package merit_tree

import (
	"aed-api-server/internal/interfaces"
	"aed-api-server/internal/interfaces/entities"
	"aed-api-server/internal/interfaces/service"
	"aed-api-server/internal/pkg/db"
	"aed-api-server/internal/pkg/response"
	"aed-api-server/internal/pkg/utils"
	"github.com/go-xorm/xorm"
	"time"
)

type walk struct {
	User service.UserService `inject:"-"`
}

//go:inject-component
func NewWalkService() *walk {
	return &walk{}
}

type WalkConvert struct {
	Id        int64
	UserId    int64
	Walk      int
	UseWalk   int
	Points    int
	CreatedAt time.Time
}

func (w *walk) getTodayWalk(req *entities.WechatDataDecryptReq) (todayWalk int, err error) {
	walks, err := w.User.GetWalks(req)
	if err != nil {
		return 0, err
	}

	return walks.GetTodayWalks(), nil
}

type ConvertedInfo struct {
	ConvertWalk     int
	ConvertedPoints int
}

func getTotalConvertedInfo(userId int64) (*ConvertedInfo, error) {
	var info ConvertedInfo
	_, err := db.SQL(`
		select
			sum(use_walk) as convert_walk,
			sum(points) as converted_points
		from points_walk
		where
			user_id = ?
			and created_at >= CURRENT_DATE()
	`, userId).Get(&info)
	return &info, err
}

//GetWalkConvertInfo 获取积分兑换信息
func (w *walk) GetWalkConvertInfo(userId int64, req *entities.WechatDataDecryptReq) (*service.WalkConvertInfo, error) {
	rst, err := getTotalConvertedInfo(userId)
	if err != nil {
		return nil, err
	}

	if req.Code == "" {
		return &service.WalkConvertInfo{
			TodayWalk:       0,
			UnConvertWalk:   0,
			ConvertedPoints: rst.ConvertedPoints,
			ConvertRatio:    entities.WalkConvertRatio,
		}, nil
	}

	todayWalk, err := w.getTodayWalk(req)
	if err != nil {
		return nil, err
	}

	unConvertWalk := todayWalk - rst.ConvertWalk
	if unConvertWalk < 0 {
		unConvertWalk = 0
	}

	utils.Go(func() {
		interfaces.S.User.RecordUserEvent(userId, entities.UserEventTypeGetWalkStep, todayWalk)
	})

	return &service.WalkConvertInfo{
		TodayWalk:       todayWalk,
		UnConvertWalk:   unConvertWalk,
		ConvertedPoints: rst.ConvertedPoints,
		ConvertRatio:    entities.WalkConvertRatio,
	}, nil
}

//ConvertWalkToPoints 兑换积分
func (*walk) ConvertWalkToPoints(
	userId int64,
	todayWalk int,
) (*service.ConvertWalkToPointsRst, error) {
	rst, err := getTotalConvertedInfo(userId)
	if err != nil {
		return nil, err
	}

	useWalk := todayWalk - rst.ConvertWalk
	if useWalk < entities.WalkConvertRatio {
		return nil, response.ErrorNotEnoughWalkToConvent
	}

	eventWalk := interfaces.S.PointsScheduler.BuildPointsEventWalk(
		userId,
		todayWalk,
		rst.ConvertWalk,
		rst.ConvertedPoints,
	)

	var eventRst *entities.DealPointsEventRst
	err = db.Transaction(func(session *xorm.Session) error {
		eventRst, err = interfaces.S.PointsScheduler.DealPointsEvent(eventWalk)
		if err != nil {
			return err
		}

		useWalk = eventRst.PeckingPointsChange * entities.WalkConvertRatio

		//存储
		_, err = db.Insert("points_walk", WalkConvert{
			UserId:    userId,
			Walk:      todayWalk,
			UseWalk:   useWalk,
			Points:    eventRst.PeckingPointsChange,
			CreatedAt: time.Now(),
		})
		return err
	})

	if err != nil {
		return nil, err
	}
	return &service.ConvertWalkToPointsRst{
		UnConvertWalk:        todayWalk - useWalk - rst.ConvertWalk,
		ConvertedPoints:      rst.ConvertedPoints + eventRst.PeckingPointsChange,
		CurrentConvertPoints: eventRst.PeckingPointsChange,
		DealPointsRst:        eventRst,
	}, nil
}
