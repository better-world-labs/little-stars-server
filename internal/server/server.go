package server

import (
	"aed-api-server/internal/interfaces"
	"aed-api-server/internal/pkg/asserts"
	"aed-api-server/internal/pkg/cache"
	"aed-api-server/internal/pkg/config"
	"aed-api-server/internal/pkg/db"
	config2 "aed-api-server/internal/pkg/domain/config"
	"aed-api-server/internal/pkg/domain/emitter"
	"aed-api-server/internal/pkg/sms"
	"aed-api-server/internal/pkg/tencent"
	"aed-api-server/internal/pkg/utils"
	"aed-api-server/internal/service/img"
	"aed-api-server/internal/service/medal"
	"aed-api-server/internal/service/speech"
	"aed-api-server/internal/service/user"
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/magiconair/properties"
	"gitlab.openviewtech.com/openview-pub/gopkg/inject"
	"gitlab.openviewtech.com/openview-pub/gopkg/log"
	"net/http"
	"time"
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

	interfaces.InitConfig(conf) //TODO 清理
	speech.Init()
	user.InitJwt(conf.JwtConfig.Secret, conf.JwtConfig.ExpiresIn)
	db.InitEngine(conf.Database)
	tencent.Init(&conf.MapConfig)
	initAsserts()
	initImg()
	medal.Init()
	sms.InitSmsClient(conf.SmsClient)

	component.Conf(p)
	loader(conf, component)
	component.Install()

	initEmitter(conf.Domain)
	//initScheduler 依赖了component，须放在component.Install()之后
	initScheduler()
	return component
}

// Start 启动阶段
func Start() {
	initRouter(conf)
	emitter.Start()
	startHttpServer()
}

// Stop 停止阶段
func Stop() {
	log.DefaultLogger().Info("shutting down server")
	stopHttpServer()
	emitter.Stop()
}

func startHttpServer() {
	httpServer = &http.Server{
		Addr:    fmt.Sprintf(":%d", conf.Server.Port),
		Handler: eng,
	}

	go func() {
		if err := httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Infof("http server error: %v", err)
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
