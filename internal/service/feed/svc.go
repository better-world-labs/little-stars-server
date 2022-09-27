package feed

import (
	"aed-api-server/internal/interfaces/entities"
	"aed-api-server/internal/interfaces/events"
	"aed-api-server/internal/interfaces/facility"
	"aed-api-server/internal/interfaces/service"
	"aed-api-server/internal/pkg/db"
	"aed-api-server/internal/pkg/utils"
	"errors"
	"fmt"
	"github.com/go-xorm/xorm"
	log "github.com/sirupsen/logrus"
)

//go:inject-component
func NewFeedSvc() service.IFeed {
	return &feedSvc{
		persistence: &persistence{},
	}
}

const DefaultPageSize = 10
const MaxPageSize = 50

type feedSvc struct {
	facility.Sender `inject:"sender"`
	Audit           service.IContentAudit   `inject:"-"`
	User            service.UserService     `inject:"-"`
	Donation        service.DonationService `inject:"-"`
	SubscribeMsg    service.SubscribeMsg    `inject:"-"`

	persistence iPersistence

	DingTalkUrl string `conf:"donation-apply-notify"`
}

//GetCatalogList 获取帖子板块列表
func (s *feedSvc) GetCatalogList() ([]*entities.Catalog, error) {
	return s.persistence.getCatalogList()
}

//GetFeeds 获取帖子流
func (s *feedSvc) GetFeeds(catalogId int, cursor string, size int, userId int64) (list []*entities.Feed, nextCursor string, err error) {
	// 参数检查
	size = limitPageSize(size)
	if catalogId > 0 && !s.persistence.existedCatalogId(catalogId) {
		return nil, "", errors.New("板块不存在")
	}

	var beforeId int64 = 0
	if cursor != "" {
		beforeId, err = parseCursor(cursor)
		if err != nil {
			return nil, "", err
		}
	}

	list, err = s.persistence.getFeeds(beforeId, catalogId, size, userId)
	if err != nil {
		return nil, "", err
	}

	if list == nil {
		list = make([]*entities.Feed, 0)
	}
	if len(list) >= size {
		endFeed := list[len(list)-1]
		nextCursor = toCursor(endFeed.Id)
	}

	//creator计算
	err = s.patchFeedCreator(list)
	if err != nil {
		return nil, "", err
	}

	if userId > 0 {
		//我是否点赞、评论计算
		err = s.patchFeedsLikeAndComment(list, userId)
		if err != nil {
			return nil, "", err
		}
	}
	return list, nextCursor, nil
}

//GetMyFeeds 获取和我相关的帖子
func (s *feedSvc) GetMyFeeds(cursor string, size int, userId int64) (list []*entities.MyFeed, nextCursor string, err error) {
	// 参数检查
	size = limitPageSize(size)
	var beforeId int64 = 0
	if cursor != "" {
		beforeId, err = parseCursor(cursor)
		if err != nil {
			return nil, "", err
		}
	}

	list, err = s.persistence.getMyFeeds(beforeId, size, userId)
	if err != nil {
		return nil, "", err
	}

	if list == nil {
		list = make([]*entities.MyFeed, 0)
	}
	if len(list) >= size {
		endFeed := list[len(list)-1]
		nextCursor = toCursor(endFeed.Id)
	}

	//creator计算
	err = s.patchMyFeedCreator(list)
	if err != nil {
		return nil, "", err
	}

	//我是否点赞、评论计算
	err = s.patchMyFeedsLikeAndComment(list, userId)
	if err != nil {
		return nil, "", err
	}

	return list, nextCursor, nil
}

func (s *feedSvc) GetMyFeedsCount(userId int64) (count int, err error) {
	return s.persistence.getMyFeedsCount(userId)
}

func (s *feedSvc) GetFeedById(feedId int64) (feed *entities.Feed, err error) {
	f, err := s.persistence.getOneFeedById(feedId)
	if err != nil {
		return
	}

	if f == nil {
		return nil, nil
	}

	return f.Feed, err
}

