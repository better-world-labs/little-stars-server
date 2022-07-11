package project

import (
	"aed-api-server/internal/interfaces/entities"
	"aed-api-server/internal/interfaces/events"
	"aed-api-server/internal/interfaces/service"
	"aed-api-server/internal/pkg/db"
	"aed-api-server/internal/pkg/domain/emitter"
	"encoding/json"
	"errors"
	"time"
)

type CourseService struct {
}

const (
	articleTableName    = "project_article"
	courseTableName     = "project_course"
	courseUserTableName = "project_course_user"
)

//go:inject-component
func NewCourseService() *CourseService {
	return &CourseService{}
}

type CourseDO struct {
	Id        int      `json:"id"`
	Name      string   `json:"name"`
	TabName   string   `json:"tabName"`
	Images    []string `json:"images"`
	Point     int      `json:"point"`
	ArticleId int64
	CreatedAt time.Time
}

type ArticleDO struct {
	Id        int64     `json:"id" xorm:"id"`
	Type      int       `json:"type" xorm:"type"`
	Content   string    `json:"content" xorm:"content"`
	CreatedAt time.Time `json:"created_at" xorm:"created_at"`
}

type CourseUserDO struct {
	Id        int64 `xorm:"id pk autoincr"`
	UserId    int64
	CourseId  int64
	CreatedAt time.Time
}

func findArticles(articleIds []int64) ([]*ArticleDO, error) {
	var articles = make([]*ArticleDO, 0)
	err := db.Table(articleTableName).In("id", articleIds).Find(&articles)
	if err != nil {
		return nil, err
	}
	return articles, nil
}

func courseDOToCourse(courses []*CourseDO) ([]*service.Course, error) {
	var articleIds []int64
	courseList := make([]*service.Course, 0)

	for _, item := range courses {
		articleIds = append(articleIds, item.ArticleId)
		courseList = append(courseList, &service.Course{
			Id:      item.Id,
			Name:    item.Name,
			TabName: item.TabName,
			Images:  item.Images,
			Point:   item.Point,
		})
	}
	return courseList, nil
}

func (s *CourseService) GetCoursesByProjectId(projectId int64) ([]*service.Course, error) {
	var courses = make([]*CourseDO, 0)
	err := db.Table(courseTableName).Where("project_id = ?", projectId).Find(&courses)
	if err != nil {
		return nil, err
	}
	return courseDOToCourse(courses)
}

func (s *CourseService) GetUserLearntCourses(userId int64) ([]*service.Course, error) {
	var courses = make([]*CourseDO, 0)
	err := db.SQL(`
		select
			course.*
		from `+courseUserTableName+` as user
		inner join `+courseTableName+` as course
			on course.id = user.course_id
		where
			user.user_id = ?
	`, userId).Find(&courses)
	if err != nil {
		return nil, err
	}
	return courseDOToCourse(courses)
}

func (s *CourseService) GetCourseByCourseId(courseId int64) (*service.Course, error) {
	var courseDO CourseDO
	existed, err := db.GetById(courseTableName, courseId, &courseDO)
	if err != nil {
		return nil, err
	}

	if !existed {
		return nil, errors.New("course not found")
	}

	article, err := s.GetArticleById(courseDO.ArticleId)
	if err != nil {
		return nil, err
	}

	return &service.Course{
		Id:      courseDO.Id,
		Name:    courseDO.Name,
		TabName: courseDO.TabName,
		Images:  courseDO.Images,
		Point:   courseDO.Point,
		Article: article,
	}, nil
}

func (s *CourseService) LearntCourseByCourseId(courseId int64, userId int64) (*entities.PointAddRst, error) {
	course, err := s.GetCourseByCourseId(courseId)
	if err != nil {
		return nil, err
	}
	if course == nil {
		return nil, errors.New("not found the course")
	}
	const PointAddDescription = "课程学习加积分"
	rst := entities.PointAddRst{
		Point:       course.Point,
		Description: PointAddDescription,
	}

	var courseUser CourseUserDO
	exist, err := db.Table(courseUserTableName).Where("course_id = ? and user_id = ?", courseId, userId).Exist(&courseUser)
	if err != nil {
		return nil, err
	}
	if !exist {
		user := CourseUserDO{
			UserId:    userId,
			CourseId:  courseId,
			CreatedAt: time.Now(),
		}
		_, err := db.Insert(courseUserTableName, &user)
		if err != nil {
			return nil, err
		}

		err = emitter.Emit(&events.PointsEvent{
			PointsEventType: entities.PointsEventTypeLearntCourse,
			UserId:          userId,
			Params: entities.PointsEventParams{
				RefTable:   courseUserTableName,
				RefTableId: user.Id,
			},
		})
		if err != nil {
			return nil, err
		}
	}
	return &rst, nil
}

func (s *CourseService) GetArticleById(articleId int64) (*service.Article, error) {
	articles, err := findArticles([]int64{articleId})
	if err != nil {
		return nil, err
	}
	if len(articles) > 0 {
		var article service.Article
		_ = json.Unmarshal([]byte(articles[0].Content), &article)
		article.Id = articles[0].Id
		return &article, nil
	}
	return nil, nil
}
