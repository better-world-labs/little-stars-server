package medal

import (
	"aed-api-server/internal/interfaces/entities"
	"aed-api-server/internal/interfaces/service"
	"aed-api-server/internal/pkg/db"
	"github.com/go-xorm/xorm"
)

func ListUsersMedal(userID int64) ([]*entities.UserMedal, error) {
	session := db.GetSession()
	defer session.Close()

	arr := make([]*entities.UserMedal, 0)
	err := session.Table("user_medal").Where("user_id = ?", userID).Desc("created").Find(&arr)
	if err != nil {
		return nil, err
	}

	return arr, nil
}

func ListAll() ([]*entities.UserMedal, error) {
	session := db.GetSession()
	defer session.Close()

	arr := make([]*entities.UserMedal, 0)
	err := session.Table("user_medal").Desc("created").Find(&arr)
	if err != nil {
		return nil, err
	}

	return arr, nil
}

func GetUserMedal(medalId int64, userId int64) (*entities.UserMedal, bool, error) {
	session := db.GetSession()
	defer session.Close()
	var u entities.UserMedal
	exists, err := session.Table("user_medal").Where("user_id=? and medal_id=?", userId, medalId).Get(&u)

	if err != nil {
		return nil, false, nil
	}

	if !exists {
		return nil, false, nil
	}

	return &u, true, nil
}

func CreateUsersMedal(session *xorm.Session, medal *entities.UserMedal) error {
	var exist entities.UserMedal
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

type userMedalImpl struct{}

//go:inject-component
func NewUserMedalService() service.IUserMedal {
	return &userMedalImpl{}
}

func (m *userMedalImpl) GetUserMedalUrl(userId int64) ([]string, error) {
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
