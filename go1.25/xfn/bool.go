package xfn

func True() bool {
	return true
}

func False() bool {
	return false
}

func One2True[T any](T) bool {
	return true
}

func One2False[T any](T) bool {
	return false
}
