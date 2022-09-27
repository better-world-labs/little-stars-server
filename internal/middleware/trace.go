package middleware

import (
	"aed-api-server/internal/pkg"
	"aed-api-server/internal/pkg/utils"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

func Trace(c *gin.Context) {
	//健康检查，不记录日志
	if c.Request.URL.Path == "/api/health-check" {
		c.AbortWithStatus(200)
		return
	}

	utils.SetTraceId("", func() {
		defer utils.TimeStat("api-stat:" + c.Request.Method + " " + c.Request.URL.Path + "$")()
		requestID := c.GetHeader(pkg.TraceHeaderKey)
		log.Infof("bind requestId:%s", requestID)
		c.Next()
	})
}
