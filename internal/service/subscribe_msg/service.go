package subscribe_msg

import (
	"aed-api-server/internal/interfaces"
	"aed-api-server/internal/interfaces/entities"
	"aed-api-server/internal/pkg/db"
	"encoding/json"
	"errors"
	"fmt"
	log "github.com/sirupsen/logrus"
	"strings"
	"time"
)

const (
	WxNotifyTemplateTaskExpired        = "tkQRItF-ipKOym08GEqe0c7rFtBolvezpjeAeG0LN1I"
	WxNotifyTemplateCouponExpired      = "Hp0_c7lQ8kicztthTNTmDTjEVqQl2YraB1WqfaqDgWY"
	WxNotifyTemplateAdventureAdventure = "Wt2SuY8pB7DoIuAEOS7UgPqAN8S51Vs_H2fi69YD7-8"
)

var (
	templateIdMap = map[entities.SubscribeMessageKey]string{
		entities.SMkPointsExpiring: WxNotifyTemplateCouponExpired,
		entities.SMkWalkExpiring:   WxNotifyTemplateTaskExpired,
		entities.SMkGamePoints:     WxNotifyTemplateAdventureAdventure,
	}

	templateEl = map[string]string{
		WxNotifyTemplateAdventureAdventure: "notice_1_click",
		WxNotifyTemplateTaskExpired:        "notice_2_click",
		WxNotifyTemplateCouponExpired:      "notice_3_click",
	}
)

var msgKeyMap = getMsgKeyMapFromTemplateIdMap(templateIdMap)

func getMsgKeyMapFromTemplateIdMap(templateIdMap map[entities.SubscribeMessageKey]string) map[string]entities.SubscribeMessageKey {
	_map := map[string]entities.SubscribeMessageKey{}
	for key, s := range templateIdMap {
		_map[s] = key
	}
	return _map
}

type svc struct{}

//go:inject-component
func NewSubscribeMsgService() *svc {
	return &svc{}
}

func (*svc) Report(userId int64, key string, templates []*entities.SubscribeTemplateSetting, setting *entities.SubscriptionsSetting) error {
	//记录事件
	interfaces.S.User.RecordUserEvent(userId, entities.GetUserEventTypeOfReport(key), templates, setting)

	//更新发消息的记录
	err := updateSendMsgTickets(userId, templates, setting)
	return err
}

func (*svc) GetLastReport(userId int64, key string) (templates []*entities.SubscribeTemplateSetting, setting *entities.SubscriptionsSetting, reportAt *time.Time, err error) {
	//查询最近时间
	event, err := interfaces.S.User.GetLastUserEventByType(userId, entities.GetUserEventTypeOfReport(key))
	if err != nil {
		return nil, nil, nil, err
	}

	if event == nil {
		return nil, nil, nil, nil
	}

	reportAt = &event.CreatedAt

	if event.EventParams == nil {
		return nil, nil, reportAt, nil
	}

	json1, _ := json.Marshal(event.EventParams[0])
	json2, _ := json.Marshal(event.EventParams[1])
	//
	err = json.Unmarshal([]byte(json1), &templates)
	if err != nil {
		return nil, nil, nil, err
	}

	err = json.Unmarshal([]byte(json2), &setting)
	if err != nil {
		return nil, nil, nil, err
	}
	return templates, setting, reportAt, nil
}
func (*svc) Send(userId int64, openId string, msgKey entities.SubscribeMessageKey, params interface{}) (bool, error) {
	templateId, ok := templateIdMap[msgKey]
	if !ok {
		return false, errors.New("SubscribeMessageKey invalid:" + string(msgKey))
	}

	suc, err := useTicket(userId, msgKey)
	if err != nil {
		return false, err
	}
	if !suc {
		return false, nil
	}
	el := templateEl[templateId]
	rst, err := interfaces.S.Wx.SendSubscribeMsg(msgKey, openId, templateId, el, params)
	if err != nil {
		return false, err
	}

	if rst.ErrCode == 43101 {
		clearTicket(userId, msgKey)
	}

	if rst.ErrCode != 0 {
		return false, nil
	}

	return true, nil
}

type UserSubscribeMessage struct {
	UserId               int64
	TicketPointsExpiring int
	TicketWalkExpiring   int
}

