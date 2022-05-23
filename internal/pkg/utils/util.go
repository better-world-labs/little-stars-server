package utils

import "gitlab.openviewtech.com/openview-pub/gopkg/log"

func MustNil(i interface{}, wrap error) {
	if i != nil {
		log.DefaultLogger().Errorf("%v", i)
		panic(wrap)
	}
}

func MustTrue(b bool, err error) {
	if !b {
		panic(err)
	}
}
