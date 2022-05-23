package async

// CompositeFutureAll 所有成功则成功，否则失败
func CompositeFutureAll(futures []*Future) (err error) {
	for _, f := range futures {
		err = f.Get(nil)
		if err != nil {
			break
		}
	}

	return
}

// CompositeFutureAny 有一个成功则为成功，全部失败则为失败
func CompositeFutureAny(futures []*Future) (err error) {
	for _, f := range futures {
		err = f.Get(nil)
		if err == nil {
			break
		}
	}

	return
}
