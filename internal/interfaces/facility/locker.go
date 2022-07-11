package facility

import "time"

type Unlock func()

// Locker 使用方式，在service中依赖注入：
// ```go
// type svc struct {
//		facility.Locker `inject:"locker"`
// }
// ...
// func (s*svc) business(){
// 		s.Lock(key, ttl)
// }
// ```
type Locker interface {
	Lock(key string, ttl time.Duration) (unlock Unlock, err error)
	LockFn(key string, fn func()) (err error)
}
