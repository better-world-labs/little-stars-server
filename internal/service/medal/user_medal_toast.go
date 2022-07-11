package medal

import (
	"aed-api-server/internal/interfaces/entities"
	"aed-api-server/internal/pkg/db"
	"github.com/go-xorm/xorm"
)

func ListUsersMedalToast(userID int64) ([]*entities.UserMedal, error) {
	session := db.GetSession()
	defer session.Close()

	arr := make([]*entities.UserMedal, 0)
	//err := session.Table("user_medal_toast").Where("user_id = ? and `read` = 0", userID).Desc("created").Find(&arr)
	err := session.Table("user_medal").
		Join("LEFT", "user_medal_toast", "user_medal.user_id = user_medal_toast.user_id and user_medal.medal_id = user_medal_toast.medal_id").
		Cols("user_medal.id", "user_medal.user_id", "user_medal.medal_id", "user_medal.created").
		Where("user_medal.user_id = ? and `read` = 0", userID).Desc("created").Find(&arr)

	if err != nil {
		return nil, err
	}

	_, err = session.Exec("update user_medal_toast set `read` = 1 where user_id = ?", userID)

	if err != nil {
		return nil, err
	}

	return arr, nil
}

func CreateUsersMedalToast(session *xorm.Session, medal *entities.UserMedal) error {
	exists, err := session.Table("user_medal_toast").Where("user_id=? and medal_id=?", medal.UserID, medal.MedalID).Exist()
	if err != nil {
		return err
	}

	if !exists {
		_, err = session.Table("user_medal_toast").Insert(medal)
	}

	return err
}
