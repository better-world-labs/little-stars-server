package merit_tree

import (
	"aed-api-server/internal/interfaces/entities"
	"aed-api-server/internal/interfaces/events"
	"aed-api-server/internal/interfaces/service"
	"aed-api-server/internal/pkg/domain/emitter"
	"aed-api-server/internal/pkg/url_args"
	"aed-api-server/internal/pkg/utils"
	"errors"
	log "github.com/sirupsen/logrus"
	"time"
)

func NewTreasureChestService() service.TreasureChestService {
	return &treasureChestService{
		p: &persistence{},
	}
}

type treasureChestService struct {
	//持久层
	p iPersistence

	//宝箱过期时间周期
	Expired time.Duration `conf:"merit-tree.treasure-chest.expired,default=2h"`

	//宝箱冷却时间周期
	CoolDown time.Duration `conf:"merit-tree.treasure-chest.cooldown,default=24h"`
}

func (s *treasureChestService) CreateTreasureChest(request *entities.TreasureChestCreateRequest) error {
	maxSort, err := s.p.findTreasureChestMaxSort()
	if err != nil {
		return err
	}

	do := treasureChestDO{
		Name:      request.Name,
		Sort:      maxSort + 1,
		Link:      request.Link,
		LinkArgs:  url_args.GetArgs(request.Link),
		Tips:      request.Tips,
		Points:    request.Points,
		CreatedAt: time.Now(),
		TaskId:    request.TaskId,
	}

	return s.p.createTreasureChest(&do)
}

func (s *treasureChestService) OpenTreasureChest(userId int64, treasureChestId int) error {
	log.Infof("user=%v,treasureChestId=%v", userId, treasureChestId)
	chest, err := s.p.findChestByUserIdAndTreasureChestIdAndStatus(userId, treasureChestId, entities.TreasureChestStatusAvailable)
	if err != nil {
		return err
	}
	if chest == nil {
		return nil
	}

	//已经过期
	if time.Now().After(chest.ExpiredAt) {
		return s.p.updateChest(chest.Id, entities.TreasureChestStatusAvailable, treasureChestUserDto{Status: entities.TreasureChestStatusExpired})
	}

	chest.CompletedAt = time.Now()
	err = s.p.updateChest(chest.Id, entities.TreasureChestStatusAvailable, treasureChestUserDto{
		Status:      entities.TreasureChestStatusExpired,
		CompletedAt: chest.CompletedAt,
	})

	if err != nil {
		return err
	}

	s.sendUserOpenTreasureChestMsg(chest)
	return nil
}

func (s *treasureChestService) sendUserOpenTreasureChestMsg(userChest *treasureChestUserDto) {
	utils.Go(func() {
		chest, err := s.p.findTreasureChestById(userChest.TreasureChestId)
		if err != nil {
			log.Error("get chest from db err", err)
			return
		}

		if chest == nil {
			log.Errorf("chest(id=%v) do not exited", userChest.TreasureChestId)
			return
		}

		err = emitter.Emit(&events.UserOpenTreasureChest{
			UserId:            userChest.UserId,
			TreasureChestName: chest.Name,
			TaskId:            chest.TaskId,
			Points:            chest.Points,
			Link:              chest.Link,
			LinkArgs:          chest.LinkArgs,
			OpenTime:          userChest.CompletedAt,
		})
		if err != nil {
			log.Error(" emitter.Emit(&events.UserOpenTreasureChest{}) err", err)
		}
	})
}

func (s *treasureChestService) GetUserTreasureChest(userId int64) (*entities.TreasureChest, error) {
	chest, err := s.doGetUserTreasureChest(userId)
	if err != nil {
		return nil, err
	}
	if chest != nil {
		var now = time.Now()
		if chest.Status == entities.TreasureChestStatusInit && chest.ValidAt.After(now) {
			sub := chest.ValidAt.Sub(now)
			chest.ValidTtl = int(sub / time.Second)
		}
		if chest.Status == entities.TreasureChestStatusAvailable && chest.ExpiredAt.After(now) {
			sub := chest.ExpiredAt.Sub(now)
			chest.ExpiredTtl = int(sub / time.Second)
		}
	}
	return chest, nil
}
func (s *treasureChestService) doGetUserTreasureChest(userId int64) (*entities.TreasureChest, error) {
	//取两个未完成[时间上可能已经过期]
	chests, err := s.p.findLast2Chests(userId)
	if err != nil {
		return nil, err
	}
	if len(chests) == 2 {
		//业务上要求同时展示的宝箱只能有一个；那么如果第一个检查到状态已经过期，第二个一定没有过期
		chest, err := s.checkAndFixStatus(userId, chests[0], nil)
		if err != nil {
			return nil, err
		}
		if chest != nil {
			return chest, nil
		}
		return s.checkAndFixStatus(userId, chests[1], chests[0])
	}

	if len(chests) == 1 {
		return s.checkAndFixStatus(userId, chests[0], nil)
	}

	return nil, nil
}

