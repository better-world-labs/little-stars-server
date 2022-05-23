package exam

import (
	"aed-api-server/internal/domains"
	"aed-api-server/internal/pkg/db"
)

type questionRepository struct {
	TableName string
}

func NewQuestionRepository() QuestionRepository {
	return &questionRepository{TableName: "question"}
}

func (q *questionRepository) BatchCreate(questions []*domains.Question) error {
	_, err := db.GetEngine().Table(q.TableName).Insert(questions)
	return err
}

func (q *questionRepository) GetByID(id int64) (*domains.Question, bool, error) {
	var qu domains.Question
	exists, err := db.GetEngine().Table(q.TableName).ID(id).Get(&qu)
	return &qu, exists, err
}

func (q *questionRepository) ListByIDs(ids []int64) ([]*domains.Question, error) {
	dom := make([]*domains.Question, 0)
	err := db.GetEngine().Table(q.TableName).In("id", ids).Find(&dom)
	return dom, err
}

func (q *questionRepository) ListByProjectIDAndType(projectID int64, _type int) ([]*domains.Question, error) {
	dom := make([]*domains.Question, 0)
	err := db.GetEngine().Table(q.TableName).Where("project_id = ? and type = ?", projectID, _type).Find(&dom)
	return dom, err
}

func (q *questionRepository) BatchDelete(ids []int64) error {
	panic("implement me")
}
