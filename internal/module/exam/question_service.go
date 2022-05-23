package exam

import (
	"io"
)

type questionService struct {
	questionPicker     QuestionPicker
	questionRepository QuestionRepository
}

func NewQuestionService() *questionService {
	return &questionService{NewExcelQuestionFilePicker(), NewQuestionRepository()}
}

func (s *questionService) ImportQuestions(reader io.ReadCloser, projectID int64, examType int) error {
	questions, err := s.questionPicker.PickQuestions(reader, projectID, examType)
	if err != nil {
		return err
	}

	err = s.questionRepository.BatchCreate(questions)
	return err
}
