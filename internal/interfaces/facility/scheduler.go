package facility

import (
	"aed-api-server/internal/interfaces/entities"
	"fmt"
)

type JobName string

const (
	UpdateAllPositionHeat JobName = "UpdateAllPositionHeat" //更新所有用户的位置热点
	FriendAddPoints       JobName = "FriendAddPoints"       //好友加成
	GameProcess           JobName = "GameProcess"           //游戏进程
)

//GetSendMsgJobLockKey 使用前缀 `jobLock:sendMsg:`
func GetSendMsgJobLockKey(key entities.SubscribeMessageKey) JobName {
	return JobName(fmt.Sprintf("sendMsg:%s", key))
}

//RunFuncOnceAt 定时跑
//@Param fn 要调用的函数
//@Param spec 调用时间 cron tab 格式
//@Param lockKey 分布式锁的key
//@Param lockTtl 锁定时长
type RunFuncOnceAt func(spec string, jobName JobName, fn func())

type Scheduler interface {

	//Cron use: Cron(run facility.RunFuncOnceAt)
	Cron(run RunFuncOnceAt)
}
