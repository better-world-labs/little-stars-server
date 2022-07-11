package exam

import (
	"aed-api-server/internal/interfaces/domains"
	"aed-api-server/internal/pkg/db"
	"github.com/go-xorm/xorm"
)

type examRepository struct {
	TableName string
}

func NewExamRepository() ExamRepository {
	return &examRepository{TableName: "exam_new"}
}

func (e *examRepository) Create(exam *domains.Exam) error {
	do := ToModel(exam)
	_, err := db.GetEngine().Table(e.TableName).Insert(do)
	exam.ID = do.ID
	return err
}

func (e *examRepository) Update(exam *domains.Exam) error {
	session := db.GetSession()
	defer session.Close()
	return e.UpdateWithSession(session, exam)
}

func (e *examRepository) CompareAndSwapExamWithSession(session *xorm.Session, excepted *domains.Exam, value *domains.Exam) (bool, error) {
	exceptedDo := ToModel(excepted)
	valueDo := ToModel(value)
	rows, err := session.Table(e.TableName).Where("id = ? and completed = ? and score = ?",
		exceptedDo.ID, exceptedDo.Completed, exceptedDo.Score).UseBool("completed").
		Update(valueDo)
	if err != nil {
		return false, err
	}

	return rows > 0, nil
}

func (e *examRepository) UpdateWithSession(session *xorm.Session, exam *domains.Exam) error {
	_, err := session.Table(e.TableName).ID(exam.ID).UseBool("completed").Update(ToModel(exam))
	return err
}

func (e *examRepository) GetByID(id int64) (*ExamDo, bool, error) {
	var m ExamDo
	exists, err := db.GetEngine().Table(e.TableName).ID(id).Get(&m)
	return &m, exists, err
}

func (e *examRepository) ListCompletedExamQuestionID(projectId, userId int64) (res []int64, err error) {

	var questionIDs []struct {
		Questions []int64 `xorm:"questions"`
	}

	err = db.GetEngine().Table(e.TableName).Cols("questions").Where("examiner = ? and project_id = ? and completed = true", userId, projectId).Find(&questionIDs)
	if err != nil {
		return nil, err
	}

	for i := 0; i < len(questionIDs); i++ {
		for j := 0; j < len(questionIDs[i].Questions); j++ {
			res = append(res, questionIDs[i].Questions[j])
		}
	}

	return
}

func (e *examRepository) ListByExaminerAndTypeCompleted(projectID, userId int64, examType int, latest int) ([]*ExamDo, error) {
	session := db.GetSession()
	defer session.Close()
	return e.ListByExaminerAndTypeCompletedWithSession(session, projectID, userId, examType, latest)
}

func (e *examRepository) GetLastByExaminerAndTypeUnCompleted(projectID, userId int64, examType int) (*ExamDo, bool, error) {
	var m ExamDo
	exists, err := db.GetEngine().Table(e.TableName).
		Where("examiner = ? and project_id = ? and type = ? and completed = 0",
			userId,
			projectID,
			examType).Get(&m)
	return &m, exists, err
}

func (e *examRepository) ListByExaminerAndTypeCompletedWithSession(session *xorm.Session, projectID int64, userId int64, examType int, latest int) ([]*ExamDo, error) {
	var arr []*ExamDo
	session.Table(e.TableName).
		Where("examiner = ? and project_id = ? and type = ? and completed = 1",
			userId,
			projectID,
			examType).
		Desc("completed_at")

	if latest > 0 {
		session = session.Limit(latest, 0)
	}

	err := session.Find(&arr)
	return arr, err
}
