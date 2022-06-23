package project

import (
	"aed-api-server/internal/interfaces/entities"
	"aed-api-server/internal/interfaces/events"
	"aed-api-server/internal/interfaces/service"
	"aed-api-server/internal/pkg/db"
	"aed-api-server/internal/pkg/domain/emitter"
	"errors"
	"fmt"
	"time"
)

const (
	projectTableName     = "project_project"
	projectUserTableName = "project_user"

	Level0 = "青铜"
	Level1 = "白银"
	Level2 = "黄金"
	Level3 = "铂金"
	Level4 = "钻石"
	Level5 = "星耀"
)

func NewProjectService() *Service {
	return &Service{}
}

type Service struct {
}

type UserDO struct {
	Id               int64 `xorm:"id pk autoincr"`
	UserId           int64
	ProjectId        int64
	Level            int
	VideoCompleted   bool
	VideoCompletedAt time.Time
	Certification    bool
	CertificationAt  time.Time
	CreatedAt        time.Time
}

func (Service) GetProjectById(projectId int64) (*service.Project, error) {
	var project service.Project

	existed, err := db.GetById(projectTableName, projectId, &project)
	project.CertPoint = 200
	//TODO 补充到积分事件定义中

	if err != nil {
		return nil, err
	}
	if !existed {
		return nil, errors.New("project is not found")
	}
	return &project, nil
}

func findUserDO(projectId int64, userId int64) (*UserDO, error) {
	var user UserDO
	existed, err := db.Table(projectUserTableName).Where("project_id = ? and user_id = ?", projectId, userId).Get(&user)
	if err != nil {
		return nil, err
	}
	if existed {
		return &user, nil
	}
	return nil, nil
}

func (Service) IsProjectVideoCompleted(projectId int64, userId int64) (bool, error) {
	user, err := findUserDO(projectId, userId)
	if err != nil {
		return false, err
	}
	if user != nil && user.VideoCompleted {
		return true, nil
	}
	return false, nil
}

func (s *Service) CompletedProjectVideo(projectId int64, userId int64) (*entities.PointAddRst, error) {
	proj, err := s.GetProjectById(projectId)
	if err != nil {
		return nil, err
	}
	if proj == nil {
		return nil, errors.New("not found the project")
	}
	const PointAddDescription = "视频学习加积分"
	rst := entities.PointAddRst{
		Point:       proj.Point,
		Description: PointAddDescription,
	}

	user, err := findUserDO(projectId, userId)
	if err != nil {
		return nil, err
	}
	if user == nil {
		user := UserDO{
			UserId:           userId,
			ProjectId:        projectId,
			Level:            0,
			VideoCompleted:   true,
			VideoCompletedAt: time.Now(),
			Certification:    false,
			CreatedAt:        time.Now(),
		}
		_, err := db.Insert(projectUserTableName, &user)
		if err != nil {
			return nil, err
		}
		err = emitter.Emit(&events.PointsEvent{
			PointsEventType: entities.PointsEventTypeLearntVideo,
			UserId:          userId,
			Params: entities.PointsEventParams{
				RefTable:   projectUserTableName,
				RefTableId: user.Id,
			},
		})

		if err != nil {
			return nil, err
		}
		return &rst, err
	} else if !user.VideoCompleted {
		_, err := db.Exec(`
			update `+projectUserTableName+`
			set video_completed = 1, video_completed_at = now()
			where
				project_id = ?
				and user_id = ?
		`, projectId, userId)
		if err != nil {
			return nil, err
		}

		err = emitter.Emit(&events.PointsEvent{
			PointsEventType: entities.PointsEventTypeLearntVideo,
			UserId:          userId,
			Params: entities.PointsEventParams{
				RefTable:   projectUserTableName,
				RefTableId: user.Id,
			},
		})
		if err != nil {
			return nil, err
		}
		return &rst, err
	}
	return nil, nil
}

func (Service) GetUserProjectLevel(projectId int64, userId int64) (*service.ProjectLevel, error) {
	user, err := findUserDO(projectId, userId)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return userLevel(0, false), err
	}
	return userLevel(user.Level, user.Certification), err
}

func (Service) DoCertification(projectId int64, userId int64) error {
	sql := fmt.Sprintf(`insert into %s (user_id, project_id, level, video_completed, certification , certification_at, created_at)
		values (?, ?, 0, 0, ?, ?, ?) 
        ON DUPLICATE KEY UPDATE 
        certification=?, certification_at=?
    `, projectUserTableName)
	_, err := db.Exec(sql, userId, projectId, true, time.Now(), time.Now(),
		true, time.Now())

	return err
}

func (Service) UpdateUserProjectLevel(projectId int64, userId int64, level int) error {
	_, err := db.Exec(`
		insert into project_user(user_id, project_id, level, video_completed, certification, created_at)
		values(?, ?, ?, 0, 0, now())
		ON DUPLICATE KEY UPDATE
		level = ?
	`, userId, projectId, level, level)
	return err
}

func (Service) IsNotLearnt(userId int64) bool {
	type Result struct {
		Learnt bool
	}
	var result Result
	_, _ = db.SQL(`
		select
		EXISTS(select 1 from project_user where user_id = ?)
		or  EXISTS(select 1 from project_course_user where user_id = ?)
		as learnt
	`, userId, userId).Get(&result)
	return !result.Learnt
}

func userLevel(level int, certification bool) *service.ProjectLevel {
	name := Level0
	switch level {
	case 0:
		name = Level0
	case 1:
		name = Level1
	case 2:
		name = Level2
	case 3:
		name = Level3
	case 4:
		name = Level4
	case 5:
		name = Level5
	}
	return &service.ProjectLevel{
		Level:         level,
		LevelName:     name,
		Certification: certification,
	}
}
