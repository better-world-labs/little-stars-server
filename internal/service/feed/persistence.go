package feed

import (
	"aed-api-server/internal/interfaces/entities"
	"aed-api-server/internal/pkg/db"
	"fmt"
	"github.com/go-sql-driver/mysql"
	"github.com/go-xorm/xorm"
	log "github.com/sirupsen/logrus"
	"strconv"
	"strings"
	"time"
)

type persistence struct{}

func (p *persistence) getCatalogList() (list []*entities.Catalog, err error) {
	err = db.Table("feed_catalog").Find(&list)
	return list, err
}

//getFeeds 获取贴子流
//未审核的文章自己可见
func (p *persistence) getFeeds(beforeId int64, catalogId int, size int, userId int64) (list []*entities.Feed, err error) {
	beforeSql := ""
	if beforeId > 0 {
		beforeSql = fmt.Sprintf(`and id < %v and top = 0`, beforeId)
	}
	catalogSql := ""
	if catalogId > 0 {
		catalogSql = fmt.Sprintf(`and catalog_id = %v`, catalogId)
	}

	sql := fmt.Sprintf(`
		select 
			*
		from feed
		where
			((created_by = ? and status = ?) or status = ?)
			%s
			%s
		order by top desc, id desc
		limit ?
	`, beforeSql, catalogSql)

	err = db.SQL(sql, userId, entities.FeedStatusNeedAudit, entities.FeedStatusAuditPass, size).Find(&list)
	return list, err
}

func (p *persistence) getMyFeeds(beforeId int64, size int, userId int64) (list []*entities.MyFeed, err error) {
	beforeSql := ""
	if beforeId > 0 {
		beforeSql = fmt.Sprintf(`and a.id < %v`, beforeId)
	}

	sql := fmt.Sprintf(`
		select 
			a.*,
			case
				when a.created_by = :userId then 0
				when b.created_by = :userId then 1
				else 2
			end as type
		from feed as a
		left join feed_comment as b
			on b.feed_id = a.id
			and b.created_by = :userId
		left join feed_like as c
			on c.feed_id = a.id
			and c.user_id = :userId
			and c.is_like = 1
		where
			(
				(a.created_by = :userId and a.status in (10,20) )
				or (b.created_by = :userId and a.status = 20)
				or (c.user_id = :userId and a.status = 20)
			)
			%s
		group by a.id
		order by a.id desc
		limit ?
	`, beforeSql)

	sql = strings.ReplaceAll(sql, ":userId", strconv.FormatInt(userId, 10))

	err = db.SQL(sql, size).Find(&list)

	return list, err
}

func (p *persistence) createFeed(req *entities.FeedCreateRequest, userId int64) (Feed, error) {
	var feed = Feed{
		Feed: &entities.Feed{
			Content:    req.Content,
			CatalogId:  req.CatalogId,
			Images:     req.Images,
			Address:    req.Address,
			Lat:        req.Lat,
			Lon:        req.Lon,
			CreatedAt:  time.Now(),
			CreatedBy:  userId,
			ShowMobile: req.ShowMobile,
			Status:     entities.FeedStatusNeedAudit,
		},
	}
	_, err := db.Table("feed").Insert(&feed)
	return feed, err
}

func (p *persistence) createComment(feedId int64, content string, userId int64) (error, *Comment) {
	var comment = Comment{
		FeedId:    feedId,
		Content:   content,
		Status:    entities.FeedStatusAuditPass,
		CreatedAt: time.Now(),
		CreatedBy: userId,
	}

	_, err := db.Table("feed_comment").Insert(&comment)
	return err, &comment
}

func (p *persistence) updateFeedCommentCount(feedId int64) error {
	_, err := db.Exec(`
		update feed
		set comment_count = (select count(1) from feed_comment where feed_id = ? and status = 20)
		where
			id = ?
	`, feedId, feedId)
	return err
}

func (p *persistence) getOneFeedById(feedId int64) (feed *Feed, err error) {
	feed = new(Feed)
	existed, err := db.Table("feed").Where("id=?", feedId).Get(feed)
	if err != nil {
		return nil, err
	}
	if !existed {
		return nil, nil
	}
	return feed, nil
}

func (p *persistence) updateFeedLike(userId int64, feedId int64, like bool, opTime time.Time) error {
	return db.Transaction(func(db *xorm.Session) error {
		tmpLike := Like{
			UserId:    userId,
			FeedId:    feedId,
			IsLike:    like,
			UpdatedAt: opTime,
		}

		_, err := db.Table("feed_like").Insert(&tmpLike)
		if err != nil {
			mysqlErr, ok := err.(*mysql.MySQLError)
			if ok && mysqlErr.Number == 1062 { //如果重复
				_, err = db.Exec(`
					update feed_like 
					set 
						is_like = ?,
						updated_at = ?
					where
						user_id = ? 
						and feed_id = ?
						and is_like != ?
						and updated_at < ?
					`, like, opTime, userId, feedId, like, opTime)
				return err
			}
			return err
		}
		return err
	})
}

