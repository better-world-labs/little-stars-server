package exam

import (
	domains2 "aed-api-server/internal/interfaces/domains"
	"github.com/go-xorm/xorm"
	"io"
)

const (
	CountQuestionsMock = 8
)

type (
	//ExamRepository  Exam 存储库
	ExamRepository interface {
		Create(exam *domains2.Exam) error
		UpdateWithSession(session *xorm.Session, exam *domains2.Exam) error
		CompareAndSwapExamWithSession(session *xorm.Session, excepted *domains2.Exam, value *domains2.Exam) (bool, error)
		Update(exam *domains2.Exam) error
		GetByID(id int64) (*ExamDo, bool, error)
		GetLastByExaminerAndTypeUnCompleted(projectID int64, userId int64, examType int) (*ExamDo, bool, error)
		ListCompletedExamQuestionID(projectId, userId int64) (res []int64, err error)
		ListByExaminerAndTypeCompleted(projectID int64, userId int64, examType int, latest int) ([]*ExamDo, error)
		ListByExaminerAndTypeCompletedWithSession(session *xorm.Session, projectID int64, userId int64, examType int, latest int) ([]*ExamDo, error)
	}

	//QuestionRepository  Question 存储库
	QuestionRepository interface {
		BatchCreate(questions []*domains2.Question) error
		ListByIDs(ids []int64) ([]*domains2.Question, error)
		GetByID(id int64) (*domains2.Question, bool, error)
		ListByProjectIDAndType(projectID int64, _type int) ([]*domains2.Question, error)
		BatchDelete(ids []int64) error
	}

	QuestionPicker interface {
		PickQuestions(reader io.Reader, projectID int64, examType int) ([]*domains2.Question, error)
	}

	//QuestionGenerator	 题目生成器
	// 根据一定规则生成一定数量的考题
	QuestionGenerator interface {
		Generate(projectID int64, userID int64, examType int) ([]*domains2.Question, error)
	}
)
