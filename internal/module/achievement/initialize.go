package achievement

import (
	"aed-api-server/internal/pkg/cache"
)

var (
	caches cache.Cache
)

func Init() {
	caches = cache.NewLocalLRUCache(1024000)
	InitMedalStore()
	InitTask()
}
