package service

import (
	"aed-api-server/internal/interfaces/entities"
	"aed-api-server/internal/interfaces/service"
	"aed-api-server/internal/module/user"
	"aed-api-server/internal/pkg/db"
	"aed-api-server/internal/pkg/domain/emitter"
	"github.com/go-xorm/xorm"
)

type traceService struct {
	Wc user.WechatClient `inject:"-"`
}

func NewTraceService() service.TraceService {
	return &traceService{}
}

func (s traceService) Create(code string, trace entities.Trace) (*entities.Trace, error) {
	session, err := s.Wc.CodeToSession(code)
	if err != nil {
		return nil, err
	}

	trace.OpenID = session.Openid
	if trace.OpenID == "" {
		return nil, err
	}

	return &trace, db.Begin(func(session *xorm.Session) error {
		_, err = db.GetEngine().Table("generalize_trace").Insert(trace)
		if err != nil {
			return err
		}
		return emitter.Emit(&trace)
	})
}

func (s traceService) GetEarliestSharerTrace(openid string) (*entities.Trace, bool, error) {
	var t entities.Trace
	exists, err := db.SQL(`select * from generalize_trace where open_id = ?
		and sharer regexp '^[0-9]+$' limit 1`, openid).Get(&t)
	if err != nil {
		return nil, exists, err
	}

	return &t, exists, nil
}
