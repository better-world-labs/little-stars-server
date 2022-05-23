package exam

import (
	"aed-api-server/internal/domains"
	"errors"
	"github.com/xuri/excelize/v2"

	//"github.com/360EntSecGroup-Skylar/excelize"
	"io"
	"regexp"
)

type excelQuestionFilePicker struct {
}

func NewExcelQuestionFilePicker() QuestionPicker {
	return &excelQuestionFilePicker{}
}

func (e excelQuestionFilePicker) PickQuestions(reader io.Reader, projectID int64, examType int) ([]*domains.Question, error) {
	var questions []*domains.Question
	f, err := excelize.OpenReader(reader)
	if err != nil {
		return nil, err
	}
	index := f.GetActiveSheetIndex()
	sheetMap := f.GetSheetMap()
	rows, _ := f.GetRows(sheetMap[index])

	for i, row := range rows {
		if i == 0 {
			continue
		}

		question, err := e.parseQuestion(row, projectID, examType)
		if err != nil {
			return nil, err
		}

		questions = append(questions, question)
	}

	return questions, nil
}

func (e excelQuestionFilePicker) parseQuestion(row []string, projectID int64, examType int) (*domains.Question, error) {
	if len(row) < 7 {
		return nil, errors.New("invalid row length")
	}

	// Options
	var options []string
	for i := 2; i < 6; i++ {
		if row[i] != "" {
			options = append(options, row[i])
		}
	}

	// Answers
	answersStr := row[6]
	match, err := regexp.Match("[^0-9]+", []byte(answersStr))
	if err != nil {
		return nil, err
	}

	if match {
		return nil, errors.New("invalid answer format")
	}

	var answers []int
	for _, a := range answersStr {
		answers = append(answers, int(a-48-1))
	}

	// Question
	question, err := domains.NewQuestion(projectID, examType, row[1], options, nil, answers)
	return question, err
}
