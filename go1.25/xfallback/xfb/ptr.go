package xfb

// TODO: rename to
// - OnNil + OnNilDeref (?)
// - OnNil + OnNilv (?)
// - Ptr + Deref (?)

func OnNil[T any](v *T, fallback *T) *T {
	if v != nil {
		return v
	}
	return fallback
}

func OnNilv[T any](v *T, fallback T) T {
	if v != nil {
		return *v
	}
	return fallback
}
