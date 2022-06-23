package merit_tree

import (
	"aed-api-server/internal/interfaces/entities"
	"aed-api-server/internal/pkg/db"
	"time"
)

type treasureChestDO struct {
	Id        int64
	Name      string
	Tips      string
	Link      string
	LinkArgs  []string
	Points    int
	TaskId    int64
	Sort      int
	CreatedAt time.Time
}

type treasureChestDto struct {
	entities.TreasureChest `xorm:"extends"`
	UserId                 int64
	UserTreasureChestId    int64
	CompletedAt            time.Time
}

type treasureChestUserDto struct {
	Id              int64
	TreasureChestId int
	Status          entities.TreasureChestStatus
	UserId          int64
	ValidAt         time.Time
	ExpiredAt       time.Time
	CompletedAt     time.Time
	CreatedAt       time.Time
	UpdatedAt       time.Time
}

type iPersistence interface {
	updateChest(userTreasureChestId int64, oldStatus entities.TreasureChestStatus, dto treasureChestUserDto) error
	findLast2Chests(userId int64) (list []*treasureChestDto, err error)
	findChestByUserIdAndTreasureChestIdAndStatus(
		userId int64,
		treasureChestId int,
		status entities.TreasureChestStatus,
	) (*treasureChestUserDto, error)

	insertTreasureChestUserDto(dto *treasureChestUserDto) error
	getLastExpiredChest(userId int64) (dto *treasureChestDto, err error)

	findTreasureChestById(id int) (*treasureChestDO, error)
	findTreasureChestMaxSort() (int, error)
	createTreasureChest(t *treasureChestDO) error
}
type persistence struct{}

func (*persistence) updateChest(userTreasureChestId int64, oldStatus entities.TreasureChestStatus, dto treasureChestUserDto) error {
	dto.UpdatedAt = time.Now()
	_, err := db.Table("treasure_chest_user").Where("id= ? and status=?", userTreasureChestId, oldStatus).Update(&dto)
	return err
}

//findLast2Chests 取两个未完成[时间上可能已经过期]
func (*persistence) findLast2Chests(userId int64) (list []*treasureChestDto, err error) {
	err = db.SQL(`
		select
			a.id,
			a.name,
			a.tips,
			a.link,
			a.points,
			b.status,
			b.user_id,
			b.id as user_treasure_chest_id,
			b.valid_at,
			b.expired_at
		from treasure_chest as a
		left join treasure_chest_user as b
			on b.user_id = ?
			and b.treasure_chest_id = a.id
		where
			(b.status is null or b.status < 20)
		order by sort asc, id asc
		limit 2
	`, userId).Find(&list)
	return list, err
}

func (*persistence) findChestByUserIdAndTreasureChestIdAndStatus(
	userId int64,
	treasureChestId int,
	status entities.TreasureChestStatus,
) (*treasureChestUserDto, error) {
	var chest treasureChestUserDto
	existed, err := db.Table("treasure_chest_user").Where("user_id = ? and treasure_chest_id=? and status = ?",
		userId, treasureChestId, status).Get(&chest)
	if err != nil {
		return nil, err
	}
	if existed {
		return &chest, nil
	}
	return nil, nil
}

func (*persistence) insertTreasureChestUserDto(dto *treasureChestUserDto) error {
	_, err := db.Table("treasure_chest_user").Insert(dto)
	return err
}

func (*persistence) getLastExpiredChest(userId int64) (dto *treasureChestDto, err error) {
	dto = new(treasureChestDto)
	_, err = db.SQL(`
		select
			treasure_chest_id as id,
			id as user_treasure_chest_id,
			status,
			user_id,
			valid_at,
			expired_at,
			completed_at
		from treasure_chest_user
		where
			user_id = ?
		order by expired_at desc
		limit 1
	`, userId).Get(dto)
	return dto, err
}

func (*persistence) findTreasureChestById(id int) (*treasureChestDO, error) {
	var chest treasureChestDO
	existed, err := db.Table("treasure_chest").Where("id = ?", id).Get(&chest)
	if err != nil {
		return nil, err
	}
	if existed {
		return &chest, nil
	}
	return nil, nil
}

func (*persistence) findTreasureChestMaxSort() (int, error) {
	count, err := db.SQL(`select ifnull(max(sort),0) from treasure_chest`).Count()
	return int(count), err
}

func (*persistence) createTreasureChest(t *treasureChestDO) error {
	_, err := db.Table("treasure_chest").Insert(t)
	return err
}
