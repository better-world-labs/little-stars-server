package controller

import (
	"aed-api-server/internal/interfaces"
	"aed-api-server/internal/pkg/response"
	"github.com/gin-gonic/gin"
)

type StatController struct{}

func (StatController) KpiStat(c *gin.Context) {
	stat, err := interfaces.S.Stat.DoKipStat()
	if err != nil {
		response.ReplyError(c, err)
		return
	}
	response.ReplyOK(c, stat)
}
