package server

import (
	"aed-api-server/internal/interfaces"
	"aed-api-server/internal/middleware"
	"aed-api-server/internal/module/achievement"
	"aed-api-server/internal/module/aid"
	"aed-api-server/internal/module/evidence"
	"aed-api-server/internal/module/exam"
	"aed-api-server/internal/module/friends"
	"aed-api-server/internal/module/img"
	"aed-api-server/internal/module/speech"
	"aed-api-server/internal/module/trace"
	"aed-api-server/internal/module/user"
	"aed-api-server/internal/pkg/asserts"
	"aed-api-server/internal/pkg/cache"
	"aed-api-server/internal/pkg/config"
	"aed-api-server/internal/pkg/db"
	config2 "aed-api-server/internal/pkg/domain/config"
	"aed-api-server/internal/pkg/domain/emitter"
	"aed-api-server/internal/pkg/environment"
	"aed-api-server/internal/pkg/star"
	"aed-api-server/internal/pkg/tencent"
	"aed-api-server/internal/pkg/utils"
	"aed-api-server/internal/service/merit_tree"
	"context"
	"fmt"
	"io/ioutil"
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"gitlab.openviewtech.com/openview-pub/gopkg/log"
)

var (
	conf       *config.AppConfig
	eng        *gin.Engine
	httpServer *http.Server
)

func SetConfig(c *config.AppConfig) {
	conf = c
}

func SetGin(engine *gin.Engine) {
	eng = engine
}

// Initialize 初始化阶段
func Initialize() {
	if conf == nil {
		panic("no config found")
	}

	log.Init(conf.Log)
	cache.InitPool(conf.Redis)
	initEmitter(conf.Domain)
	interfaces.InitConfig(conf)
	prefix := environment.GetDomainPrefix(conf.Server.Environment)
	speech.SetAidCaller(speech.NewAidCaller(prefix, speech.NewPathGenerator(speech.NewTokenService())))
	speech.SetUserFinder(speech.NewUserFinder(conf.Notifier.UserFinder))
	user.InitJwt(conf.JwtConfig.Secret, conf.JwtConfig.ExpiresIn)
	db.InitEngine(conf.Database)
	tencent.Init(&conf.MapConfig)
	exam.Init(conf)
	trace.Init(conf)
	friends.Init()
	merit_tree.InitWalk(conf)
	merit_tree.InitEarly()
	initAsserts()
	evidence.Init(conf)
	initImg()
	initGin()
	star.Init(conf.MiniProgramQrcode)
	achievement.Init()
}

func initGin() {
	// middleware
	eng.Use(middleware.Trace)
	eng.Use(middleware.AccessLog)
	eng.Use(middleware.Recovery)
	eng.Use(middleware.Cors())
	routerGroup := eng.Group("/api")
	adminGroup := eng.Group("/admin-api")
	pageGroup := eng.Group("/p")
	wechatFileGroup := eng.Group("/share")
	eng.GET("/79pnqPgC5T.txt", func(c *gin.Context) {
		bytes, err := ioutil.ReadFile("assert/wechat/79pnqPgC5T.txt")
		utils.MustNil(err, err)
		_, err = c.Writer.Write(bytes)
	})

	InitWechatFileRoutes(wechatFileGroup)
	InitAuthorizedRoutes(routerGroup, conf, user.NewWechatClient(&conf.Wechat))
	InitRoutes(routerGroup, conf, user.NewWechatClient(&conf.Wechat))
	InitAdminRoutes(adminGroup, conf, user.NewWechatClient(&conf.Wechat))
	InitPageRoutes(pageGroup, conf, user.NewWechatClient(&conf.Wechat))
}

func InitPageRoutes(g *gin.RouterGroup, c *config.AppConfig, client user.WechatClient) {
	controller := aid.NewController(aid.NewService(user.NewService(client)))
	g.GET("/c/:token/:aid", controller.ActionAidCalledPage)
}

// Start 启动阶段
func Start() {
	emitter.Start()
	startHttpServer()
}

// Stop 停止阶段
func Stop() {
	log.DefaultLogger().Info("shutting down server")

	g := sync.WaitGroup{}
	g.Add(2)

	go func() { stopHttpServer(); g.Done() }()
	go func() { emitter.Stop(); g.Done() }()

	g.Wait()
}

func startHttpServer() {
	httpServer = &http.Server{
		Addr:    fmt.Sprintf(":%d", conf.Server.Port),
		Handler: eng,
	}
	group := sync.WaitGroup{}
	group.Add(1)
	go func() {
		if err := httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Info("http server error: %v", err)
			panic(err)
		}
		group.Done()
	}()

	log.DefaultLogger().Infof("HTTP server started on port %d", conf.Server.Port)
	group.Wait()
}

func stopHttpServer() {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := httpServer.Shutdown(ctx); err != nil {
		log.DefaultLogger().Infof("server forced to shutdown: %v", err)
	}
}

func initEmitter(c config2.DomainEventConfig) {
	emitter.SetContext(context.Background())
	emitter.SetConfig(&c)
	initEventHandler()
}

func initAsserts() {
	err := asserts.LoadResourceDir("assert")
	if err != nil {
		panic(err)
	}
}

func initImg() {
	err := img.Init()
	if err != nil {
		panic(err)
	}
}
