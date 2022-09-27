package feed

import (
	"aed-api-server/internal/interfaces/entities"
	"aed-api-server/internal/interfaces/events"
	"aed-api-server/internal/interfaces/facility"
	"aed-api-server/internal/pkg/db"
	"aed-api-server/internal/pkg/domain/emitter"
	"aed-api-server/internal/pkg/utils"
	"errors"
	"fmt"
	"github.com/go-xorm/xorm"
	log "github.com/sirupsen/logrus"
	"time"
)

func (s *feedSvc) Listen(on facility.OnEvent) {
	on(&events.UserEvent{}, s.dealUserReport)
	on(&events.PostFeed{}, s.dealPostEvent)
	on(&events.PostFeedComment{}, s.sendNewCommentMsg)
}

func (s *feedSvc) dealUserReport(event emitter.DomainEvent) error {
	userEvent, ok := event.(*events.UserEvent)
	if !ok {
		log.Error()
		return nil
	}

	switch userEvent.EventType {
	case entities.GetUserEventTypeOfReport(entities.Report_feedLike):
		if feedId := getFeedIdFromEventParams(userEvent.EventParams); feedId != 0 {
			return s.userLikeOrNot(userEvent.UserId, feedId, true, userEvent.CreatedAt)
		}

	case entities.GetUserEventTypeOfReport(entities.Report_feedLikeCancel):
		if feedId := getFeedIdFromEventParams(userEvent.EventParams); feedId != 0 {
			return s.userLikeOrNot(userEvent.UserId, feedId, false, userEvent.CreatedAt)
		}
	case entities.GetUserEventTypeOfReport(entities.Report_feedShare):
		if feedId := getFeedIdFromEventParams(userEvent.EventParams); feedId != 0 {
			return s.recordShareEvent(userEvent.Id, userEvent.UserId, feedId)
		}
	}
	return nil
}

func (s *feedSvc) userLikeOrNot(userId, feedId int64, isLike bool, opTime time.Time) error {
	return db.Transaction(func(db *xorm.Session) error {
		err := s.persistence.updateFeedLike(userId, feedId, isLike, opTime)
		if err != nil {
			return err
		}
		return s.persistence.updateFeedLikeCount(feedId)
	})
}

func getFeedIdFromEventParams(params []interface{}) int64 {
	if len(params) > 0 {
		param, ok := params[0].(map[string]interface{})
		if !ok {
			return 0
		}
		f, ok := param["feedId"].(float64)
		if !ok {
			return 0
		}

		return int64(f)
	}
	return 0
}

func (s *feedSvc) dealPostEvent(event emitter.DomainEvent) error {
	postFeed := event.(*events.PostFeed)

	var status = entities.FeedStatusAuditPass
	textPass, err := s.checkText(postFeed)
	if err != nil {
		return err
	}

	imagePass, err := s.checkImage(postFeed)
	if err != nil {
		return err
	}

	if !textPass || !imagePass {
		status = entities.FeedStatusAuditDeny
	}

	return s.persistence.updateFeedStatus(postFeed.FeedId, status)
}

func (s *feedSvc) checkText(feed *events.PostFeed) (bool, error) {
	res, err := s.Audit.ScanText(feed.Content)
	if err != nil {
		return false, err
	}

	pass := res.CheckPass()
	if !pass {
		return pass, s.handleBlock(feed, res)
	}

	return pass, nil
}

func (s *feedSvc) checkImage(feed *events.PostFeed) (bool, error) {
	for _, img := range feed.Images {
		res, err := s.Audit.ScanImage(img)
		if err != nil {
			return false, err
		}

		pass := res.CheckPass()
		if !pass {
			return pass, s.handleBlock(feed, res)
		}
	}

	return true, nil

}

func (s *feedSvc) handleBlock(feed *events.PostFeed, r *entities.AuditResult) error {
	user, exists, err := s.User.GetUserById(feed.UserId)
	if err != nil {
		return err
	}

	if !exists {
		return errors.New("user not found")
	}

	return s.DingTalkNotify(feed, user, r)
}

func (s *feedSvc) DingTalkNotify(feed *events.PostFeed, user *entities.SimpleUser, a *entities.AuditResult) error {
	suggestion, label := a.Suggestion()
	return utils.SendDingTalkBot(s.DingTalkUrl, &utils.DingTalkMsg{
		Msgtype: "markdown",
		Markdown: utils.Markdown{
			Title: "帖子审核不通过",
			Text: fmt.Sprintf(`### 帖子审核不通过
- 用户: %s(%d)
- 帖子ID: %d
- 帖子正文: %s
- 审核建议: %s
- 违规原因: %s
                `, user.Nickname, user.ID, feed.FeedId, feed.Content, suggestion, label),
		},
	})
}

func (s *feedSvc) sendNewCommentMsg(event emitter.DomainEvent) error {
	postFeedComment, ok := event.(*events.PostFeedComment)
	if !ok {
		log.Warnf("get wrong event:event%v", event)
	}

	//查询帖子相关用户
	userIds, err := s.persistence.getFeedUsers(postFeedComment.FeedId)
	if err != nil {
		return err
	}

	feed, err := s.persistence.getOneFeedById(postFeedComment.FeedId)
	if err != nil {
		return err
	}

	user, b, err := s.User.GetUserById(postFeedComment.UserId)
	if err != nil {
		log.Error("sendNewCommentMsgToUser GetUserById err:", err)
		return err
	}
	if !b {
		log.Warnf("user(userId=%v) not exited", postFeedComment.UserId)
		return err
	}

	for _, userId := range userIds {
		if userId == postFeedComment.UserId {
			continue
		}
		s.sendNewCommentMsgToUser(userId, postFeedComment.FeedId, feed.Content, postFeedComment.Content, postFeedComment.CreatedAt, user)
	}

	return err
}
func (s *feedSvc) sendNewCommentMsgToUser(userId, feedId int64, feedContent string, commentContent string, at time.Time, by *entities.SimpleUser) {

	openId, err := s.User.GetUserOpenIdById(userId)
	if err != nil {
		log.Error("sendNewCommentMsgToUser GetUserOpenIdById err:", err)
	}

	_, err = s.SubscribeMsg.Send(userId, openId, entities.SMKFeedComment, map[string]interface{}{
		"thing1": map[string]interface{}{
			"value": stringCutAndEclipse(feedContent),
		},
		"thing2": map[string]interface{}{
			"value": stringCutAndEclipse(commentContent),
		},
		"time3": map[string]interface{}{
			"value": at.Format("2006年01月02日 15:04"),
		},
		"thing4": map[string]interface{}{
			"value": by.Nickname,
		},
	}, fmt.Sprintf("/subcontract/feed/detail?id=%d", feedId))

	if err != nil {
		log.Error("sendNewCommentMsgToUser SubscribeMsg.Send err", err)
	}
}
