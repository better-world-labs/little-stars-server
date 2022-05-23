package service

import "aed-api-server/internal/interfaces/entities"

type ContentItem struct {
	Type  string `json:"t"`
	Value string `json:"v"`
}

type Section struct {
	Title   string         `json:"title"`
	Content []*ContentItem `json:"content"`
}

type Article struct {
	Id       int64      `json:"id"`
	Title    string     `json:"title"`
	Abstract *Section   `json:"abstract"`
	Sections []*Section `json:"sections"`
}

type Course struct {
	Id      int      `json:"id"`
	Name    string   `json:"name"`
	TabName string   `json:"tabName"`
	Images  []string `json:"images"`
	Point   int      `json:"point"`
	Article *Article `json:"article" xorm:"-"`
}

type CourseService interface {
	GetCoursesByProjectId(projectId int64) ([]*Course, error)
	GetUserLearntCourses(userId int64) ([]*Course, error)

	GetCourseByCourseId(courseId int64) (*Course, error)
	LearntCourseByCourseId(courseId int64, userId int64) (*entities.PointAddRst, error)
	GetArticleById(articleId int64) (*Article, error)
}
