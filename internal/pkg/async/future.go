package async

type Future[T any] struct {
	resultChan chan interface{}
	result     T
	err        error
	completed  bool
}

func newFuture[T any]() *Future[T] {
	future := Future[T]{
		resultChan: make(chan interface{}, 1),
	}

	return &future
}

func (f *Future[T]) Get() (T, error) {
	var t T
	if !f.completed {
		res := <-f.resultChan
		switch res.(type) {
		case error:
			f.err = res.(error)

		default:
			if assert, ok := res.(T); ok {
				f.result = assert
			} else {
				//TODO 想办法处理下，不能报错
				f.result = t
			}
		}

		f.completed = true
	}

	return t, f.err
}
