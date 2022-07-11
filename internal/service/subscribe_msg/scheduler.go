package subscribe_msg

import (
	"aed-api-server/internal/interfaces"
	"aed-api-server/internal/interfaces/entities"
	"aed-api-server/internal/interfaces/facility"
	"aed-api-server/internal/pkg/utils"
	"fmt"
	log "github.com/sirupsen/logrus"
	"time"
)

const pointsExpiredTime = 2 * 60 * 60 * 1000 //用户积分即将在  ${pointsExpiredTime} 过期

var JobLockKeySendPointsExpiringMsg = facility.GetSendMsgJobLockKey(entities.SMkPointsExpiring)
var JobLockKeySendWalkConvertExpiringMsg = facility.GetSendMsgJobLockKey(entities.SMkWalkExpiring)

func getSendTimeKey(key entities.SubscribeMessageKey) string {
	return fmt.Sprintf("msg-time:%s", key)
}

//Cron 安装定时任务
func (*svc) Cron(run facility.RunFuncOnceAt) {
	//8:00-21:00
	run("0 8-21 * * *", JobLockKeySendPointsExpiringMsg, sendPointsExpiringMsg)

	//20:00
	run("0 20 * * *", JobLockKeySendWalkConvertExpiringMsg, sendWalkConvertExpiringMsg)
}

func sendPointsExpiringMsg() {
	//1. 获取所有2小时内有积分将要过期的用户
	userIds, err := interfaces.S.Points.GetAllPointsExpiringUserIds(pointsExpiredTime)
	if err != nil {
		log.Error("GetAllPointsExpiringUserIds error:", err)
		return
	}

	for i := 0; i < len(userIds); i++ {
		userId := userIds[i]
		if hasSendMsgToday(entities.SMkPointsExpiring, userId) || notHasSendMsgTicket(entities.SMkPointsExpiring, userId) {
			continue
		}
		doSendPointsExpiringMsg(userId)
	}
}

func doSendPointsExpiringMsg(userId int64) {
	openId, err := interfaces.S.User.GetUserOpenIdById(userId)
	if err != nil {
		log.Error("getUserOpenIdById", err)
		return
	}
	var points int
	var expiredTime time.Time

	points, expiredTime, err = interfaces.S.Points.StatExpiringPoints(userId)
	if err != nil {
		log.Error("StatExpiringPoints err:", err, ",userId=", userId)
		return
	}
	if points == 0 {
		log.Info("not has expiring points")
		return
	}

	suc, err := interfaces.S.SubscribeMsg.Send(userId, openId, entities.SMkPointsExpiring, map[string]interface{}{
		"thing1": map[string]interface{}{
			"value": fmt.Sprintf("你还有%d积分未领取", points),
		},
		"thing2": map[string]interface{}{
			"value": "请及时领取哦，过期失效~",
		},
		"time3": map[string]interface{}{
			"value": expiredTime.Format("15:04:05"),
		},
	})
	if err != nil {
		log.Error("SubscribeMsg.Send error:", err, ",userId=", userId)
		return
	}
	if suc {
		recordMsgSendTime(userId, entities.SMkPointsExpiring)
	}
}

func recordMsgSendTime(userId int64, key entities.SubscribeMessageKey) {
	var userConfigPointsExpiringKey = getSendTimeKey(key)
	_, err := interfaces.S.UserConfig.PutValueToConfig(userId, userConfigPointsExpiringKey, time.Now())
	if err != nil {
		log.Error("PutValueToConfig(userId, userConfigPointsExpiringKey, time.Now()) error", err, ",userId=", userId)
	}
}

func hasSendMsgToday(msgKey entities.SubscribeMessageKey, userId int64) bool {
	var t time.Time
	err := interfaces.S.UserConfig.GetConfigToValue(userId, getSendTimeKey(msgKey), &t)
	if err != nil {
		log.Error("GetConfigToValue(userId, getSendTimeKey(msgKey), &t) error:", err)
		return false
	}
	return t.After(utils.TodayBegin())
}

func sendWalkConvertExpiringMsg() {
	interfaces.S.User.TraverseSubscribeMessageTicketUser(entities.SMkWalkExpiring, func(user []*entities.UserDTO) {
		userIds := userDtoMapUserId(user)
		todayHasPointFlowUserIds, err := interfaces.S.Points.IsTodayHasPointFlowOfTypeBatched(userIds, entities.PointsEventTypeWalk)
		if err != nil {
			log.Error("IsTodayHasPointFlowOfTypeBatched error", err.Error())
			return
		}

		todayHasPointFlowUserIdSet := utils.NewInt64Set()
		todayHasPointFlowUserIdSet.AddAll(todayHasPointFlowUserIds)

		events, err := interfaces.S.User.BatchGetLastUserEventByType(userIds, entities.UserEventTypeGetWalkStep)
		if err != nil {
			log.Error("BatchGetLastUserEventByType error", err.Error())
			return
		}

		for _, u := range user {
			if _, ok := events[u.ID]; !ok {
				log.Debug("user", u.ID, "has no event report. skip")
				continue
			}

			_, err = interfaces.S.SubscribeMsg.Send(u.ID, u.Openid, entities.SMkWalkExpiring, map[string]interface{}{
				"thing1": map[string]interface{}{
					"value": "你还有步数未兑换",
				},
				"time2": map[string]interface{}{
					"value": "23:59:59",
				},
			})
			if err != nil {
				log.Error("Send(user.ID, user.Openid, service.SMkWalkExpiring ... error:", err, ",userId=", u.ID)
			}
		}
	})
}

func userDtoMapUserId(u []*entities.UserDTO) []int64 {
	var arr []int64

	for _, e := range u {
		arr = append(arr, e.ID)
	}

	return arr
}
