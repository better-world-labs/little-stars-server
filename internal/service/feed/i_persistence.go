package feed

import (
	"aed-api-server/internal/interfaces/entities"
	"time"
)

type Feed struct {
	*entities.Feed `xorm:"extends"`
	Top            bool
}

type Comment struct {
	Id        int64 `json:"id" xorm:"id pk autoincr"`
	FeedId    int64
	Content   string
	Status    entities.FeedStatus
	CreatedAt time.Time
	CreatedBy int64
}

type Like struct {
	FeedId    int64
	UserId    int64
	IsLike    bool
	UpdatedAt time.Time
}

type Report struct {
	FeedId    int64
	Content   string
	Type      string
	CreatedAt time.Time
	CreatedBy int64
}

type UserMarkedFeedInfo struct {
	FeedId      int64
	IsLike      bool
	IsCommented bool
}

type iPersistence interface {
	getCatalogList() ([]*entities.Catalog, error)
	getFeeds(beforeId int64, catalogId int, size int, userId int64) ([]*entities.Feed, error)
	createFeed(req *entities.FeedCreateRequest, userId int64) (feed Feed, err error)
	createComment(fee int64, content string, userId int64) (error, *Comment)
	getMyFeeds(beforeId int64, size int, userId int64) ([]*entities.MyFeed, error)
	updateFeedCommentCount(feedId int64) error
	getOneFeedById(feedId int64) (feed *Feed, err error)
	updateFeedLike(userId int64, feedId int64, like bool, opTime time.Time) error
	updateFeedLikeCount(feedId int64) error
	updateFeedShareCount(feedId int64) error
	existedCatalogId(catalogId int) bool
	updateFeedStatus(feedId int64, status entities.FeedStatus) error
	findUserMarkedInFeeds(userId int64, feedIds []int64) (list []*UserMarkedFeedInfo, err error)
	getAllComments(feedId int64) ([]*entities.FeedComment, error)
	getMyFeedsCount(userId int64) (int, error)
	userEventIsExisted(userEventId int64) (bool, error)
	recordShareEvent(userEventId int64, userId int64, feedId int64) error
	getCatalogById(catalogId int) (*entities.Catalog, error)
	getFeedUsers(feedId int64) ([]int64, error)
	markFeedUserDeleted(feedId int64) error
}
