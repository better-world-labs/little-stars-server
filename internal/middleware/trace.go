package middleware

import (
	"aed-api-server/internal/pkg"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func Trace(c *gin.Context) {
	requestID := c.GetHeader(pkg.TraceHeaderKey)
	if requestID == "" {
		requestID = uuid.New().String()
	}

	remoteIP := c.GetHeader("X-Forwarded-For")
	if remoteIP == "" {
		ip, _ := c.RemoteIP()
		remoteIP = ip.String()
	}
	c.Set(pkg.TraceHeaderKey, requestID)
	c.Writer.Header().Set(pkg.TraceHeaderKey, requestID)
	c.Next()
}
