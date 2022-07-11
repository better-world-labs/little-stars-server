package user

import (
	"aed-api-server/internal/interfaces/facility"
	"gitlab.openviewtech.com/openview-pub/gopkg/log"
)

func (s *Service) Cron(run facility.RunFuncOnceAt) {

	//0点运行，锁定30分钟
	run("0 0 * * *", facility.UpdateAllPositionHeat, func() {
		err := s.UpdateAllPositionHeat()
		if err != nil {
			log.Error("UpdateAllPositionHeat err:", err)
		}
	})
}
