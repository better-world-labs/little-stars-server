package entities

import "time"

type Catalog struct {
	Id   int    `json:"id"`
	Name string `json:"name"`
}

type FeedCreator struct {
	UserId         int64  `json:"userId"`
	Nickname       string `json:"nickname"`
	Mobile         string `json:"mobile"`
	Avatar         string `json:"avatar"`
	DonationPoints int    `json:"donationPoints"`
}

type Feed struct {
	Id        int64      `json:"id" xorm:"id pk autoincr"`
	CatalogId int        `json:"catalogId"`
	Content   string     `json:"content"`
	Images    []string   `json:"images"`
	Address   string     `json:"address"`
	Lon       float64    `json:"lon"`
	Lat       float64    `json:"lat"`
	Status    FeedStatus `json:"status"`

	ShowMobile bool `json:"-"`

	Creator   *FeedCreator `json:"creator" xorm:"-"`
	CreatedBy int64        `json:"-"`
	CreatedAt time.Time    `json:"createdAt"`

	CommentCount int `json:"commentCount"`
	ShareCount   int `json:"shareCount"`
	LikeCount    int `json:"likeCount"`

	IsCommented bool `json:"isCommented" xorm:"-"`
	IsLike      bool `json:"isLike" xorm:"-"`
}

type MyFeedType int

const (
	MyFeedTypeICreate  MyFeedType = 0
	MyFeedTypeIComment MyFeedType = 1
	MyFeedTypeILike    MyFeedType = 2
	MyFeedTypeIPrivate MyFeedType = 3
)

type FeedStatus int

const (
	FeedStatusNeedAudit   FeedStatus = 10
	FeedStatusAuditPass   FeedStatus = 20 //2x为帖子可正常正式状态
	FeedStatusAuditDeny   FeedStatus = 30 //3x为不可展示状态
	FeedStatusDeleted     FeedStatus = 31
	FeedStatusUserDeleted FeedStatus = 32
)

type MyFeed struct {
	*Feed `xorm:"extends"`
	Type  MyFeedType `json:"type"`
}

type FeedComment struct {
	Content      string    `json:"content"`
	CreatedAt    time.Time `json:"createdAt"`
	UserNickname string    `json:"userNickname"`
	UserAvatar   string    `json:"userAvatar"`
}

type FeedWithComment struct {
	*Feed
	CatalogName string         `json:"catalogName"`
	Comments    []*FeedComment `json:"comments"`
}

type FeedCreateRequest struct {
	CatalogId  int      `json:"catalogId" binding:"required"`
	Content    string   `json:"content" binding:"required"`
	Images     []string `json:"images"`
	Address    string   `json:"address"`
	Lon        float64  `json:"lon"`
	Lat        float64  `json:"lat"`
	ShowMobile bool     `json:"showMobile"`
}
