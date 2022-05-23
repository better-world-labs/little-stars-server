package achievement

import (
	"aed-api-server/internal/module/aid"
	"aed-api-server/internal/module/user"
	"aed-api-server/internal/pkg/cache"
)

var (
	caches     cache.Cache
	aidService aid.Service
)

func Init() {
	caches = cache.NewLocalLRUCache(1024000)
	aidService = aid.NewService(user.NewService(nil))
	InitMedalStore()
	InitTask()
}
