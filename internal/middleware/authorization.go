package middleware

import (
	"aed-api-server/internal/interfaces"
	"aed-api-server/internal/pkg"
	"aed-api-server/internal/pkg/response"
	"github.com/gin-gonic/gin"
)

func Authorize(c *gin.Context) {
	authorization := c.GetHeader(pkg.AuthorizationHeaderKey)
	user, err := interfaces.S.User.ParseInfoFromJwtToken(authorization)
	if err != nil {
		response.HTTPComplete(c, 401, response.NewResponse(-1, err.Error(), nil))
		c.Abort()
		return
	}

	c.Set(pkg.AccountIDKey, user.ID)
	c.Set(pkg.AccountKey, user)
	c.Next()
}
