package controller

import (
	"aed-api-server/internal/interfaces"
	"aed-api-server/internal/interfaces/service"
	"aed-api-server/internal/pkg/cache"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gitlab.openviewtech.com/openview-pub/gopkg/route"
	"time"
)

const (
	ShortLinkKey = "SHORT_LINK"
)

type LinkParam struct {
	Path  string `json:"path" binding:"required"`
	Query string `json:"query"`
}

type ShortLinkController struct {
	Wx service.IWechat `inject:"-"`

	Host string `conf:"server.host"`
}

func (c ShortLinkController) GenerateMPLink(ctx *gin.Context) (interface{}, error) {
	var param LinkParam
	err := ctx.ShouldBindJSON(&param)
	if err != nil {
		return nil, err
	}

	link, err := CreateLink(param)
	if err != nil {
		return nil, err
	}

	return fmt.Sprintf("https://%s/api/l/%s", c.Host, link), nil
}

func CreateLink(param interface{}) (string, error) {
	conn := cache.GetConn()
	defer conn.Close()

	b, err := json.Marshal(param)
	if err != nil {
		return "", err
	}
	code := genericLinkCode()
	r, err := conn.Do("SET", key(code), string(b), "NX")
	if err != nil {
		return "", err
	}

	if r == nil {
		return CreateLink(param)
	}

	return code, nil
}

func GetLinkParam(code string) (LinkParam, bool, error) {
	var param LinkParam

	conn := cache.GetConn()
	defer conn.Close()

	r, err := conn.Do("GET", key(code))
	if err != nil {
		return param, false, err
	}

	if r == nil {
		return param, false, nil
	}

	err = json.Unmarshal(r.([]byte), &param)
	if err != nil {
		return param, false, err
	}

	return param, true, nil
}

func (c ShortLinkController) ProxyLink(ctx *gin.Context) (interface{}, error) {
	code := ctx.Param("code")

	param, exists, err := GetLinkParam(code)
	if err != nil {
		return nil, err
	}

	if !exists {
		return "链接已失效", nil
	}

	link, err := c.Wx.GenericUrlLink(param.Path, param.Query)
	if err != nil {
		return nil, err
	}

	ctx.Redirect(302, link)
	return nil, nil
}

//go:inject-component
func NewShortLinkController() *ShortLinkController {
	return &ShortLinkController{}
}

func (c ShortLinkController) MountNoAuthRouter(r *route.Router) {
	voteGroup := r.Group("/l")
	voteGroup.POST("/generate", c.GenerateMPLink)
	voteGroup.GET("/:code", c.ProxyLink)
}

func genericLinkCode() string {
	u := uuid.NewString()
	str := []byte(u)[:6]
	return base64.StdEncoding.EncodeToString([]byte(fmt.Sprintf("%s%d", str, time.Now().Unix())))
}

func key(code string) string {
	return fmt.Sprintf("%s-%s-%s", interfaces.GetConfig().Server.Env, ShortLinkKey, code)
}
