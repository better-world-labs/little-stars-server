package controller

import (
	"aed-api-server/internal/interfaces/entities"
	"aed-api-server/internal/interfaces/service"
	"aed-api-server/internal/pkg"
	"aed-api-server/internal/pkg/utils"
	"errors"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"gitlab.openviewtech.com/openview-pub/gopkg/route"
)

//go:inject-component
func NewFeedController() *FeedController {
	return &FeedController{}
}

type FeedController struct {
	FeedService service.IFeed       `inject:"-"`
	User        service.UserService `inject:"-"`
}

func (ctr *FeedController) MountAuthRouter(r *route.Router) {
	g := r.Group("/feeds")

	g.GET("/:id", ctr.getFeedDetail)
	g.GET("/mine", ctr.getMyFeeds)
	g.POST("/", ctr.postFeed)
	g.PUT("/set-private/:id", ctr.setPrivate)
	g.DELETE("/:id", ctr.deleteFeed)
	g.POST("/:id/comments", ctr.commitComment)
}

func (ctr *FeedController) MountNoAuthRouter(r *route.Router) {
	g := r.Group("/feeds")
	g.GET("/catalogs", ctr.getCatalogs)
	g.GET("/", ctr.getFeedsPage)
}

func (ctr *FeedController) getCatalogs(context *gin.Context) (interface{}, error) {
	return ctr.FeedService.GetCatalogList()
}

func (ctr *FeedController) getFeedsPage(context *gin.Context) (interface{}, error) {
	authorization := context.GetHeader(pkg.AuthorizationHeaderKey)
	var userId int64
	if authorization != "" {
		user, err := ctr.User.ParseInfoFromJwtToken(authorization)
		if err != nil {
			log.Infof("ParseInfoFromJwtToken error: %v", err)
		} else {
			userId = user.ID
			context.Set(pkg.AccountIDKey, userId)
			context.Set(pkg.AccountKey, user)
		}
	}

	type Req struct {
		CatalogId int    `form:"catalogId"`
		Cursor    string `form:"cursor"`
		Size      int    `form:"size"`
	}

	var req Req

	if err := context.ShouldBindQuery(&req); err != nil {
		return nil, err
	}

	feeds, cursor, err := ctr.FeedService.GetFeeds(req.CatalogId, req.Cursor, req.Size, userId)
	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"records":    feeds,
		"nextCursor": cursor,
	}, nil
}

func (ctr *FeedController) getFeedDetail(context *gin.Context) (interface{}, error) {
	type Req struct {
		Id int64 `uri:"id" binding:"required"`
	}

	var req Req
	if err := context.ShouldBindUri(&req); err != nil {
		return nil, err
	}

	return ctr.FeedService.GetFeedDetail(req.Id, context.MustGet(pkg.AccountIDKey).(int64))
}

func (ctr *FeedController) getMyFeeds(context *gin.Context) (interface{}, error) {
	type Req struct {
		Cursor string `form:"cursor"`
		Size   int    `form:"size"`
	}

	var req Req

	if err := context.ShouldBindQuery(&req); err != nil {
		return nil, err
	}

	feeds, cursor, err := ctr.FeedService.GetMyFeeds(req.Cursor, req.Size, context.MustGet(pkg.AccountIDKey).(int64))
	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"records":    feeds,
		"nextCursor": cursor,
	}, nil
}

func (ctr *FeedController) postFeed(context *gin.Context) (interface{}, error) {
	var req entities.FeedCreateRequest
	if err := context.ShouldBindJSON(&req); err != nil {
		return nil, err
	}

	err := ctr.FeedService.PostFeed(&req, context.MustGet(pkg.AccountIDKey).(int64))
	if err != nil {
		return nil, err
	}
	return nil, nil
}

func (ctr *FeedController) commitComment(context *gin.Context) (interface{}, error) {
	type Req struct {
		FeedId  int64  `uri:"id" binding:"required"`
		Content string `json:"content"`
	}

	var req Req

	if err := context.ShouldBindUri(&req); err != nil {
		return nil, err
	}
	if err := context.ShouldBindJSON(&req); err != nil {
		return nil, err
	}

	if req.Content == "" {
		return nil, errors.New("content 是必填字段")
	}

	err := ctr.FeedService.CommitFeedComment(req.FeedId, req.Content, context.MustGet(pkg.AccountIDKey).(int64))
	if err != nil {
		return nil, err
	}

	return nil, nil
}

func (ctr *FeedController) deleteFeed(ctx *gin.Context) (interface{}, error) {
	feedId, err := utils.GetContextPathParamInt64(ctx, "id")
	if err != nil {
		return nil, err
	}

	err = ctr.FeedService.DeleteFeed(feedId)
	if err != nil {
		return nil, err
	}

	return nil, nil
}

func (ctr *FeedController) setPrivate(ctx *gin.Context) (interface{}, error) {
	panic("not implements")
}
