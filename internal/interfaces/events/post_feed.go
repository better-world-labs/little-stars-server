package events

import (
	"aed-api-server/internal/pkg/domain/emitter"
	"encoding/json"
	"time"
)

type PostFeed struct {
	FeedId    int64     `json:"feedId"`
	UserId    int64     `json:"userId"`
	Content   string    `json:"content"`
	Images    []string  `json:"images"`
	CreatedAt time.Time `json:"createdAt"`
}

func (e *PostFeed) GetUserId() int64 {
	return e.UserId
}

func (e *PostFeed) Decode(bytes []byte) (emitter.DomainEvent, error) {
	var f PostFeed
	err := json.Unmarshal(bytes, &f)
	return &f, err
}

func (e *PostFeed) Encode() ([]byte, error) {
	return json.Marshal(e)
}

type PostFeedComment struct {
	FeedCommentId int64     `json:"feedCommentId"` //帖子评论ID
	FeedId        int64     `json:"feedId"`        //帖子ID
	UserId        int64     `json:"userId"`        //用户ID
	Content       string    `json:"content"`       //评论内容
	CreatedAt     time.Time `json:"createdAt"`     //评论时间
}

func (e *PostFeedComment) GetUserId() int64 {
	return e.UserId
}

func (e *PostFeedComment) Decode(bytes []byte) (emitter.DomainEvent, error) {
	var f PostFeedComment
	err := json.Unmarshal(bytes, &f)
	return &f, err
}

func (e *PostFeedComment) Encode() ([]byte, error) {
	return json.Marshal(e)
}
