package utils

import log "github.com/sirupsen/logrus"

func MustNil(i interface{}, wrap error) {
	if i != nil {
		log.Errorf("%v", i)
		panic(wrap)
	}
}

func MustTrue(b bool, err error) {
	if !b {
		panic(err)
	}
}
