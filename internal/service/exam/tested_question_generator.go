package exam

import (
	domains2 "aed-api-server/internal/interfaces/domains"
	"aed-api-server/internal/pkg/utils"
	"math/rand"
)

type testedQuestionGenerator struct {
	repository QuestionRepository
	exam       ExamRepository
}

func NewTestedQuestionGenerator(repository QuestionRepository, exam ExamRepository) QuestionGenerator {
	return &testedQuestionGenerator{repository: repository, exam: exam}
}

func (q *testedQuestionGenerator) Generate(projectID int64, userID int64, examType int) ([]*domains2.Question, error) {
	//TODO 针对模拟考试详细制定出题策略
	allQuestions, err := q.repository.ListByProjectIDAndType(projectID, examType)

	if err != nil {
		return nil, err
	}

	if examType == domains2.ExamTypeFormal {
		return allQuestions, nil
	}

	return q.DoGenerateForMock(allQuestions, userID, projectID)
}

func (q *testedQuestionGenerator) DoGenerateForMock(questions []*domains2.Question, userID int64, projectID int64) ([]*domains2.Question, error) {
	set := utils.NewSet[int64]()
	set.AddAll([]int64{96, 97, 98, 99, 100, 153, 154, 155})

	var result []*domains2.Question
	// 不够 8 道直接用

	for _, q := range questions {
		if set.Contains(q.ID) {
			result = append(result, q)
		}
	}

	return result, nil
}

func (q *testedQuestionGenerator) shuffle(questions []*domains2.Question) {
	for i := len(questions) - 1; i >= 0; i-- {
		randomI := rand.Int() % (i + 1)
		questions[i], questions[randomI] = questions[randomI], questions[i]
	}
}