func (p *persistence) updateFeedLikeCount(feedId int64) error {
	return db.Transaction(func(db *xorm.Session) error {
		_, err := db.Exec(`
			update feed
			set like_count = (select count(1) from feed_like where feed_id = ? and is_like = 1)
			where
				id = ?
		`, feedId, feedId)
		return err
	})
}

func (p *persistence) userEventIsExisted(userEventId int64) (bool, error) {
	return db.Table("feed_share_event").Where("event_id = ?", userEventId).Exist()
}

func (p *persistence) recordShareEvent(userEventId int64, userId int64, feedId int64) error {
	return db.Transaction(func(session *xorm.Session) error {
		_, err := session.Table("feed_share_event").Insert(map[string]interface{}{
			"event_id": userEventId,
			"feed_id":  feedId,
			"user_id":  userId,
		})
		return err
	})
}

func (p *persistence) updateFeedShareCount(feedId int64) error {
	return db.Transaction(func(db *xorm.Session) error {
		_, err := db.Exec(`
			update feed
			set share_count = (select count(1) from feed_share_event where feed_id = ?)
			where
				id = ?
		`, feedId, feedId)
		return err
	})
}

func (p *persistence) existedCatalogId(catalogId int) bool {
	exist, err := db.Table("feed_catalog").Exist(&entities.Catalog{Id: catalogId})
	if err != nil {
		log.Error("err in existedCatalogId:", err)
	}
	return exist
}

func (p *persistence) updateFeedStatus(feedId int64, status entities.FeedStatus) error {
	_, err := db.Table("feed").Where("id = ?", feedId).Update(map[string]interface{}{
		"status": status,
	})
	return err
}

func (p *persistence) findUserMarkedInFeeds(userId int64, feedIds []int64) (list []*UserMarkedFeedInfo, err error) {
	sql := `
		select
			a.id as feed_id,
			ifnull(b.is_like,0) as is_like,
			ifnull(c.is_commented,0) as is_commented
		from feed as a
		left join (
			SELECT
				feed_id,
				is_like
			FROM
				feed_like 
			WHERE
				feed_id in (:feedIds)
				AND user_id = :userId
		) as b on b.feed_id = a.id
		left join (
			SELECT
				feed_id,
				count(1)>0 as is_commented
			FROM
				feed_comment 
			WHERE
				feed_id in (:feedIds)
				AND created_by = :userId
				AND STATUS = 20
			group by feed_id
		) as c on c.feed_id = a.id
		where
			a.id in (:feedIds)
`
	err = db.Sqlx(sql, db.NameMap{
		"feedIds": feedIds,
		"userId":  userId,
	}).Find(&list)
	return list, err
}

func (p *persistence) getAllComments(feedId int64) (list []*entities.FeedComment, err error) {
	err = db.SQL(`
		select
			a.content,
			a.created_at,
			b.nickname as user_nickname,
			b.avatar as user_avatar
		from feed_comment as a
		left join account as b on b.id = a.created_by
		where
			a.feed_id = ?
			and a.status = 20
		order by a.id desc
    `, feedId).Find(&list)
	return list, err
}

func (p *persistence) getMyFeedsCount(userId int64) (int, error) {
	// 我创建 + 我点赞（不包含我创建）+ 我评论（不包含我创建和我点赞的）
	sql := `
		select 
		(
			select
				count(1)
			from feed
			where
				created_by = :userId
				and status in (10, 20)
		)
		+
		(
			select
				count(1)
			from feed_like as a
			inner join feed as b
				on b.id = a.feed_id
				and b.status = 20
				and b.created_by != :userId
			where
				a.user_id = :userId
				and a.is_like = 1
		)
		+
		(
			select
				count(distinct a.feed_id)
			from feed_comment  as a
			inner join feed as b
				on b.id = a.feed_id
				and b.status = 20
				and b.created_by != :userId
			left join feed_like as c
				on c.feed_id = a.id
				and c.user_id = :userId
			where
				a.created_by = :userId
				and a.status = 20
				and ifnull(c.is_like,0) = 0
		)
	`
	count, err := db.Sqlx(sql, db.NameMap{
		"userId": userId,
	}).Count()
	return int(count), err
}

func (p *persistence) getCatalogById(catalogId int) (catalog *entities.Catalog, err error) {
	catalog = new(entities.Catalog)
	existed, err := db.Table("feed_catalog").Where("id = ?", catalogId).Get(catalog)
	if err == nil && !existed {
		catalog = nil
	}
	return catalog, err
}

func (p *persistence) getFeedUsers(feedId int64) ([]int64, error) {
	type Record struct {
		UserId int64
	}

	list := make([]*Record, 0)

	err := db.SQL(`
		select created_by as user_id from feed where id = ?
		union
		select created_by as user_id from feed_comment where feed_id = ?
	`, feedId, feedId).Find(&list)

	if err != nil {
		return nil, err
	}

	userIds := make([]int64, 0, len(list))
	for i := range list {
		userIds = append(userIds, list[i].UserId)
	}
	return userIds, nil
}

func (p *persistence) markFeedUserDeleted(feedId int64) error {
	_, err := db.Exec("update feed set status = ? where id = ?", entities.FeedStatusUserDeleted, feedId)
	return err
}
