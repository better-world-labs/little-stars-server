package middleware

import (
	"aed-api-server/internal/interfaces/entities"
	"aed-api-server/internal/module/user"
	"aed-api-server/internal/pkg"
	"aed-api-server/internal/pkg/db"
	"aed-api-server/internal/pkg/response"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"strings"
)

func Authorize(c *gin.Context) {
	authorization := c.GetHeader(pkg.AuthorizationHeaderKey)
	split := strings.Split(authorization, " ")

	if len(split) != 2 || split[0] != "Bearer" || split[1] == "" {
		response.HTTPComplete(c, 401, response.NewResponse(-1, "invalid token", nil))
		c.Abort()
		return
	}

	token := split[1]
	claims, err := user.ParseToken(token)
	if err != nil {
		log.Errorf("handle authorization: %v", err)
		response.HTTPComplete(c, 401, response.NewResponse(-1, "invalid token", nil))
		c.Abort()
		return
	}

	session := db.GetSession()
	defer session.Close()

	var acc entities.User
	exists, err := session.Table("account").Where("id=?", claims.ID).Get(&acc)
	if err != nil {
		log.Errorf("handle authorization: %v", err)
		response.HTTPComplete(c, 401, response.NewResponse(-1, "invalid token", nil))
		c.Abort()
		return
	}

	if !exists {
		response.HTTPComplete(c, 401, response.NewResponse(-1, "invalid token", nil))
		return
	}

	c.Set(pkg.AccountIDKey, claims.ID)
	c.Set(pkg.AccountKey, &acc)
	c.Next()
}
