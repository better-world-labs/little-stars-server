package exam

import (
	domains2 "aed-api-server/internal/interfaces/domains"
	"aed-api-server/internal/pkg/utils"
	"math/rand"
)

type questionGenerator struct {
	repository QuestionRepository
	exam       ExamRepository
}

func NewQuestionGenerator(repository QuestionRepository, exam ExamRepository) QuestionGenerator {
	return &questionGenerator{repository: repository, exam: exam}
}

func (q *questionGenerator) Generate(projectID int64, userID int64, examType int) ([]*domains2.Question, error) {
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

func (q *questionGenerator) DoGenerateForMock(questions []*domains2.Question, userID int64, projectID int64) ([]*domains2.Question, error) {
	size := len(questions)

	// 不够 8 道直接用
	if size < CountQuestionsMock {
		return questions, nil
	}

	// 乱序
	q.shuffle(questions)
	toBeExcluded, err := q.exam.ListCompletedExamQuestionID(projectID, userID)
	if err != nil {
		return nil, err
	}

	toBeExcludedSet := utils.NewSet[int64]()
	toBeExcludedSet.AddAll(toBeExcluded)
	lenToBeExcluded := len(toBeExcludedSet)

	// 排除考过的题后不够 8 道？那不排了
	if size-lenToBeExcluded < CountQuestionsMock {
		return questions[:CountQuestionsMock], nil
	}

	// 排除考过的题
	toBeUsedQuestions := make([]*domains2.Question, 0, size)
	for _, q := range questions {
		if !toBeExcludedSet.Contains(q.ID) {
			toBeUsedQuestions = append(toBeUsedQuestions, q)
		}
	}

	return toBeUsedQuestions[:CountQuestionsMock], nil
}

func (q *questionGenerator) shuffle(questions []*domains2.Question) {
	for i := len(questions) - 1; i >= 0; i-- {
		randomI := rand.Int() % (i + 1)
		questions[i], questions[randomI] = questions[randomI], questions[i]
	}
}
