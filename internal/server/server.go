package server

import (
	"aed-api-server/internal/interfaces"
	"aed-api-server/internal/module/achievement"
	"aed-api-server/internal/module/friends"
	"aed-api-server/internal/module/img"
	"aed-api-server/internal/module/speech"
	"aed-api-server/internal/module/user"
	"aed-api-server/internal/pkg/asserts"
	"aed-api-server/internal/pkg/cache"
	"aed-api-server/internal/pkg/config"
	"aed-api-server/internal/pkg/db"
	config2 "aed-api-server/internal/pkg/domain/config"
	"aed-api-server/internal/pkg/domain/emitter"
	"aed-api-server/internal/pkg/environment"
	"aed-api-server/internal/pkg/sms"
	"aed-api-server/internal/pkg/star"
	"aed-api-server/internal/pkg/tencent"
	"aed-api-server/internal/pkg/utils"
	"aed-api-server/internal/service/merit_tree"
	"aed-api-server/internal/service/subscribe_msg"
	"context"
	"fmt"
	"github.com/magiconair/properties"
	"gitlab.openviewtech.com/openview-pub/gopkg/inject"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"gitlab.openviewtech.com/openview-pub/gopkg/log"
)

var (
	conf       *config.AppConfig
	eng        *gin.Engine
	httpServer *http.Server
	component  *inject.Component
)

func SetConfig(c *config.AppConfig) {
	conf = c
}

func SetGin(engine *gin.Engine) {
	eng = engine
}

type loggerRequestGetter struct{}

func (*loggerRequestGetter) Getter() string {
	return utils.GetTraceId()
	//return utils.RequestId()
}

func initLog(conf *config.AppConfig) {
	conf.Log.RequestIdGetter = &loggerRequestGetter{}
	log.Init(conf.Log)
}

// Initialize 初始化阶段
func Initialize(loader func(conf *config.AppConfig, component *inject.Component), p *properties.Properties) *inject.Component {
	if conf == nil {
		panic("no config found")
	}

	component = &inject.Component{}
	initLog(conf)
	cache.InitPool(conf.Redis)
	initEmitter(conf.Domain)
	interfaces.InitConfig(conf) //TODO 清理
	prefix := environment.GetDomainPrefix(conf.Server.Env)
	speech.SetAidCaller(speech.NewAidCaller(prefix, speech.NewPathGenerator(speech.NewTokenService())))
	speech.SetUserFinder(speech.NewUserFinder(conf.Notifier.UserFinder))
	user.InitJwt(conf.JwtConfig.Secret, conf.JwtConfig.ExpiresIn)
	db.InitEngine(conf.Database)
	tencent.Init(&conf.MapConfig)
	friends.Init()
	merit_tree.InitWalk(conf)
	initAsserts()
	initImg()
	achievement.Init()
	sms.InitSmsClient(conf.SmsClient)
	star.Init(conf.MiniProgramQrcode)
	component.Conf(p)
	loader(conf, component)
	component.Install()
	return component
}

// Start 启动阶段
func Start() {
	initRouter(conf)
	subscribe_msg.InitScheduler()
	emitter.Start()
	interfaces.S.Cron.Start()
	startHttpServer()
}

// Stop 停止阶段
func Stop() {
	log.DefaultLogger().Info("shutting down server")
	stopHttpServer()
	emitter.Stop()
	//停止定时器
	interfaces.S.Cron.Stop()
}

func startHttpServer() {
	httpServer = &http.Server{
		Addr:    fmt.Sprintf(":%d", conf.Server.Port),
		Handler: eng,
	}

	go func() {
		if err := httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Info("http server error: %v", err)
			panic(err)
		}
	}()

	log.DefaultLogger().Infof("HTTP server started on port %d", conf.Server.Port)
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
