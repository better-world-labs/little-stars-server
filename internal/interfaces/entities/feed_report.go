package entities

import "time"

const (
	FeedReportStatusTodo     FeedReportStatus = 0
	FeedReportStatusResolved FeedReportStatus = 1
)

type (
	FeedReportStatus int

	FeedReport struct {
		Id        int64            `json:"id"`
		Type      string           `json:"type" binding:"required"`
		FeedId    int64            `json:"feedId"`
		Content   string           `json:"content" binding:"required"`
		Status    FeedReportStatus `json:"status"`
		CreatedBy int64            `json:"createdBy"`
		CreatedAt time.Time        `json:"createdAt"`
	}
)
