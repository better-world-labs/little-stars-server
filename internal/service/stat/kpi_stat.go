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
)

func Init() {
	interfaces.S.Stat = &service{}
}

type service struct{}

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
	})

	if err != nil {
		return nil, err
	}

	return &entities.KpiStatItem{
		UserCount:            arrRst[0].(int) + arrRst[1].(int),
		MiniProgramUserCount: arrRst[0].(int),
		OffiaccountUserCount: arrRst[1].(int),
		LearntUserCount:      arrRst[2].(int),
		DailyUserCount:       arrRst[3].(int),
		UserNextDayUseRatio:  arrRst[4].(float64),
		DeviceCount:          arrRst[5].(int),
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
