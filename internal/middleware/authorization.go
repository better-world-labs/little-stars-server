package middleware

import (
	"aed-api-server/internal/module/user"
	"aed-api-server/internal/pkg"
	"aed-api-server/internal/pkg/db"
	"aed-api-server/internal/pkg/global"
	"aed-api-server/internal/pkg/response"
	"aed-api-server/internal/pkg/utils"
	"github.com/gin-gonic/gin"
	"gitlab.openviewtech.com/openview-pub/gopkg/log"
	"strings"
)

func Authorize(c *gin.Context) {
	authorization := c.GetHeader(pkg.AuthorizationHeaderKey)
	split := strings.Split(authorization, " ")

	if len(split) != 2 || split[0] != "Bearer" || split[1] == "" {
		response.ReplyError(c, global.ErrorInvalidAccessToken)
		c.Abort()
		return
	}

	token := split[1]
	claims, err := user.ParseToken(token)
	if err != nil {
		log.DefaultLogger().Errorf("handle authorization: %v", err)
		response.ReplyError(c, global.ErrorInvalidAccessToken)
		c.Abort()
		return
	}

	session := db.GetSession()
	defer session.Close()

	var acc user.User
	exists, err := session.Table("account").Where("id=?", claims.ID).Get(&acc)
	if err != nil {
		log.DefaultLogger().Errorf("handle authorization: %v", err)
		response.ReplyError(c, global.ErrorInvalidAccessToken)
		c.Abort()
		return
	}

	utils.MustTrue(exists, global.ErrorAccountNotFound)
	c.Set(pkg.AccountIDKey, claims.ID)
	c.Set(pkg.AccountKey, &acc)
	c.Next()
}
