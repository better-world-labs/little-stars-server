package service

import "aed-api-server/internal/interfaces/entities"

type IFeed interface {
	//GetCatalogList 获取帖子板块列表
	GetCatalogList() ([]*entities.Catalog, error)

	//GetFeeds 获取帖子流
	GetFeeds(catalogId int, cursor string, size int, userId int64) (list []*entities.Feed, nextCursor string, err error)

	//GetMyFeeds 获取和我相关的帖子
	GetMyFeeds(cursor string, size int, userId int64) (list []*entities.MyFeed, nextCursor string, err error)

	//GetMyFeedsCount 获取和我相关的帖子的数量
	GetMyFeedsCount(userId int64) (count int, err error)

	// GetFeedById 根据ID读取帖子简单信息
	GetFeedById(feedId int64) (feed *entities.Feed, err error)

	//GetFeedDetail 获取帖子详情
	GetFeedDetail(feedId int64, userId int64) (feed *entities.FeedWithComment, err error)

	//PostFeed 发布帖子
	PostFeed(req *entities.FeedCreateRequest, userId int64) error

	//DeleteFeed 删除帖子
	DeleteFeed(feedId int64) error

	//CommitFeedComment  提交帖子评论
	CommitFeedComment(feedId int64, content string, userId int64) error
}
