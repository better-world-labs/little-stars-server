package achievement

import (
	"aed-api-server/internal/interfaces"
	"aed-api-server/internal/interfaces/entities"
	service2 "aed-api-server/internal/interfaces/service"
	"aed-api-server/internal/module/aid/track"
	"aed-api-server/internal/pkg/cache"
	"fmt"
	log "github.com/sirupsen/logrus"
	"time"
)

var t *time.Timer

const (
	H = 0
	M = 0
	S = 0
)

func InitTask() {
	duration := getDuration(H, M, S)
	t = time.NewTimer(duration)
	go run()
}

func aidService() service2.AidService {
	return interfaces.S.Aid
}

func getDuration(h int, m int, s int) time.Duration {
	now := time.Now()
	next := time.Date(now.Year(), now.Month(), now.Day(), h, m, s, 0, now.Location())
	if next.Sub(now) < 0 {
		next = next.Add(24 * time.Hour)
	}

	fmt.Printf("task run on %s\n", next)
	return next.Sub(time.Now())
}

func run() {
	for {
		<-t.C
		doTask()
		t.Reset(getDuration(H, M, S))
	}
}

func doTask() {
	log.Info("start medal issue")
	lock, err := cache.GetDistributeLock("aed-lock-medal-issue", 5000)
	if err != nil {
		log.Error("medal issue error: %v", err)
		return
	}

	defer lock.Release()

	if !lock.Locked() {
		log.Warnf("get distribute lock failed, give up")
		return
	}

	doIssue()
}

func doIssue() {
	arr, err := aidService().ListHelpInfosInner24h()
	if err != nil {
		log.Errorf("do issue error: %v", err)
		return
	}

	for _, helpInfo := range arr {
		doIssueForUser(helpInfo.Publisher, helpInfo)
		t, err := track.GetService().GetAidDeviceGotTracksSorted(helpInfo.ID)
		if err != nil {
			log.Errorf("do issue error: %v", err)
			continue
		}

		for _, e := range t {
			doIssueForUser(e.UserID, helpInfo)
		}
	}
}

func doIssueForUser(userID int64, helpInfo *entities.HelpInfo) {
	log.Infof("do issue for user %d", helpInfo.Publisher)

	err := interfaces.S.Medal.AwardMedalSaveLife(userID, helpInfo.ID)
	if err != nil {
		log.Errorf("do issue for user %s failed:%v", userID, err)
	}
}
