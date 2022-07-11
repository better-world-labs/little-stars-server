package command

import (
	"aed-api-server/internal/interfaces"
	"aed-api-server/internal/pkg/config"
	"gitlab.openviewtech.com/openview-pub/gopkg/inject"
)

func loadServices(c *config.AppConfig, component *inject.Component) {
	component.Load(interfaces.S)
}
