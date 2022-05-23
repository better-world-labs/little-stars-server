package achievement

import (
	"aed-api-server/internal/interfaces"
	"aed-api-server/internal/pkg/db"
	"aed-api-server/internal/pkg/global"
	"github.com/go-xorm/xorm"
)

type UserMedal struct {
	ID         int64                `xorm:"id pk autoincr" json:"id,string"`
	MedalID    int64                `xorm:"medal_id" json:"medalId,string"`
	UserID     int64                `xorm:"user_id" json:"-"`
	BusinessId string               `xorm:"business_id" json:"-"`
	Created    global.FormattedTime `xorm:"created" json:"created"`
}

func ListUsersMedal(userID int64) ([]*UserMedal, error) {
	session := db.GetSession()
	defer session.Close()

	arr := make([]*UserMedal, 0)
	err := session.Table("user_medal").Where("user_id = ?", userID).Desc("created").Find(&arr)
	if err != nil {
		return nil, err
	}

	return arr, nil
}

func ListAll() ([]*UserMedal, error) {
	session := db.GetSession()
	defer session.Close()

	arr := make([]*UserMedal, 0)
	err := session.Table("user_medal").Desc("created").Find(&arr)
	if err != nil {
		return nil, err
	}

	return arr, nil
}

func GetUserMedal(medalId int64, userId int64) (*UserMedal, bool, error) {
	session := db.GetSession()
	defer session.Close()
	var u UserMedal
	exists, err := session.Table("user_medal").Where("user_id=? and medal_id=?", userId, medalId).Get(&u)

	if err != nil {
		return nil, false, nil
	}

	if !exists {
		return nil, false, nil
	}

	return &u, true, nil
}

func CreateUsersMedal(session *xorm.Session, medal *UserMedal) error {
	var exist UserMedal
	exists, err := session.Table("user_medal").Where("user_id=? and medal_id=?", medal.UserID, medal.MedalID).Get(&exist)
	if err != nil {
		return err
	}

	if !exists {
		_, err := session.Table("user_medal").Insert(medal)
		return err
	} else {
		medal.ID = exist.ID
		_, err := session.Table("user_medal").ID(medal.ID).Update(medal)
		return err
	}
}

type service struct{}

func init() {
	interfaces.S.UserMedal = &service{}
}

func (m *service) GetUserMedalUrl(userId int64) ([]string, error) {
	modalUrls := make([]string, 0)

	medalResource, err := ListMedals()
	if err != nil {
		return nil, err
	}
	list, err := ListUsersMedal(userId)
	if err != nil {
		return nil, err
	}
	for _, medal := range list {
		for _, r := range medalResource {
			if r.ID == medal.MedalID {
				modalUrls = append(modalUrls, r.ActiveIcon)
				break
			}
		}
	}
	return modalUrls, nil
}
