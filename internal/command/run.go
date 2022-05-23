package command

import (
	"aed-api-server/internal/pkg/config"
	"aed-api-server/internal/server"
	"github.com/gin-gonic/gin"
	"github.com/urfave/cli/v2"
	"os"
	"os/signal"
	"syscall"
)

func Run(context *cli.Context) error {
	f := context.String("c")
	println("config file:", f)
	c, err := config.LoadConfig(f)
	if err != nil {
		return err
	}

	StartServer(c)

	quit := make(chan os.Signal)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	<-quit
	server.Stop()
	return nil
}

func StartServer(c *config.AppConfig) {
	gin.SetMode(gin.ReleaseMode)
	eng := gin.New()
	server.SetGin(eng)
	server.SetConfig(c)
	server.Initialize()
	go server.Start()
}
