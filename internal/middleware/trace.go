package middleware

import (
	"aed-api-server/internal/pkg"
	"aed-api-server/internal/pkg/utils"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

func Trace(c *gin.Context) {
	utils.SetTraceId("", func() {
		requestID := c.GetHeader(pkg.TraceHeaderKey)
		log.Infof("bind requestId:%s", requestID)
		c.Next()
	})
}
