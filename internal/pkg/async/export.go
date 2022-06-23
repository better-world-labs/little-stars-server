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

func RunTask(f func() (interface{}, error)) *Future {
	future := newFuture()
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
