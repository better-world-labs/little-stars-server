package report

import (
	"aed-api-server/internal/interfaces/entities"
	"aed-api-server/internal/interfaces/service"
	"aed-api-server/internal/pkg/db"
	"aed-api-server/internal/pkg/utils"
	"errors"
	"fmt"
	"github.com/go-xorm/xorm"
	"time"
)

type FeedReportImpl struct {
	persistence IPersistence

	DingTalkUrl string              `conf:"donation-apply-notify"`
	Feed        service.IFeed       `inject:"-"`
	User        service.UserService `inject:"-"`
}

//go:inject-component
func NewFeedReport() service.IFeedReport {
	return &FeedReportImpl{
		persistence: NewPersistenceImpl("feed_report"),
	}
}

func (f *FeedReportImpl) dingTalkNotify(report *entities.FeedReport, user *entities.SimpleUser, feed *entities.Feed) error {
	return utils.SendDingTalkBot(f.DingTalkUrl, &utils.DingTalkMsg{
		Msgtype: "markdown",
		Markdown: utils.Markdown{
			Title: "帖子举报",
			Text: fmt.Sprintf(`### 帖子举报
- 用户: %s(%d)
- 帖子ID: %d
- 帖子正文: %s
- 举报类型: %s
- 举报内容: %s
                `, user.Nickname, user.ID, feed.Id, feed.Content, report.Type, report.Content),
		},
	})
}

func (f *FeedReportImpl) Create(report *entities.FeedReport) error {
	report.CreatedAt = time.Now()

	user, exists, err := f.User.GetUserById(report.CreatedBy)
	if err != nil {
		return err
	}

	if !exists {
		return errors.New("user not found")
	}

	feed, err := f.Feed.GetFeedById(report.FeedId)
	if err != nil {
		return err
	}

	if feed == nil {
		return errors.New("feed not found")
	}

	return db.Transaction(func(session *xorm.Session) error {
		err := f.persistence.Create(report)
		if err != nil {
			return err
		}

		return f.dingTalkNotify(report, user, feed)
	})
}

func (f *FeedReportImpl) List() ([]*entities.FeedReport, error) {
	return f.persistence.List()
}

func (f *FeedReportImpl) GetById(id int64) (*entities.FeedReport, bool, error) {
	return f.persistence.GetById(id)
}

func (f *FeedReportImpl) UpdateStatus(id int64, status entities.FeedReportStatus) error {
	report, exists, err := f.GetById(id)
	if err != nil {
		return err
	}

	if !exists {
		return errors.New("not found")
	}

	report.Status = status
	return f.persistence.Update(report)
}
