package middleware

import (
	"aed-api-server/internal/pkg/response"
	"aed-api-server/internal/pkg/utils"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

func Recovery(c *gin.Context) {
	traceID := utils.GetTraceId()
	defer func() {
		if err := recover(); err != nil {
			log.Errorf("[%s] handle panic: %v, %s", traceID, err, utils.PanicTrace(2))
			response.ReplyError(c, err)
			c.Abort()
		}
	}()
	c.Next()
}
