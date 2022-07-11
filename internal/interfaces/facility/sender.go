package facility

import "aed-api-server/internal/pkg/domain/emitter"

//Sender 使用方式，在service中依赖注入：
// ```go
// type svc struct {
//		facility.Sender `inject:"sender"`
// }
// ...
// func (s*svc) business(){
// 		s.Send(&event)
// }
// ```
type Sender interface {
	Send(events ...emitter.DomainEvent) error
}
