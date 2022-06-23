package command

import (
	"aed-api-server/internal/pkg/config"
	"aed-api-server/internal/server"
	"github.com/gin-gonic/gin"
	"github.com/magiconair/properties"
	"github.com/urfave/cli/v2"
	"gitlab.openviewtech.com/openview-pub/gopkg/inject"
	"os"
	"os/signal"
	"syscall"
)

func Run(context *cli.Context) error {
	baseDir := context.String("b")
	env := context.String("e")
	println("env:", env)
	c, p, err := config.LoadConfigX(baseDir, env)
	if err != nil {
		return err
	}

	StartServer(c, p)

	quit := make(chan os.Signal)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	<-quit
	server.Stop()
	return nil
}

func StartServer(c *config.AppConfig, p *properties.Properties) {
	gin.SetMode(gin.ReleaseMode)
	eng := gin.New()
	server.SetGin(eng)
	server.SetConfig(c)
	server.Initialize(LoadComponents, p)
	go server.Start()
}

func LoadComponents(c *config.AppConfig, component *inject.Component) {
	LoadServices(c, component)
	LoadController(c, component)
	autoLoad(component)
}
