package utils

import (
	"errors"
	"github.com/google/uuid"
	log "github.com/sirupsen/logrus"
)

type PromiseProcessor = func() (interface{}, error)

func PromiseAll(fnList ...PromiseProcessor) ([]interface{}, error) {
	type Pair struct {
		Rst   interface{}
		Index int
	}

	var n = len(fnList)
	pairChan := make(chan Pair, n)
	errChan := make(chan error, 1)

	for i := range fnList {
		pair := Pair{
			Index: i,
		}
		go func(fn func() (interface{}, error)) {
			defer func() {
				if err := recover(); err != nil {
					panicId := uuid.New().String()
					log.Errorf("panic(%s): %v, %s", panicId, err, PanicTrace(1))
					errChan <- errors.New("has panic(" + panicId + ")")
				}
			}()

			rst, err := fn()
			if err != nil {
				errChan <- err
				return
			}
			pair.Rst = rst
			pairChan <- pair
		}(fnList[i])
	}
	rstList := make([]interface{}, n, n)
	i := 0
	for {
		var err error
		select {
		case pair := <-pairChan:
			rstList[pair.Index] = pair.Rst
			i++
		case err = <-errChan:

		}
		if err != nil {
			println("error end")
			return nil, err
		}
		if i == n {
			println("end")
			return rstList, nil
		}
	}
}

func PromiseAllArr(fnList []PromiseProcessor) ([]interface{}, error) {
	return PromiseAll(fnList...)
}
