package server

import (
	"aed-api-server/internal/interfaces/facility"
	"aed-api-server/internal/pkg/utils"
	"fmt"
	"github.com/robfig/cron/v3"
	log "github.com/sirupsen/logrus"
)

func initScheduler() {
	cronTab := cron.New()
	comp := component.FindInstanceById("locker")
	locker := comp.(facility.Locker)

	for _, instances := range component.GetInstances() {
		if o, ok := instances.Obj.(facility.Scheduler); ok {
			o.Cron(func(spec string, jobName facility.JobName, fn func()) {
				lockKey := fmt.Sprintf("lock-job:%s", jobName)

				wrapFn := utils.GoWrapWithNewTraceId(func() {
					err := locker.LockFn(lockKey, fn)
					if err != nil {
						log.Error("cron get lock error")
					}
				})
				_, err := cronTab.AddFunc(spec, wrapFn)
				if err != nil {
					panic("cron.AddFunc for " + string(jobName) + " err:" + err.Error())
				}
				log.Infof("Add cron item: %s :%s", spec, jobName)
			})
		}
	}
	cronTab.Start()
}
