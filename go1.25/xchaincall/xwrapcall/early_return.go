package xwrapcall

import "errors"

type EarlyReturnAction[CtxT Context] func(EarlyReturnInfo[CtxT]) (CtxT, error)

type EarlyReturnInfo[CtxT Context] interface {
	StackIndex() StackInd
	Middleware() Middleware[CtxT]
	ReturnedValues() (CtxT, error)
}

func IgnoreEarlyReturn[CtxT Context](info EarlyReturnInfo[CtxT]) (CtxT, error) {
	return info.ReturnedValues()
}

func PanicEarlyReturn[CtxT Context](info EarlyReturnInfo[CtxT]) (CtxT, error) {
	panic(ErrEarlyReturn)
}

var _ EarlyReturnAction[Context] = IgnoreEarlyReturn[Context]
var _ EarlyReturnAction[Context] = PanicEarlyReturn[Context]

func earlyReturnAction_default[CtxT Context](info EarlyReturnInfo[CtxT]) (CtxT, error) {
	c, err := info.ReturnedValues()
	if ac, ok := Context(c).(AbortableContext); ok {
		// Any early return is assumed as abort by default.
		// So we call Abort() if middleware didn't.
		if !ac.Aborted() {
			ac.Abort(WithReason("middleware has not called neither next nor context's Abort()"), WithFields(
				Field{"mw_ind", info.StackIndex()},
				Field{"mw_func", getFuncFullname(info.Middleware())},
				Field{"mw_err", err.Error()},
				// TODO: add callstack of middlewares to fields (mw1 -> mw2 -> hasn't call next + call stack if given)
				// TODO: allow to customize which fields are added here
			))
		}
		return c, err
	}
	if err == nil {
		// We panic if middleware didn't call next and returned nil error to avoid
		// accidental mistakes when one forgot to call next.
		// To ignore this case one may use AbortContext or provide custom early return action (e.g. IgnoreEarlyPanic).
		panic(ErrEarlyReturn)
	}
	return c, err
}

var ErrEarlyReturn = errors.New("middleware has not called next and returned no error")
