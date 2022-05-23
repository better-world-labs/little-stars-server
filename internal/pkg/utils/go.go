package utils

import (
	"github.com/jtolds/gls"
	"gitlab.openviewtech.com/openview-pub/gopkg/log"
)

func Go(fn func()) {
	gls.Go(func() {
		defer func() {
			if err := recover(); err != nil {
				log.DefaultLogger().Errorf("handle panic: %v, %s", err, PanicTrace(2))
			}
		}()
		fn()
	})
}
