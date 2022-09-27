package async

import "aed-api-server/internal/pkg/utils"

//var pool *Pool

//func Start(config Config) error {
//	pool = NewPool(config)
//	pool.Start()
//	return nil
//}
//
//func Stop() {
//	pool.Stop()
//}

type Fn func() (interface{}, error)

func RunTask[T any](f Fn) *Future[T] {
	future := newFuture[T]()
	utils.Go(func() {
		defer close(future.resultChan)

		result, err := f()
		if err != nil {
			future.resultChan <- err
		} else {
			future.resultChan <- result
		}
	})

	return future
}
