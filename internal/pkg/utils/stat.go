package utils

import (
	log "github.com/sirupsen/logrus"
	"time"
)

func TimeStat(processName string) func() {
	beginTime := time.Now()
	return func() {
		log.Debugf("stat <%s> process use time:%v\n", processName, time.Now().Sub(beginTime))
	}
}
