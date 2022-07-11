package utils

import (
	"aed-api-server/internal/pkg"
	"github.com/gin-gonic/gin"
	"strconv"
)

func GetContextUserId(ctx *gin.Context) int64 {
	return ctx.GetInt64(pkg.AccountIDKey)
}

func GetContextPathParamInt64(ctx *gin.Context, key string) (int64, error) {
	param := ctx.Param(key)
	return strconv.ParseInt(param, 10, 64)
}
