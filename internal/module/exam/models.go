package exam

import (
	"aed-api-server/internal/domains"
	"time"
)

type ExamDo struct {
	ID          int64           `xorm:"id pk autoincr"`
	Type        int             `xorm:"type"`
	ProjectID   int64           `xorm:"project_id"`
	Examiner    int64           `xorm:"examiner"`
	Questions   []int64         `xorm:"questions"`
	Answers     map[int64][]int `xorm:"answers"`
	CreatedAt   time.Time       `xorm:"created_at"`
	CompletedAt *time.Time      `xorm:"completed_at"`
	Completed   bool            `xorm:"completed"`
	Score       int             `xorm:"score"`
}

func ToModel(exam *domains.Exam) *ExamDo {
	size := len(exam.QuestionsSorted)
	questions := make([]int64, 0, size)
	paper := make(map[int64][]int, size)

	for _, q := range exam.QuestionsSorted {
		if question, exists := exam.GetQuestionPaper(q); exists {
			questions = append(questions, q)
			paper[q] = question.Answers
		}
	}

	return &ExamDo{
		ID:          exam.ID,
		Type:        exam.Type,
		ProjectID:   exam.ProjectID,
		Examiner:    exam.Examiner,
		CreatedAt:   exam.CreatedAt,
		CompletedAt: exam.CompletedAt,
		Completed:   exam.Completed,
		Score:       exam.Score,
		Questions:   questions,
		Answers:     paper,
	}
}
