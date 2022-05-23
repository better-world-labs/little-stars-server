package achievement

import (
	"aed-api-server/internal/interfaces"
	"aed-api-server/internal/interfaces/entities"
	"aed-api-server/internal/module/aid/track"
	"aed-api-server/internal/pkg/cache"
	"fmt"
	"gitlab.openviewtech.com/openview-pub/gopkg/log"
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
	log.DefaultLogger().Info("start medal issue")
	b, lock, err := cache.GetDistributeLock("aed-lock-medal-issue", 5000)
	if err != nil {
		log.DefaultLogger().Error("medal issue error: %v", err)
		return
	}

	if !b {
		log.DefaultLogger().Warnf("get distribute lock failed, give up")
		return
	}

	defer lock.Release()

	doIssue()
}

func doIssue() {
	arr, err := aidService.ListHelpInfosInner24h()
	if err != nil {
		log.DefaultLogger().Errorf("do issue error: %v", err)
		return
	}

	for _, helpInfo := range arr {
		doIssueForUser(helpInfo.Publisher, helpInfo)
		t, err := track.GetService().GetAidDeviceGotTracksSorted(helpInfo.ID)
		if err != nil {
			log.DefaultLogger().Errorf("do issue error: %v", err)
			continue
		}

		for _, e := range t {
			doIssueForUser(e.UserID, helpInfo)
		}
	}
}

func doIssueForUser(userID int64, helpInfo *entities.HelpInfo) {
	log.DefaultLogger().Infof("do issue for user %d", helpInfo.Publisher)

	err := interfaces.S.Medal.AwardMedalSaveLife(userID, helpInfo.ID)
	if err != nil {
		log.DefaultLogger().Errorf("do issue for user %s failed:%v", userID, err)
	}
}
