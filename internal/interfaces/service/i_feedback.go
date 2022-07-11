package service

import (
	"aed-api-server/internal/interfaces/entities"
	"io"
	"time"
)

type FeedbackService interface {
	SubmitFeedback(userId int64, feedback *entities.Feedback) error
	ExportFeedback(beginDate time.Time, endDate time.Time, writer *io.PipeWriter)
}