//GetFeedDetail 获取帖子详情
func (s *feedSvc) GetFeedDetail(feedId int64, userId int64) (feed *entities.FeedWithComment, err error) {
	feedDo, err := s.persistence.getOneFeedById(feedId)
	if err != nil {
		return nil, err
	}

	if !(feedDo.CreatedBy == userId && feedDo.Status == entities.FeedStatusNeedAudit) {
		if err = s.checkFeedStatus(feedDo); err != nil {
			return nil, err
		}
	}

	catalog, err := s.persistence.getCatalogById(feedDo.CatalogId)
	if err != nil {
		return nil, err
	}
	if catalog == nil {
		return nil, errors.New("没找到相关板块")
	}

	warp := []*entities.Feed{feedDo.Feed}
	if err = s.patchFeedCreator(warp); err != nil {
		return nil, err
	}

	if err = s.patchFeedsLikeAndComment(warp, userId); err != nil {
		return nil, err
	}

	comments, err := s.persistence.getAllComments(feedId)
	if err != nil {
		return nil, err
	}
	return &entities.FeedWithComment{
		Feed:        feedDo.Feed,
		Comments:    comments,
		CatalogName: catalog.Name,
	}, nil
}

//PostFeed 发布帖子
func (s *feedSvc) PostFeed(req *entities.FeedCreateRequest, userId int64) error {
	if err := s.userNotBindMobile(userId); err != nil {
		return err
	}
	if !s.persistence.existedCatalogId(req.CatalogId) {
		return errors.New("板块不存在")
	}

	//创建帖子
	feed, err := s.persistence.createFeed(req, userId)
	if err != nil {
		return err
	}

	log.Info("feed:", feed)

	//发送消息
	err = s.Send(&events.PostFeed{
		FeedId:    feed.Id,
		UserId:    userId,
		Content:   feed.Content,
		Images:    feed.Images,
		CreatedAt: feed.CreatedAt,
	})
	if err != nil {
		log.Errorf("发送消息events.PostFeed{%v}失败", feed.Id)
	}
	return nil
}
func (s *feedSvc) DeleteFeed(feedId int64) error {
	return s.persistence.markFeedUserDeleted(feedId)
}

func (s *feedSvc) DingTalkNotifyComment(feedId int64, content string, user *entities.SimpleUser, a *entities.AuditResult) error {
	suggestion, label := a.Suggestion()
	return utils.SendDingTalkBot(s.DingTalkUrl, &utils.DingTalkMsg{
		Msgtype: "markdown",
		Markdown: utils.Markdown{
			Title: "评论审核不通过",
			Text: fmt.Sprintf(`### 评论审核不通过
- 用户: %s(%d)
- 帖子ID: %d
- 评论内容: %s
- 审核建议: %s
- 违规原因: %s
                `, user.Nickname, user.ID, feedId, content, suggestion, label),
		},
	})
}

//CommitFeedComment  提交帖子评论
func (s *feedSvc) CommitFeedComment(feedId int64, content string, userId int64) (err error) {
	if err := s.userNotBindMobile(userId); err != nil {
		return err
	}

	//调用机审接口 审核内容
	res, err := s.Audit.ScanText(content)
	if err != nil {
		return err
	}

	if !res.CheckPass() {
		log.Info("评论内容：", content)
		if user, exists, err := s.User.GetUserById(userId); exists && err == nil {
			err := s.DingTalkNotifyComment(feedId, content, user, res)
			if err != nil {
				log.Errorf("DingTalkNotifyComment: error %v", err)
			}
		}

		return errors.New("含有敏感信息，评论失败")
	}

	//检查帖子状态
	if err = s.checkFeedStatusById(feedId); err != nil {
		return err
	}

	err, comment := s.persistence.createComment(feedId, content, userId)
	if err != nil {
		return err
	}

	//更新帖子评论数
	err = s.persistence.updateFeedCommentCount(feedId)
	if err != nil {
		return err
	}

	//发送领域事件
	err = s.Send(&events.PostFeedComment{
		FeedCommentId: comment.Id,
		FeedId:        feedId,
		UserId:        userId,
		Content:       content,
		CreatedAt:     comment.CreatedAt,
	})

	return err
}

func (s *feedSvc) checkFeedStatusById(feedId int64) error {
	feed, err := s.persistence.getOneFeedById(feedId)
	if err != nil {
		return err
	}
	return s.checkFeedStatus(feed)
}

func (s *feedSvc) checkFeedStatus(feed *Feed) error {
	if feed == nil {
		return errors.New("帖子不存在")
	}
	if feed.Status == entities.FeedStatusNeedAudit {
		return errors.New("帖子内容在审核中，禁止操作")
	}
	if feed.Status == entities.FeedStatusAuditDeny {
		return errors.New("帖子内容审核不通过，禁止操作")
	}
	if feed.Status == entities.FeedStatusDeleted {
		return errors.New("帖子内容被举报已经删除，禁止操作")
	}
	return nil
}

