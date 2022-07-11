package stat

import (
	"aed-api-server/internal/interfaces"
	"aed-api-server/internal/interfaces/entities"
	"aed-api-server/internal/pkg/db"
	"aed-api-server/internal/pkg/utils"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)

type service struct{}

//go:inject-component
func NewService() *service {
	return &service{}
}

func (*service) DoKipStat() (*entities.KpiStatItem, error) {
	arrRst, err := utils.PromiseAll(func() (interface{}, error) {
		return statMiniProgramUserCount()
	}, func() (interface{}, error) {
		return getOffiaccountUserCount()
	}, func() (interface{}, error) {
		return getLearntUserCount()
	}, func() (interface{}, error) {
		return getTodayUseUserCount()
	}, func() (interface{}, error) {
		return getUserNextDayUseRatio()
	}, func() (interface{}, error) {
		return getDeviceCount()
	}, func() (interface{}, error) {
		return interfaces.S.User.StatUser()
	})

	if err != nil {
		return nil, err
	}

	userStat := arrRst[6].(entities.UserStat)

	return &entities.KpiStatItem{
		UserCount:            arrRst[0].(int) + arrRst[1].(int),
		MiniProgramUserCount: arrRst[0].(int),
		OffiaccountUserCount: arrRst[1].(int),
		LearntUserCount:      arrRst[2].(int),
		DailyUserCount:       arrRst[3].(int),
		UserNextDayUseRatio:  arrRst[4].(float64),
		DeviceCount:          arrRst[5].(int),
		RegUserCount:         userStat.TotalCount,
		MobileUserCount:      userStat.MobileCount,
	}, nil
}

func getDeviceCount() (int, error) {
	count, err := db.SQL(`select count(1) from device where created < unix_timestamp(CURRENT_DATE)`).Count()
	return int(count), err
}

func getUserNextDayUseRatio() (float64, error) {
	type Rst struct {
		YesterdayCount int
		TodayCount     int
	}

	var rst Rst

	_, err := db.SQL(`
		select
			count(distinct a.user_id) as yesterday_count,
			count(distinct b.user_id) as today_count
		from user_event_record as a
		left join (
			select
				distinct user_id
			from user_event_record
			WHERE
				created_at > TIMESTAMPADD(day,-1,CURRENT_DATE)
		) as b
			on b.user_id = a.user_id
		WHERE
			a.created_at BETWEEN TIMESTAMPADD(day,-2,CURRENT_DATE) and TIMESTAMPADD(day,-1,CURRENT_DATE)
	`).Get(&rst)

	if err != nil {
		return 0, err
	}
	if rst.YesterdayCount == 0 {
		return 0, nil
	}
	return float64(rst.TodayCount) / float64(rst.YesterdayCount) * 100, nil
}

func getTodayUseUserCount() (int, error) {
	count, err := db.SQL(`
		select count(distinct user_id) as count 
		from user_event_record
		WHERE
			created_at BETWEEN TIMESTAMPADD(day,-1,CURRENT_DATE) and CURRENT_DATE
	`).Count()
	return int(count), err
}

func getLearntUserCount() (int, error) {
	count, err := db.SQL(`
		select
			count(distinct user_id) as count
		from (
			select distinct
				user_id
			from project_user
			where
				created_at < CURRENT_DATE
			Union 
			select distinct
				user_id
			from project_course_user
			where
				created_at < CURRENT_DATE
		) as a
	`).Count()
	return int(count), err
}

func statMiniProgramUserCount() (int, error) {
	count, err := db.SQL(`
		select
			count(distinct open_id) as count 
		from generalize_trace
		where
			created_at < CURRENT_DATE
	`).Count()
	return int(count), err
}

func getOffiaccountUserCount() (int, error) {
	token, err := getWechatToken()
	if err != nil {
		return 0, err
	}

	res, err := http.Get(fmt.Sprintf(`https://api.weixin.qq.com/cgi-bin/user/get?access_token=%s`, token))
	all, err := ioutil.ReadAll(res.Body)
	rst := map[string]interface{}{}
	err = json.Unmarshal(all, &rst)

	total, ok := rst["total"]
	if ok {
		return int(total.(float64)), nil
	}
	if rst["errcode"] != 0 {
		return 0, errors.New(rst["errmsg"].(string))
	}
	return rst["total"].(int), nil
}

