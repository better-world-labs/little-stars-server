package domains

import (
	"errors"
	"time"
)

const (
	CategorySingleChoice   = 1
	CategoryMultipleChoice = 2
)

type (
	QuestionAnswers []int

	Question struct {
		ID        int64           `xorm:"id"`
		ProjectID int64           `xorm:"project_id"`
		Type      int             `xorm:"type"`
		Title     string          `xorm:"question"`
		Options   []string        `xorm:"options"`
		Images    []string        `xorm:"images"`
		Category  int             `json:"category"`
		Answers   QuestionAnswers `xorm:"answers"`
		CreatedAt time.Time       `xorm:"created_at"`
	}

	QuestionPaper struct {
		Question *Question
		Answers  QuestionAnswers
	}
)

func NewQuestion(projectId int64,
	examType int,
	title string,
	options []string,
	images []string,
	answers []int) (*Question, error) {
	// Exam type check
	if examType != ExamTypeMock && examType != ExamTypeFormal {
		return nil, errors.New("invalid type")
	}

	// Answers format check
	for _, a := range answers {
		if a < 0 || a > len(options) {
			return nil, errors.New("invalid answers")
		}
	}

	var q Question
	q.ProjectID = projectId
	q.Title = title
	q.Options = options
	q.Images = images
	q.Answers = answers
	q.Type = examType
	if len(answers) > 1 {
		q.Category = CategoryMultipleChoice
	} else {
		q.Category = CategorySingleChoice
	}

	q.CreatedAt = time.Now()

	return &q, nil
}

func NewQuestionPaper(question *Question, answers QuestionAnswers) *QuestionPaper {
	return &QuestionPaper{Question: question, Answers: answers}
}

func NewEmptyQuestionPaper(question *Question) *QuestionPaper {
	return &QuestionPaper{Question: question, Answers: []int{}}
}

func (p *QuestionPaper) CheckAnswerRight() bool {
	if len(p.Answers) != len(p.Question.Answers) {
		return false
	}

	for _, a := range p.Answers {
		if !p.isInRightAnswer(a) {
			return false
		}
	}

	return true
}

func (p *QuestionPaper) isInRightAnswer(answer int) bool {
	for _, a := range p.Question.Answers {
		if a == answer {
			return true
		}
	}

	return false
}
