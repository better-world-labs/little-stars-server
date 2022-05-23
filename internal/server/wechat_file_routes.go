package server

import (
	"aed-api-server/internal/pkg/utils"
	"github.com/gin-gonic/gin"
	"io/ioutil"
)

func InitWechatFileRoutes(g *gin.RouterGroup) {
	g.GET("/79pnqPgC5T.txt", func(c *gin.Context) {
		bytes, err := ioutil.ReadFile("assert/wechat/79pnqPgC5T.txt")
		utils.MustNil(err, err)
		_, err = c.Writer.Write(bytes)
	})

	g.GET("/cert/79pnqPgC5T.txt", func(c *gin.Context) {
		bytes, err := ioutil.ReadFile("assert/wechat/79pnqPgC5T.txt")
		utils.MustNil(err, err)
		_, err = c.Writer.Write(bytes)
	})

}
