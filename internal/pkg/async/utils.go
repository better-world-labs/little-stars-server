package async

import (
	"aed-api-server/internal/pkg/utils"
	"errors"
)

func RunAsync[T any](f func() (T, error)) func() (T, error) {
	resChan := make(chan interface{}, 1)

	utils.Go(func() {
		defer close(resChan)
		res, err := f()
		if err == nil {
			resChan <- res
		} else {
			resChan <- err
		}
	})

	return func() (T, error) {
		var o T
		res := <-resChan

		switch res.(type) {
		case T:
			return res.(T), nil
		case error:
			return o, res.(error)
		default:
			return o, errors.New("invalid type from chan")
		}
	}
}
