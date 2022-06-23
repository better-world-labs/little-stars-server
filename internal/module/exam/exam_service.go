package exam

import (
	"aed-api-server/internal/domains"
	"aed-api-server/internal/interfaces"
	"aed-api-server/internal/interfaces/entities"
	"aed-api-server/internal/interfaces/events"
	"aed-api-server/internal/interfaces/service"
	"aed-api-server/internal/pkg"
	"aed-api-server/internal/pkg/base"
	"aed-api-server/internal/pkg/config"
	"aed-api-server/internal/pkg/db"
	"aed-api-server/internal/pkg/response"
	"fmt"
	"github.com/go-xorm/xorm"
)

type examService struct {
	examRepository     ExamRepository
	questionRepository QuestionRepository
	questionsGenerator QuestionGenerator

	ProjectService service.ProjectService `inject:"-"`
	CertService    service.CertService    `inject:"-"`
}

func NewExamService(c *config.AppConfig) *examService {
	repository := NewQuestionRepository()
	exam := NewExamRepository()
	var generator QuestionGenerator
	if c.Exam.Debug {
		generator = NewTestedQuestionGenerator(repository, exam)
	} else {
		generator = NewQuestionGenerator(repository, exam)
	}

	return &examService{
		examRepository:     exam,
		questionRepository: repository,
		questionsGenerator: generator,
	}
}

func (e *examService) Start(projectId int64, userId int64, examType int) (*domains.Exam, error) {
	//TODO get Unsubmitted use union index
	_, exists, err := e.examRepository.GetLastByExaminerAndTypeUnCompleted(projectId, userId, examType)
	if err != nil {
		return nil, err
	}

	if exists {
		return nil, response.ErrorExamUnSubmit
	}

	return e.StartNewExam(projectId, userId, examType)
}

func (e *examService) ContinueLastExam(last *ExamDo) (*domains.Exam, error) {
	questions, err := e.questionRepository.ListByIDs(last.Questions)
	if err != nil {
		return nil, err
	}

	exam, err := ParseExam(last, last.Questions, questions)
	if err != nil {
		return nil, err
	}

	return exam, nil
}

func (e *examService) StartNewExam(projectId, userId int64, examType int) (*domains.Exam, error) {
	questions, err := e.questionsGenerator.Generate(projectId, userId, examType)
	if err != nil {
		return nil, err
	}

	exam := domains.NewExam(examType, projectId, userId, nil, questions)
	err = e.examRepository.Create(exam)

	return exam, err
}

func (e *examService) Save(examID, userId int64, paper map[int64][]int) error {
	exam, exists, err := e.GetByID(examID)
	if err != nil {
		return err
	}

	if !exists {
		return base.NewError("exam.service", "exam not found")
	}

	if exam.Examiner != userId {
		return response.ErrorExamOwnerError
	}

	err = exam.SaveExam(paper)
	if err != nil {
		return err
	}

	err = e.examRepository.Update(exam)
	return err
}

func (e *examService) Submit(examID, userId int64, paper map[int64][]int) (string, string, []*entities.DealPointsEventRst, error) {
	var pointsRst []*entities.DealPointsEventRst
	exam, exists, err := e.GetByID(examID)
	if err != nil {
		return "", "", pointsRst, err
	}

	if !exists {
		return "", "", pointsRst, base.NewError("exam.service", "exam not found")
	}

	if exam.Examiner != userId {
		return "", "", pointsRst, response.ErrorExamOwnerError
	}

	excepted := *exam
	err = exam.SubmitExam(paper)
	if err != nil {
		return "", "", pointsRst, err
	}

	var certImg string
	var certNum string
	err = db.Begin(func(session *xorm.Session) (err error) {
		ok, err := e.examRepository.CompareAndSwapExamWithSession(session, &excepted, exam)
		if err != nil {
			return err
		}

		if !ok {
			return response.ErrorConcurrentOperation
		}

		certImg, certNum, err = e.handleExamCompleted(session, &pointsRst, exam)
		if err != nil {
			return
		}

		return
	})

	return certImg, certNum, pointsRst, err
}

func (e *examService) ListLatestSubmitted(projectId int64, userId int64, examType int, latest int) ([]*domains.Exam, error) {
	//TODO 方法有点长，后边封装一下
	list, err := e.examRepository.ListByExaminerAndTypeCompleted(projectId, userId, examType, latest)
	if err != nil {
		return nil, err
	}

	var exams []*domains.Exam
	questionsCache := make(map[int64]*domains.Question, len(list))
	var questionIDs []int64
	for _, e := range list {
		for _, q := range e.Questions {
			questionIDs = append(questionIDs, q)
		}
	}

	questions, err := e.questionRepository.ListByIDs(questionIDs)
	if err != nil {
		return nil, err
	}

	for _, q := range questions {
		questionsCache[q.ID] = q
	}

	for _, e := range list {
		var questionGroup []*domains.Question
		for _, q := range e.Questions {
			question, exists := questionsCache[q]
			if !exists {
				continue
			}
			questionGroup = append(questionGroup, question)
		}

		exam, err := ParseExam(e, e.Questions, questionGroup)
		if err != nil {
			return nil, err
		}

		exams = append(exams, exam)
	}

	return exams, nil
}

