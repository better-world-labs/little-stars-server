package facility

import "time"

// Cache 使用方式，在service中依赖注入：
// ```go
// type svc struct {
//		facility.Cache `inject:"cache"`
// }
// ...
// func (s*svc) business(){
// 		s.Get(key, ttl)
// }
// ```
type Cache interface {
	Get(key string, v interface{}) (existed bool, err error)
	Put(key string, v interface{}, ttl time.Duration) (err error)
	Remove(key string) (err error)
}