func (s *feedSvc) patchMyFeedCreator(list []*entities.MyFeed) error {
	feeds := make([]*entities.Feed, 0, len(list))
	for i := range list {
		feeds = append(feeds, list[i].Feed)
	}

	return s.patchFeedCreator(feeds)
}

func (s *feedSvc) patchMyFeedsLikeAndComment(list []*entities.MyFeed, userId int64) error {
	feeds := make([]*entities.Feed, 0, len(list))
	for i := range list {
		feeds = append(feeds, list[i].Feed)
	}
	return s.patchFeedsLikeAndComment(feeds, userId)
}

func (s *feedSvc) patchFeedCreator(list []*entities.Feed) error {
	if len(list) == 0 {
		return nil
	}
	userIds := make([]int64, 0)
	userMap := make(map[int64][]*entities.Feed)
	for i := range list {
		feed := list[i]
		arr, ok := userMap[feed.CreatedBy]
		if !ok {
			arr = make([]*entities.Feed, 0)
			userIds = append(userIds, feed.CreatedBy)
		}
		arr = append(arr, feed)
		userMap[feed.CreatedBy] = arr
	}

	creatorList, err := s.getFeedCreatorByUserIds(userIds)
	if err != nil {
		return err
	}

	for i := range creatorList {
		creator := creatorList[i]
		feeds, ok := userMap[creator.UserId]
		if ok {
			for j := range feeds {
				feed := feeds[j]

				feed.Creator = &entities.FeedCreator{
					UserId:         creator.UserId,
					Nickname:       creator.Nickname,
					Avatar:         creator.Avatar,
					DonationPoints: creator.DonationPoints,
				}
				if feed.ShowMobile {
					feed.Creator.Mobile = creator.Mobile
				}
			}
		}
	}
	return nil
}

func (s *feedSvc) getFeedCreatorByUserIds(userIds []int64) (list []*entities.FeedCreator, err error) {
	users, err := s.User.GetListUserByIDs(userIds)
	if err != nil {
		return nil, err
	}

	donations, err := s.Donation.StatUsersDonations(userIds)
	if err != nil {
		return nil, err
	}

	creatorMap := make(map[int64]*entities.FeedCreator)
	for i := range users {
		user := users[i]
		userId := user.ID
		_, ok := creatorMap[userId]
		if !ok {
			creator := entities.FeedCreator{
				UserId:   user.ID,
				Nickname: user.Nickname,
				Mobile:   user.Mobile,
				Avatar:   user.Avatar,
			}
			creatorMap[userId] = &creator
			list = append(list, &creator)
		}
	}

	for i := range donations {
		donation := donations[i]
		userId := donation.UserId
		creator, ok := creatorMap[userId]
		if ok {
			creator.DonationPoints = donation.Points
		}
	}
	return list, nil
}

func (s *feedSvc) patchFeedsLikeAndComment(list []*entities.Feed, userId int64) error {
	if len(list) == 0 {
		return nil
	}
	feedIds := make([]int64, 0)
	feedMap := make(map[int64]*entities.Feed)
	for i := range list {
		feedId := list[i].Id
		feedIds = append(feedIds, feedId)
		feedMap[feedId] = list[i]
	}

	feeds, err := s.persistence.findUserMarkedInFeeds(userId, feedIds)
	if err != nil {
		return err
	}
	for i := range feeds {
		feedInfo := feeds[i]
		feedId := feedInfo.FeedId
		feed, ok := feedMap[feedId]
		if ok {
			feed.IsLike = feedInfo.IsLike
			feed.IsCommented = feedInfo.IsCommented
		}
	}
	return nil
}

func (s *feedSvc) recordShareEvent(userEventId int64, userId int64, feedId int64) error {
	return db.Transaction(func(session *xorm.Session) error {
		//保证幂等性
		exited, err := s.persistence.userEventIsExisted(userEventId)
		if err != nil {
			return err
		}
		if exited {
			return nil
		}

		err = s.persistence.recordShareEvent(userEventId, userId, feedId)
		if err != nil {
			return err
		}

		return s.persistence.updateFeedShareCount(feedId)
	})
}

func (s *feedSvc) userNotBindMobile(userId int64) error {
	user, b, err := s.User.GetUserById(userId)
	if err != nil {
		return err
	}
	if !b {
		return errors.New("用户不存在")
	}
	if user.Mobile == "" {
		return errors.New("必须绑定手机")
	}
	return nil
}
