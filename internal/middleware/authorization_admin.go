package middleware

import (
	"aed-api-server/internal/module/user"
	"aed-api-server/internal/pkg"
	"aed-api-server/internal/pkg/config"
	"aed-api-server/internal/pkg/db"
	"aed-api-server/internal/pkg/response"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"strings"
)

type AuthorizationAdmin struct {
	backend config.Backend
}

func NewAuthorizationAdmin(c config.Backend) *AuthorizationAdmin {
	return &AuthorizationAdmin{backend: c}
}

func (a *AuthorizationAdmin) AuthorizeAdmin(c *gin.Context) {
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

	if claims.ID != a.backend.Id {
		log.Errorf("handle authorization: %v", err)
		response.HTTPComplete(c, 401, response.NewResponse(-1, "invalid token", nil))
		c.Abort()
		return
	}

	c.Set(pkg.AccountIDKey, claims.ID)
	c.Next()
}
