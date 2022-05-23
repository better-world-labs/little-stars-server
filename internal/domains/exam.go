package domains

import (
	"aed-api-server/internal/pkg/global"
	"time"
)

var (
//ErrorAlreadyCompleted = base.NewError("exam", "exam already completed")
//ErrorUnknownQuestion  = base.NewError("exam", "unknown question in this exam")
)

const (
	ExamTypeMock   = 1
	ExamTypeFormal = 2
	ExamPassScore  = 100
)

type (
	Exam struct {
		ID          int64
		Type        int
		ProjectID   int64
		Examiner    int64
		CreatedAt   time.Time
		CompletedAt *time.Time
		Completed   bool
		Score       int
		QuestionsSorted []int64

		//questionsIndex map[int64]int
		questions   map[int64]*QuestionPaper
	}
)

func ParseExam(id int64,
	Type int,
	projectID int64,
	examiner int64,
	sortedIndex []int64,
	questions []*Question,
	paper map[int64][]int,
	createdAt time.Time,
	completedAt *time.Time,
	completed bool,
	score int,
) (*Exam, error) {
	exam := NewExam(Type, projectID, examiner,sortedIndex, questions)
	err := exam.InsertAnswer(paper)
	if err != nil {
		return nil, err
	}

	exam.ID = id
	exam.CreatedAt = createdAt
	exam.CompletedAt = completedAt
	exam.Completed = completed
	exam.Score = score

	return exam, nil
}

func NewExam(examType int, projectID int64, examinerID int64, sortedIndex []int64, questions []*Question) *Exam {
	count := len(questions)
	questionsIndexLen := len(sortedIndex)

	// 没有指定题目 index 则生成一个
	if count != questionsIndexLen {
		for _, q := range questions {
			sortedIndex = append(sortedIndex, q.ID)
		}
	}

	var e Exam
	e.questions = make(map[int64]*QuestionPaper)

	e.ProjectID = projectID
	e.Examiner = examinerID
	e.Type = examType
	e.CreatedAt = time.Now()
	e.QuestionsSorted = sortedIndex

	for _, q := range questions {
		e.questions[q.ID] = NewEmptyQuestionPaper(q)
	}

	return &e
}

func (e *Exam) GetQuestionPaper(questionID int64) (*QuestionPaper, bool) {
	paper, exists := e.questions[questionID]
	return paper, exists
}

func (e *Exam) SaveExam(paper map[int64][]int) error {
	if e.Completed {
		return global.ErrorAlreadyCompleted
	}

	return e.InsertAnswer(paper)
}

func (e *Exam) InsertAnswer(paper map[int64][]int) error {
	for questionID, answers := range paper {
		q, exists := e.GetQuestionPaper(questionID)
		if !exists {
			return global.ErrorUnknownQuestion
		}

		q.Answers = answers
	}
	return nil
}

func (e *Exam) SubmitExam(paper map[int64][]int) error {
	err := e.SaveExam(paper)
	if err != nil {
		return err
	}

	return e.complete()
}

func (e *Exam) CheckPass() bool {
	return e.Completed && e.Score >= ExamPassScore
}

func (e *Exam) complete() error {
	if e.Completed {
		return global.ErrorAlreadyCompleted
	}

	now := time.Now()
	e.Completed = true
	e.CompletedAt = &now
	e.Score = e.evaluateScore()

	return nil
}

func (e *Exam) evaluateScore() (score int) {
	right := 0

	for _, p := range e.questions {
		if p.CheckAnswerRight() {
			right++
		}
	}

	return int(100 * (float64(right) / float64(len(e.questions))))
}
