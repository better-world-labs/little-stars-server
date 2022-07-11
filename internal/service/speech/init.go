package speech

import (
	"aed-api-server/internal/interfaces"
	"aed-api-server/internal/interfaces/events"
	"aed-api-server/internal/interfaces/facility"
	"aed-api-server/internal/pkg/environment"
)

var (
	caller AidCenterCaller
	finder UserFinder
)

func Init() {
	conf := interfaces.GetConfig()
	prefix := environment.GetDomainPrefix(conf.Server.Env)
	SetAidCaller(NewAidCaller(prefix, NewPathGenerator(NewTokenService())))
	SetUserFinder(NewUserFinder(conf.Notifier.UserFinder))
}

func (t tokenService) Listen(on facility.OnEvent) {
	on(&events.HelpInfoPublishedEvent{}, onAidPublished)
}

func SetAidCaller(centerCaller AidCenterCaller) {
	caller = centerCaller
}

func SetUserFinder(userFinder UserFinder) {
	finder = userFinder
}