func updateSendMsgTickets(userId int64, templates []*entities.SubscribeTemplateSetting, setting *entities.SubscriptionsSetting) error {
	_switchMap := make(map[entities.SubscribeMessageKey]bool, 0)
	if len(templates) == 0 {
		for key := range templateIdMap {
			templateId := templateIdMap[key]
			templates = append(templates, &entities.SubscribeTemplateSetting{
				TemplateId: templateId,
				Status:     "reject",
			})
		}
	}

	for _, tpl := range templates {
		msgKey, ok := msgKeyMap[tpl.TemplateId]
		if ok {
			_switchMap[msgKey] = tpl.Status == "accept"
		}
	}

	if !setting.MainSwitch {
		for key := range _switchMap {
			_switchMap[key] = false
		}
	} else {
		for _, tpl := range setting.Templates {
			msgKey, ok := msgKeyMap[tpl.TemplateId]
			if ok {
				_, ok = _switchMap[msgKey]
				if ok && tpl.Status != "accept" {
					_switchMap[msgKey] = false
				}
			}
		}
	}

	boolToString := func(input bool) string {
		if input {
			return "1"
		} else {
			return "0"
		}
	}

	switchSql := func(flag bool, fieldName string) string {
		if flag {
			return fmt.Sprintf("%s = %s +1", fieldName, fieldName)
		} else {
			return fmt.Sprintf("%s = 0", fieldName)
		}
	}

	fields := make([]string, 0)
	inserts := make([]string, 0)
	updates := make([]string, 0)
	for key, switcher := range _switchMap {
		tableFieldName := getTableFieldFromSubscribeMessageKey(key)
		fields = append(fields, tableFieldName)
		inserts = append(inserts, boolToString(switcher))
		updates = append(updates, switchSql(switcher, tableFieldName))
	}

	sql := fmt.Sprintf(`insert into user_subscribe_message(user_id, %s) values(?, %s)on duplicate key update %s`,
		strings.Join(fields, ","),
		strings.Join(inserts, ","),
		strings.Join(updates, ","),
	)

	_, err := db.Exec(sql, userId)
	if err != nil {
		return err
	}
	return nil
}

func getTableFieldFromSubscribeMessageKey(msgKey entities.SubscribeMessageKey) string {
	return string(msgKey)
}

func ListUsersSubscriptionValidTicket(msgKey entities.SubscribeMessageKey) ([]*entities.NotifiedSubscription, error) {
	tableFieldName := getTableFieldFromSubscribeMessageKey(msgKey)
	var arr []*entities.NotifiedSubscription
	err := db.Table("user_subscribe_message").
		Alias("m").
		Join("LEFT", []string{"account", "a"}, "m.user_id = a.id").Cols("m.*", "a.openid").
		Where(fmt.Sprintf("%s > 0", tableFieldName)).
		Find(&arr)
	return arr, err
}

func reduceWalkTicket(userId int64) error {
	_, err := db.Exec(`
		update user_subscribe_message set walk_expiring = walk_expiring - 1 where user_id = ?
    `, userId)

	return err
}

//使用一次发送通知的机会
func useTicket(userId int64, msgKey entities.SubscribeMessageKey) (bool, error) {
	tableFieldName := getTableFieldFromSubscribeMessageKey(msgKey)

	rst, err := db.Exec(fmt.Sprintf(`
		update user_subscribe_message
			set %s = %s -1
		where %s > 0
		and 
			user_id = ?
	`, tableFieldName, tableFieldName, tableFieldName), userId)
	if err != nil {
		return false, err
	}

	affected, err := rst.RowsAffected()
	if err != nil {
		return false, err
	}
	return 1 == affected, nil
}

func clearTicket(userId int64, msgKey entities.SubscribeMessageKey) {
	_, err := db.Exec(fmt.Sprintf(`update user_subscribe_message set %s = 0 where user_id = ?`, msgKey), userId)
	if err != nil {
		log.Error("clearTicket userId=", userId, " err:", err)
	}
}

func notHasSendMsgTicket(msgKey entities.SubscribeMessageKey, userId int64) bool {
	tableFieldName := getTableFieldFromSubscribeMessageKey(msgKey)
	count, err := db.SQL(
		fmt.Sprintf(`select %s from user_subscribe_message where user_id =?`, tableFieldName),
		userId,
	).Count()

	if err != nil {
		log.Error("get ticket count error for msgKey:", msgKey)
		return false
	}
	return count <= 0
}