func (e *examService) GetLatestUnSubmitted(projectId int64, userId int64, examType int) (*domains.Exam, bool, error) {
	latest, exists, err := e.examRepository.GetLastByExaminerAndTypeUnCompleted(projectId, userId, examType)
	if err != nil {
		return nil, false, err
	}

	if !exists {
		return nil, exists, nil
	}

	questions, err := e.questionRepository.ListByIDs(latest.Questions)
	if err != nil {
		return nil, false, err
	}

	exam, err := ParseExam(latest, latest.Questions, questions)
	return exam, true, err
}

func (e *examService) GetByID(examID int64) (*domains.Exam, bool, error) {
	examDo, exists, err := e.examRepository.GetByID(examID)
	if err != nil || !exists {
		return nil, exists, err
	}

	questions, err := e.questionRepository.ListByIDs(examDo.Questions)
	if err != nil {
		return nil, exists, err
	}

	exam, err := ParseExam(examDo, examDo.Questions, questions)
	return exam, exists, err
}

func (e *examService) handleFormalExamPassed(exam *domains.Exam) (string, string, error) {
	err := e.ProjectService.DoCertification(exam.ProjectID, exam.Examiner)
	if err != nil {
		return "", "", nil
	}

	certImg, certNum, err := e.CertService.CreateCert(exam.ProjectID, exam.Examiner)
	if err != nil {
		return "", "", err
	}

	//err = producer.Publish(record.CreateRecordProjectCertificated(exam.Examiner, exam.ProjectID))
	//if err != nil {
	//	return err
	//}

	return certImg, certNum, nil
}

func (e *examService) handleExamCompleted(session *xorm.Session, pointsRst *[]*entities.DealPointsEventRst, exam *domains.Exam) (certImg, certNum string, err error) {
	// 模拟考试完成
	if exam.Type == domains.ExamTypeMock {
		exams, err := e.examRepository.ListByExaminerAndTypeCompletedWithSession(session,
			exam.ProjectID, exam.Examiner, exam.Type, 7,
		)
		if err != nil {
			return "", "", err
		}

		level := MatchLevel(exams)
		fmt.Printf("change level to %d", level)
		err = e.ProjectService.UpdateUserProjectLevel(exam.ProjectID, exam.Examiner, level)
		if err != nil {
			return "", "", err
		}
		pointsEvent := interfaces.S.PointsScheduler.BuildPointsEventTypeMockedExam(exam.Examiner, exam.ID, exam.Score)

		times, err := interfaces.S.Points.GetUserPointsEventTimes(exam.Examiner, entities.PointsEventTypeExam)
		if err != nil {
			return "", "", err
		}

		if times < pkg.UserPointsMaxTimesMockExam {
			event, err := interfaces.S.PointsScheduler.DealPointsEvent(pointsEvent)
			if err != nil {
				return "", "", err
			}

			*pointsRst = append(*pointsRst, event)
		}
	}

	// 认证考试通过
	if exam.Type == domains.ExamTypeFormal && exam.CheckPass() {
		certImg, certNum, err = e.handleFormalExamPassed(exam)
		if err != nil {
			return "", "", err
		}

		times, err := interfaces.S.Points.GetUserPointsEventTimes(exam.Examiner, entities.PointsEventTypeCertificated)
		if err != nil {
			return "", "", err
		}

		if times < 1 {
			event, err := interfaces.S.PointsScheduler.DealPointsEvent(&events.PointsEvent{
				PointsEventType: entities.PointsEventTypeCertificated,
				UserId:          exam.Examiner,
				Params: map[string]interface{}{
					"examId": exam.ID,
				},
			})
			if err != nil {
				return "", "", err
			}

			*pointsRst = append(*pointsRst, event)
		}
	}

	return
}

func ParseExam(model *ExamDo, sortedIndex []int64, questions []*domains.Question) (*domains.Exam, error) {
	return domains.ParseExam(
		model.ID,
		model.Type,
		model.ProjectID,
		model.Examiner,
		sortedIndex,
		questions,
		model.Answers,
		model.CreatedAt,
		model.CompletedAt,
		model.Completed,
		model.Score)
}
