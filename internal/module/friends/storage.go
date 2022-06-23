package friends

import (
	"aed-api-server/internal/interfaces"
	"aed-api-server/internal/interfaces/entities"
	"aed-api-server/internal/interfaces/events"
	"aed-api-server/internal/pkg"
	"aed-api-server/internal/pkg/db"
	"aed-api-server/internal/pkg/domain/emitter"
	"github.com/go-xorm/xorm"
	"time"
)

const (
	TableNameRelationship = "friends_relationship"
	TableNameProspective  = "friends_prospective"
)

func createRelationship(relationship *Relationship) error {
	var current Relationship
	exists, err := db.Table(TableNameRelationship).
		Where("user_id = ? ", relationship.UserId).
		Exist(&current)
	if err != nil {
		return err
	}

	if !exists {
		_, err := db.Table(TableNameRelationship).Insert(relationship)
		return err
	}

	//TODO Loop Check
	return nil
}

func listByParentID(parentID int64) (relationships []*Relationship, err error) {
	err = db.Table(TableNameRelationship).Where("parent_id = ?", parentID).Find(&relationships)
	return
}

func listAll() ([]*Relationship, error) {
	var slice []*Relationship
	err := db.Table(TableNameRelationship).Find(&slice)
	return slice, err
}

func hasFriend(userId int64) (bool, error) {
	return db.Table(TableNameRelationship).Where("parent_id = ?", userId).Exist()
}

func createProspective(prospective *Prospective) error {
	_, exists, err := getProspective(prospective.OpenId)
	if err != nil {
		return err
	}

	if !exists {
		return db.Begin(func(session *xorm.Session) error {
			prospective.CreateAt = time.Now()
			_, err := db.Table(TableNameProspective).Insert(prospective)
			if err != nil {
				return err
			}

			times, err := interfaces.S.Points.GetUserPointsEventTimes(prospective.ParentId, entities.PointsEventTypeInvite)
			if err != nil {
				return err
			}

			if times < pkg.UserPointsMaxTimesInvite {
				err = emitter.Emit(&events.PointsEvent{
					PointsEventType: entities.PointsEventTypeInvite,
					UserId:          prospective.ParentId,
					Params: map[string]interface{}{
						"beInvitedOpenId": prospective.OpenId,
					},
				})
			}

			return err
		})
	}

	return nil
}

func getProspective(openid string) (*Prospective, bool, error) {
	var p Prospective
	exists, err := db.Table(TableNameProspective).Where("open_id = ?", openid).Get(&p)
	return &p, exists, err
}
