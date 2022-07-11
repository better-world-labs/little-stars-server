package game

import (
	"aed-api-server/internal/interfaces/entities"
	"aed-api-server/internal/pkg/db"
	"github.com/go-xorm/xorm"
)

type PlayerPersistenceImpl struct {
	table string
}

func NewPlayerPersistence() IPlayerPersistence {
	return &PlayerPersistenceImpl{"game_process"}
}

func (p *PlayerPersistenceImpl) ListUserProcessesCompleted(gameId int64) ([]*entities.GameProcess, error) {
	var g []*entities.GameProcess
	err := db.Table(p.table).Where("game_id = ? and completed = true", gameId).Find(&g)
	return g, err
}

func (p *PlayerPersistenceImpl) ListTopUserProcessesCompleted(gameId int64, top int) ([]*entities.GameProcess, error) {
	var g []*entities.GameProcess
	err := db.Table(p.table).Where("game_id = ? and completed = true", gameId).
		Asc("updated_at").
		Limit(top, 0).
		Find(&g)
	return g, err
}

func (p *PlayerPersistenceImpl) ListUserProcesses(userId int64) ([]*entities.GameProcess, error) {
	var g []*entities.GameProcess
	err := db.Table(p.table).Where("user_id = ?", userId).Find(&g)
	return g, err
}

func (p *PlayerPersistenceImpl) ListUserProcessesUncompleted(userId int64) ([]*entities.GameProcess, error) {
	var g []*entities.GameProcess
	err := db.Table(p.table).Where("user_id = ? and completed = false", userId).Find(&g)
	return g, err
}

func (p *PlayerPersistenceImpl) GetById(id int64) (*entities.GameProcess, bool, error) {
	var g entities.GameProcess
	exists, err := db.Table(p.table).Where("id = ?", id).Get(&g)
	return &g, exists, err
}

func (p *PlayerPersistenceImpl) CountUserCompleted(gameId int64) (int64, error) {
	return p.tableSession().Where("game_id = ?", gameId).
		Where("game_id = ? and completed = true", gameId).
		Distinct("user_id").
		Cols("user_id").
		Count()
}

func (p *PlayerPersistenceImpl) ListRankUserIds(gameId int64) ([]int64, error) {
	top, err := p.ListTopBySteps(gameId, 0)
	if err != nil {
		return nil, err
	}

	var res []int64
	for _, u := range top {
		res = append(res, u.UserId)
	}

	return res, nil
}

func (p *PlayerPersistenceImpl) ListTopBySteps(gameId int64, top int) (games []*entities.GameProcess, err error) {
	session := p.tableSession().
		Where("game_id = ?", gameId).
		OrderBy("(history_steps + active_steps) desc")

	if top > 0 {
		session.Limit(top, 0)
	}

	err = session.
		Find(&games)

	return
}

func (p *PlayerPersistenceImpl) GetByGameIdAndUserId(gameId, userId int64) (*entities.GameProcess, bool, error) {
	var gp entities.GameProcess
	exists, err := p.tableSession().Where("user_id = ? and game_id = ?", userId, gameId).Get(&gp)
	return &gp, exists, err
}

func (p *PlayerPersistenceImpl) ListJoinedGameIds(userId int64) ([]int64, error) {
	var ids []int64
	var gameIds []*struct{ GameId int64 }

	err := p.tableSession().
		Cols("game_id").
		Where("user_id = ?", userId).
		Distinct("game_id").Find(&gameIds)

	if err != nil {
		return nil, err
	}

	for _, g := range gameIds {
		ids = append(ids, g.GameId)
	}

	return ids, nil
}

func (p *PlayerPersistenceImpl) Create(player *entities.GameProcess) error {
	return db.Transaction(func(session *xorm.Session) error {
		_, err := session.Table(p.table).Insert(player)
		return err
	})
}

func (p *PlayerPersistenceImpl) CompareAndSet(excepted entities.GameProcess, player *entities.GameProcess) (bool, error) {
	var updated int64
	var err error

	err = db.Transaction(func(session *xorm.Session) error {
		updated, err = session.Table(p.table).
			Where("id = ? and history_steps = ? and active_steps = ? and history_steps_date = ? and completed = ?",
				player.Id, excepted.HistorySteps,
				excepted.ActiveSteps,
				excepted.HistoryStepsDate,
				excepted.Completed).
			MustCols("active_steps", "history_steps", "completed").
			Update(player)
		return err
	})

	return updated > 0, err
}

func (p *PlayerPersistenceImpl) Update(player *entities.GameProcess) error {
	return db.Transaction(func(session *xorm.Session) error {
		_, err := session.Table(p.table).ID(player.Id).UseBool("completed").Update(player)
		return err
	})
}

func (p *PlayerPersistenceImpl) Delete(id int64) error {
	err := db.Transaction(func(session *xorm.Session) error {
		_, err := session.Table(p.table).Delete(&entities.GameProcess{Id: id})
		return err
	})

	return err
}

func (p *PlayerPersistenceImpl) tableSession() *xorm.Session {
	return db.Table(p.table)
}
