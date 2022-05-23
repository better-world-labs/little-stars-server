package service

import "aed-api-server/internal/interfaces/entities"

type Project struct {
	Id              int    `json:"id"`
	Name            string `json:"name"`
	VideoUrl        string `json:"videoUrl"`
	VideoName       string `json:"videoName"`
	VideoEndTip     string `json:"videoEndTip"`
	Point           int    `json:"point"`
	CertPoint       int    `json:"certPoint"       xorm:"-"`
	CertDescription string `json:"certDescription"`
}

type ProjectLevel struct {
	Level         int    `json:"level"`
	LevelName     string `json:"levelName"`
	Certification bool   `json:"certification"`
}

type ProjectService interface {
	GetProjectById(projectId int64) (*Project, error)
	IsProjectVideoCompleted(projectId int64, userId int64) (bool, error)
	CompletedProjectVideo(projectId int64, userId int64) (*entities.PointAddRst, error)
	GetUserProjectLevel(projectId int64, userId int64) (*ProjectLevel, error)

	// DoCertification 记录认证
	DoCertification(projectId int64, userId int64) error

	UpdateUserProjectLevel(projectId int64, userId int64, level int) error

	IsNotLearnt(userId int64) bool
}
