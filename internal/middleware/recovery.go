package middleware

import (
	"aed-api-server/internal/pkg"
	"aed-api-server/internal/pkg/response"
	"aed-api-server/internal/pkg/utils"
	"github.com/gin-gonic/gin"
	"gitlab.openviewtech.com/openview-pub/gopkg/log"
)

func Recovery(c *gin.Context) {
	requestID := c.MustGet(pkg.TraceHeaderKey)
	defer func() {
		if err := recover(); err != nil {
			log.DefaultLogger().Errorf("[%s] handle panic: %v, %s", requestID, err, utils.PanicTrace(2))
			response.ReplyError(c, err)
			c.Abort()
		}
	}()
	c.Next()
}
