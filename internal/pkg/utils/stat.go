package utils

import (
	"gitlab.openviewtech.com/openview-pub/gopkg/log"
	"time"
)

func TimeStat(processName string) func() {
	beginTime := time.Now()
	return func() {
		log.DefaultLogger().Debugf("stat <%s> process use time:%v\n", processName, time.Now().Sub(beginTime))
	}
}
