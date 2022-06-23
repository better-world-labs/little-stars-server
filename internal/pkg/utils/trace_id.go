package utils

import (
	"github.com/google/uuid"
	"github.com/jtolds/gls"
)

var (
	mgr        = gls.NewContextManager()
	traceIdKey = gls.GenSym()
)

func SetTraceId(traceId string, fn func()) {
	if traceId == "" {
		traceId = uuid.New().String()
	}
	mgr.SetValues(gls.Values{traceIdKey: traceId}, fn)
}

func GetTraceId() string {
	if traceId, ok := mgr.GetValue(traceIdKey); ok {
		return traceId.(string)
	} else {
		return ""
	}
}
