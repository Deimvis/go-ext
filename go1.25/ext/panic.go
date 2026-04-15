package ext

import (
	"fmt"
	"log"
)

// OnPanic runs the given function only if code is paniced.
// It doesn't stop panic, it will continue panic after function is executed.
// First argument of wrapped function is a result of recover.
// NOTE: it must be called directly by defer statement, otherwise panic won't be recovered: https://stackoverflow.com/a/49344592.
func OnPanic(fn func(r any)) {
	// starting from Golang 1.21 panic(nil) doesn't return nil on recover: https://pkg.go.dev/runtime@master#PanicNilError
	if r := recover(); r != nil {
		fn(r)
		panic(r)
	}
}

// OnPanicX runs the given function only if code is paniced.
// It doesn't stop panic, it will continue panic after function is executed.
// First argument of wrapped function is a result of recover.
// If wrapped function returns an error, then new panic combines
// current and new errors.
// NOTE: it must be called directly by defer statement, otherwise panic won't be recovered: https://stackoverflow.com/a/49344592.
func OnPanicX(fn func(r any) error) {
	// starting from Golang 1.21 panic(nil) doesn't return nil on recover: https://pkg.go.dev/runtime@master#PanicNilError
	if r := recover(); r != nil {
		err := fn(r)
		if err != nil {
			switch rr := r.(type) {
			case error:
				r = fmt.Errorf("%w\nwDuring handling of the above exception, another exception occurred:\n%w", rr, err)
			default:
				r = fmt.Errorf("panic(%v)\nwDuring handling of the above exception, another exception occurred:\n%w", rr, err)
			}
		}
		panic(r)
	}
}

func WithRecovery(fn func()) {
	defer func() {
		if r := recover(); r != nil {
			msg := ""
			switch rr := r.(type) {
			case error:
				msg = rr.Error()
			default:
				msg = fmt.Sprintf("%v", rr)
			}
			// TODO: add stack trace
			log.Printf("Recovered from panic: %s\n", msg)
		}
	}()
	fn()
}
