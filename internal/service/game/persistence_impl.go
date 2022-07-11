package game

import (
	"aed-api-server/internal/interfaces/entities"
	"aed-api-server/internal/pkg/db"
	"github.com/go-xorm/xorm"
)

type GamePersistenceImpl struct {
	table string
}

func NewGamePersistence() IGamePersistence {
	return &GamePersistenceImpl{table: "game"}
}

func (p *GamePersistenceImpl) GetById(id int64) (g *entities.Game, exists bool, err error) {
	g = &entities.Game{}
	exists, err = p.tableSession().Where("id = ?", id).Get(g)
	return
}

func (p *GamePersistenceImpl) ListGames() (games []*entities.Game, err error) {
	err = p.tableSession().
		Asc("start_at").
		Find(&games)
	return
}

func (p *GamePersistenceImpl) ListByIds(ids []int64) (games []*entities.Game, err error) {
	err = p.tableSession().In("id", ids).Find(&games)
	return
}

func (p *GamePersistenceImpl) Create(game *entities.Game) error {
	err := db.Transaction(func(session *xorm.Session) error {
		_, err := session.Insert(game)
		return err
	})

	return err
}

func (p *GamePersistenceImpl) Update(game *entities.Game) error {
	err := db.Transaction(func(session *xorm.Session) error {
		_, err := session.ID(game.Id).UseBool("settled").Update(game)
		return err
	})

	return err
}

func (p *GamePersistenceImpl) Delete(id int64) error {
	err := db.Transaction(func(session *xorm.Session) error {
		_, err := session.Where("id = ?", id).Delete(&entities.Game{})
		return err
	})

	return err
}

func (p *GamePersistenceImpl) tableSession() *xorm.Session {
	return db.Table(p.table)
}