func getWechatToken() (string, error) {
	config := interfaces.GetConfig()
	res, err := http.Get(fmt.Sprintf(`https://api.weixin.qq.com/cgi-bin/token?grant_type=client_credential&appid=%s&secret=%s`,
		config.WechatOffiaccountAppid,
		config.WechatOffiaccountSecret),
	)
	if err != nil {
		return "", err
	}

	all, err := ioutil.ReadAll(res.Body)
	rst := map[string]interface{}{}
	err = json.Unmarshal(all, &rst)
	if err != nil {
		return "", err
	}

	token, ok := rst["access_token"]
	if ok {
		return token.(string), nil
	}

	if rst["errcode"] != 0 {
		return "", errors.New(rst["errmsg"].(string))
	}
	return rst["access_token"].(string), nil
}

func (*service) StatPointsTop() (*entities.UserPointsTop, error) {
	ranks, err := getTopN(10)
	if err != nil {
		return nil, err
	}

	text := getTopNMarkdownText(ranks)

	return &entities.UserPointsTop{
		List: ranks,
		Text: text,
	}, nil
}

func getTopN(topN int) ([]*entities.UserPointsRank, error) {
	ranks := make([]*entities.UserPointsRank, 0)
	err := db.SQL(`
		select
			user_id,
			sum(points) as points_amount,
			count(points) as points_count
		from point_flow
		WHERE
			created_at >= TIMESTAMPADD(day,-1,CURRENT_DATE)
			and created_at < CURRENT_DATE
			and status = 1
			and points > 0
		group by user_id
		order by points_amount desc, points_count desc
		limit ?
	`, topN).Find(&ranks)
	if err != nil {
		return nil, err
	}
	return ranks, nil
}

func getTopNMarkdownText(ranks []*entities.UserPointsRank) string {
	var title = "### 昨日积分排行\n"
	if len(ranks) == 0 {
		return title
	}

	userIds := make([]int64, 0)
	for i := range ranks {
		userIds = append(userIds, ranks[i].UserId)
	}
	ds, err := interfaces.S.User.GetListUserByIDs(userIds)
	if err != nil {
		return title
	}

	m := make(map[int64]*entities.SimpleUser)
	for i := range ds {
		user := ds[i]
		m[user.ID] = user
	}

	userRankStrList := make([]string, 0)
	for i := range ranks {
		rank := ranks[i]
		user, ok := m[rank.UserId]
		if !ok {
			user = &entities.SimpleUser{
				ID:       rank.UserId,
				Nickname: "未知用户",
			}
		}
		userRankStrList = append(userRankStrList, fmt.Sprintf("#### %d). %d(%d次，%s=%d)  \n",
			i+1, rank.PointsAmount, rank.PointsCount, user.Nickname, user.ID),
			userPointsDis(rank.UserId),
			"\n",
		)
	}

	return title + strings.Join(userRankStrList, "")
}

func userPointsDis(userId int64) string {
	type Rst struct {
		Name   string
		Points int
		Count  int
	}

	list := make([]*Rst, 0)
	err := db.SQL(`
		select
			b.name,
			count(1) as count,
			a.points
		from point_flow as a
		left join point_event_define as b on b.points_event_type = a.points_event_type
		where 
			a.user_id = ?
			and a.points > 0
			and a.created_at > TIMESTAMPADD(day,-1,CURRENT_DATE)
			and a.created_at < CURRENT_DATE
		group by a.points_event_type
	`, userId).Find(&list)
	if err != nil {
		return ""
	}
	textList := make([]string, 0)
	for i := range list {
		r := list[i]
		textList = append(textList, fmt.Sprintf("- %s(%d分,%d次)\n", r.Name, r.Points, r.Count))
	}

	return strings.Join(textList, "")
}
