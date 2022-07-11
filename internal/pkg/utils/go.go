package utils

import (
	"github.com/jtolds/gls"
	log "github.com/sirupsen/logrus"
)

func Go(fn func()) {
	gls.Go(func() {
		defer func() {
			if err := recover(); err != nil {
				log.Errorf("handle panic: %v, %s", err, PanicTrace(2))
			}
		}()
		fn()
	})
}

func GoWrapWithNewTraceId(fn func()) func() {
	return func() {
		go func() {
			defer func() {
				if err := recover(); err != nil {
					log.Errorf("handle panic: %v, %s", err, PanicTrace(2))
				}
			}()
			SetTraceId("", fn)
		}()
	}
}
