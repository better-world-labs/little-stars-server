package async

type Future struct {
	resultChan chan interface{}
	result     interface{}
	err        error
	completed  bool
}

func newFuture() *Future {
	future := Future{
		resultChan: make(chan interface{}, 1),
	}

	return &future
}

func (f *Future) Get(binder interface{}) error {
	if !f.completed {
		res := <-f.resultChan
		switch res.(type) {
		case error:
			f.err = res.(error)

		default:
			f.result = res
		}
		f.completed = true
	}

	if f.err != nil {
		return f.err
	}

	binder = f.result
	return nil
}
