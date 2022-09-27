package async

// CompositeFutureAll 所有成功则成功，否则失败
func CompositeFutureAll[T any](futures []*Future[T]) (err error) {
	for _, f := range futures {
		_, err = f.Get()
		if err != nil {
			break
		}
	}

	return
}

// CompositeFutureAny 有一个成功则为成功，全部失败则为失败
func CompositeFutureAny[T any](futures []*Future[T]) (err error) {
	for _, f := range futures {
		_, err = f.Get()
		if err == nil {
			break
		}
	}

	return
}
