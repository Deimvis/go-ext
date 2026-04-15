package xwrapcall

import "fmt"

type Action[CtxT Context] func(CtxT) error
type Middleware[CtxT Context] func(CtxT, Next[CtxT]) (CtxT, error)
type Next[CtxT Context] func(CtxT) (CtxT, error)

//              ∧  | Silence()
// Generalize() |  ∨

type SilentAction[CtxT Context] func(CtxT)
type SilentMiddleware[CtxT Context] func(CtxT, SilentNext[CtxT]) CtxT
type SilentNext[CtxT Context] func(CtxT) CtxT

// Generalize converts SilentAction -> Action
func (sa SilentAction[CtxT]) Generalize() Action[CtxT] {
	return func(c CtxT) error {
		sa(c)
		return nil
	}
}

// Generalize converts SilentMiddleware -> Middleware
func (smw SilentMiddleware[CtxT]) Generalize() Middleware[CtxT] {
	return func(c CtxT, next Next[CtxT]) (CtxT, error) {
		return smw(c, next.Silence()), nil
	}
}

// Generalize converts SilentNext -> Next
func (sn SilentNext[CtxT]) Generalize() Next[CtxT] {
	return func(c CtxT) (CtxT, error) {
		return sn(c), nil
	}
}

// Silence converts Action -> SilentAction
func (a Action[CtxT]) Silence() SilentAction[CtxT] {
	return func(c CtxT) {
		err := a(c)
		checkSilencedError(err)
	}
}

// Silence converts Middleware -> SilentAction
func (mw Middleware[CtxT]) Silence() SilentMiddleware[CtxT] {
	return func(c CtxT, next SilentNext[CtxT]) CtxT {
		c, err := mw(c, next.Generalize())
		checkSilencedError(err)
		return c
	}
}

// Silence converts Next -> SilentNext
func (n Next[CtxT]) Silence() SilentNext[CtxT] {
	return func(c CtxT) CtxT {
		c, err := n(c)
		checkSilencedError(err)
		return c
	}
}

func checkSilencedError(err error) {
	if err != nil {
		panic(fmt.Errorf(errErrorSilencedFmt, err))
	}
}

var (
	errErrorSilencedFmt = "silenced error: %w"
)
