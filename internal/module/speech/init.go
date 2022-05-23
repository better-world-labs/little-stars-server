package speech

import (
	"aed-api-server/internal/interfaces/events"
	"aed-api-server/internal/pkg/domain/emitter"
)

var (
	caller AidCenterCaller
	finder UserFinder
)

func SetAidCaller(centerCaller AidCenterCaller) {
	caller = centerCaller
}

func SetUserFinder(userFinder UserFinder) {
	finder = userFinder
}

func InitEventHandler() {
	emitter.On(&events.HelpInfoPublishedEvent{}, onAidPublished)
}
