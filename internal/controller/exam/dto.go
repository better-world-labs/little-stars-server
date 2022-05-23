package exam

import (
	"aed-api-server/internal/domains"
	"aed-api-server/internal/pkg/global"
)

type (
	StartExamDto struct {
		ID        int64         `json:"id"`
		Questions []QuestionDto `json:"questions"`
	}
	SimpleListInfoDto struct {
		ID        int64                 `json:"id"`
		StartAt   global.FormattedTime  `json:"startAt"`
		SubmitAt  *global.FormattedTime `json:"submitAt,omitempty"`
		Score     int                   `json:"score"`
		Completed bool                  `json:"completed"`
	}
	SimpleListInfoWithQuestionsDto struct {
		SimpleListInfoDto

		Questions []QuestionWithAnswerDto `json:"questions"`
	}
	QuestionDto struct {
		ID       int64    `json:"id"`
		Question string   `json:"question"`
		Category int      `json:"category"`
		Options  []string `json:"options"`
		Images   []string `json:"images"`
	}

	QuestionWithAnswerDto struct {
		QuestionDto

		Answers []int `json:"answers"`
	}

	AnswerDetail struct {
		QuestionWithAnswerDto
		Right          bool  `json:"right"`
		CorrectAnswers []int `json:"correctAnswers"`
	}

	DetailDto struct {
		SimpleListInfoDto

		IncorrectedIndex []int          `json:"incorrectedIndex"`
		Questions        []AnswerDetail `json:"questions"`
	}
)

func NewStartExamVo(exam domains.Exam) StartExamDto {
	var dto StartExamDto
	dto.ID = exam.ID
	dto.Questions = make([]QuestionDto, 0, len(exam.QuestionsSorted))

	for _, q := range exam.QuestionsSorted {
		if question, exists := exam.GetQuestionPaper(q); exists {
			qDto := NewQuestionDto(q,
				question.Question.Title,
				question.Question.Category,
				question.Question.Options,
				question.Question.Images)

			dto.Questions = append(dto.Questions, *qDto)
		}
	}

	return dto
}

func ParseSimpleListInfoDto(exam *domains.Exam) (dto SimpleListInfoDto) {
	dto.ID = exam.ID
	dto.StartAt = global.FormattedTime(exam.CreatedAt)
	dto.Score = exam.Score
	dto.Completed = exam.Completed
	if exam.CompletedAt != nil {
		time := global.FormattedTime(*exam.CompletedAt)
		dto.SubmitAt = &time
	}

	return
}

func ParseSimpleListInfoWithQuestionDto(exam *domains.Exam) (dto SimpleListInfoWithQuestionsDto) {
	simple := ParseSimpleListInfoDto(exam)
	answers := make([]QuestionWithAnswerDto, 0, len(exam.QuestionsSorted))

	for _, v := range exam.QuestionsSorted {
		if question, exists := exam.GetQuestionPaper(v); exists {
			answer := NewQuestionWithAnswersDto(v,
				question.Question.Title,
				question.Question.Category,
				question.Question.Options,
				question.Question.Images,
				question.Answers)
			answers = append(answers, *answer)
		}
	}

	dto.SimpleListInfoDto = simple
	dto.Questions = answers

	return
}

func ParseDetailDTo(exam *domains.Exam) (dto DetailDto) {
	simple := ParseSimpleListInfoDto(exam)
	questions := make([]AnswerDetail, 0, len(exam.QuestionsSorted))
	incorrectIndex := make([]int, 0)
	for index, q := range exam.QuestionsSorted {
		if question, exists := exam.GetQuestionPaper(q); exists {
		detail := NewAnswerDetail(q,
			question.Question.Title,
			question.Question.Category,
			question.Question.Options,
			question.Question.Images,
			question.Answers,
			question.CheckAnswerRight(),
			question.Question.Answers)

		if !detail.Right {
			incorrectIndex = append(incorrectIndex, index)
		}

		questions = append(questions, *detail)
		}
	}

	dto.SimpleListInfoDto = simple
	dto.Questions = questions
	dto.IncorrectedIndex = incorrectIndex

	return
}

func NewQuestionDto(
	id int64,
	question string,
	category int,
	options []string,
	images []string,
) *QuestionDto {
	if images == nil {
		images = make([]string, 0)
	}
	return &QuestionDto{
		ID:       id,
		Question: question,
		Category: category,
		Options:  options,
		Images:   images,
	}
}

func NewQuestionWithAnswersDto(
	id int64,
	question string,
	category int,
	options []string,
	images []string,
	answers []int) *QuestionWithAnswerDto {
	if answers == nil {
		answers = make([]int, 0)
	}
	return &QuestionWithAnswerDto{
		*NewQuestionDto(id, question, category, options, images),
		answers,
	}
}

func NewAnswerDetail(
	id int64,
	question string,
	category int,
	options []string,
	images []string,
	answers []int,
	Right bool,
	CorrectAnswers []int,
) *AnswerDetail {
	if CorrectAnswers == nil {
		CorrectAnswers = make([]int, 0)
	}
	return &AnswerDetail{
		*NewQuestionWithAnswersDto(id, question, category, options, images, answers),
		Right,
		CorrectAnswers,
	}
}