func (s *treasureChestService) checkAndFixStatus(userId int64, dto *treasureChestDto, lastExpired *treasureChestDto) (*entities.TreasureChest, error) {
	//宝箱未分配给用户
	if 0 == dto.UserId {
		return s.assignChestToUser(userId, dto, lastExpired)
	}

	now := time.Now()

	//待展示
	if entities.TreasureChestStatusInit == dto.Status {
		chest := dto.TreasureChest
		//检查生效时间
		if now.After(dto.ValidAt) {
			expiredAt := time.Now().Add(s.Expired)
			err := s.p.updateChest(dto.UserTreasureChestId, entities.TreasureChestStatusInit, treasureChestUserDto{
				Status:    entities.TreasureChestStatusAvailable,
				ExpiredAt: expiredAt,
			})
			if err != nil {
				return nil, err
			}
			chest.Status = entities.TreasureChestStatusAvailable
			chest.ExpiredAt = expiredAt
		}
		return &chest, nil
	}

	//展示中
	if entities.TreasureChestStatusAvailable == dto.Status {
		//检查失效时间
		if now.After(dto.ExpiredAt) {
			err := s.p.updateChest(dto.UserTreasureChestId, entities.TreasureChestStatusAvailable, treasureChestUserDto{Status: entities.TreasureChestStatusExpired})
			if err != nil {
				return nil, err
			}
			return nil, nil
		}
		return &dto.TreasureChest, nil

	}
	return nil, errors.New("not supported status")
}

func (s *treasureChestService) assignChestToUser(userId int64, dto *treasureChestDto, lastExpired *treasureChestDto) (*entities.TreasureChest, error) {
	var err error
	if lastExpired == nil {
		lastExpired, err = s.p.getLastExpiredChest(userId)
		if err != nil {
			return nil, err
		}
	}

	var userChest *treasureChestUserDto
	if lastExpired.Id == 0 {
		userChest, err = s.assignFirstChest(userId, dto.Id)
	} else {
		userChest, err = s.assignChestToUserRefLastExpiredTime(userId, dto.Id, min(lastExpired.ExpiredAt, lastExpired.CompletedAt))
	}

	if err != nil {
		return nil, err
	}

	chest := dto.TreasureChest
	chest.ValidAt = userChest.ValidAt
	chest.Status = userChest.Status
	chest.ExpiredAt = userChest.ExpiredAt
	return &chest, nil
}

func (s *treasureChestService) assignChestToUserRefLastExpiredTime(userId int64, treasureChestId int, lastChestExpiredTime time.Time) (*treasureChestUserDto, error) {
	validAt := lastChestExpiredTime.Add(s.CoolDown)
	status := entities.TreasureChestStatusInit
	expiredAt := time.Time{}
	if time.Now().After(validAt) {
		status = entities.TreasureChestStatusAvailable
		expiredAt = time.Now().Add(s.Expired)
	}

	dto := treasureChestUserDto{
		UserId:          userId,
		TreasureChestId: treasureChestId,
		Status:          status,
		ValidAt:         validAt,
		ExpiredAt:       expiredAt,
		CreatedAt:       time.Now(),
		UpdatedAt:       time.Now(),
	}

	if err := s.p.insertTreasureChestUserDto(&dto); err != nil {
		return nil, err
	}
	return &dto, nil
}

func (s *treasureChestService) assignFirstChest(userId int64, treasureChestId int) (*treasureChestUserDto, error) {
	dto := treasureChestUserDto{
		UserId:          userId,
		TreasureChestId: treasureChestId,
		Status:          entities.TreasureChestStatusAvailable,
		ExpiredAt:       time.Now().Add(s.Expired),
		CreatedAt:       time.Now(),
		UpdatedAt:       time.Now(),
	}

	if err := s.p.insertTreasureChestUserDto(&dto); err != nil {
		return nil, err
	}
	return &dto, nil
}

//min 不是0的最小时间
func min(t1 time.Time, t2 time.Time) time.Time {
	zero := time.Time{}
	if t1 == zero {
		return t2
	}
	if t2 == zero {
		return t1
	}
	if t1.After(t2) {
		return t2
	} else {
		return t1
	}
}
