package xwrapcall

import (
	"context"
	"errors"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNew(t *testing.T) {
	type C = Context
	t.Run("smoke", func(t *testing.T) {
		fn := func(c Context) error { return nil }
		fn = New[Context]().Do(fn)
		err := fn(context.Background())
		require.NoError(t, err)
	})
	t.Run("action/invoked", func(t *testing.T) {
		invoked := false
		fn := func(c Context) error {
			invoked = true
			return nil
		}
		fn = New[Context]().Do(fn)
		err := fn(context.Background())
		require.NoError(t, err)
		require.True(t, invoked)
	})
	t.Run("action/error-passed", func(t *testing.T) {
		e := errors.New("123")
		fn := func(c Context) error {
			return e
		}
		fn = New[Context]().Do(fn)
		err := fn(context.Background())
		require.Equal(t, e, err)
	})
	t.Run("middlewares/smoke", func(t *testing.T) {
		fn := func(c Context) error { return nil }
		fn = New[Context]().
			With(
				func(c Context, next Next[Context]) (Context, error) {
					return next(c)
				},
			).
			Do(fn)
		err := fn(context.Background())
		require.NoError(t, err)
	})
	t.Run("middlewares/invoked", func(t *testing.T) {
		invoked := false
		fn := func(c Context) error { return nil }
		fn = New[Context]().
			With(
				func(c Context, next Next[Context]) (Context, error) {
					invoked = true
					return next(c)
				},
			).
			Do(fn)
		err := fn(context.Background())
		require.NoError(t, err)
		require.True(t, invoked)
	})
	t.Run("middlewares/error-passed", func(t *testing.T) {
		e := errors.New("middleware error")
		fn := func(c Context) error { return nil }
		fn = New[Context]().
			With(
				func(c Context, next Next[Context]) (Context, error) {
					return c, e
				},
			).
			Do(fn)
		err := fn(context.Background())
		require.Equal(t, e, err)
	})
	t.Run("middlewares/context-passed", func(t *testing.T) {
		type key struct{}
		value := "value"
		fn := func(c Context) error { return nil }
		fn = New[Context]().
			With(
				func(c Context, next Next[Context]) (Context, error) {
					c = context.WithValue(c, key{}, value)
					return next(c)
				},
				func(c Context, next Next[Context]) (Context, error) {
					require.Equal(t, value, c.Value(key{}))
					return next(c)
				},
			).
			Do(fn)
		err := fn(context.Background())
		require.NoError(t, err)
	})
	t.Run("middlewares/early-return-default", func(t *testing.T) {
		t.Run("pass-error", func(t *testing.T) {
			e := errors.New("early return")
			fn := func(c Context) error { return nil }
			fn = New[Context]().
				With(
					func(c Context, next Next[Context]) (Context, error) {
						return c, e
					},
				).
				Do(fn)
			err := fn(context.Background())
			require.Equal(t, e, err)
		})
		t.Run("panic-on-nil-error", func(t *testing.T) {
			fn := func(c Context) error { return nil }
			fn = New[Context]().
				With(
					func(c Context, next Next[Context]) (Context, error) {
						return c, nil
					},
				).
				Do(fn)
			func() {
				defer func() {
					r := recover()
					require.NotNil(t, r)
				}()
				_ = fn(context.Background())
			}()
		})
		t.Run("abortable-ctx-pass-abort", func(t *testing.T) {
			reason := "my reason"
			fn := func(c AbortableContext) error { return nil }
			fn = New[AbortableContext]().
				With(
					func(c AbortableContext, next Next[AbortableContext]) (AbortableContext, error) {
						c, err := next(c)
						require.True(t, c.Aborted())
						require.Equal(t, reason, c.AbortInfo().Reason())
						return c, err
					},
					func(c AbortableContext, next Next[AbortableContext]) (AbortableContext, error) {
						c.Abort(WithReason(reason))
						return c, nil
					},
				).
				Do(fn)
			err := fn(actxNew())
			require.NoError(t, err)
		})
		t.Run("abortable-ctx-auto-abort", func(t *testing.T) {
			e := errors.New("early return")
			fn := func(c AbortableContext) error { return nil }
			fn = New[AbortableContext]().
				With(
					func(c AbortableContext, next Next[AbortableContext]) (AbortableContext, error) {
						c, err := next(c)
						require.True(t, c.Aborted())
						require.Equal(t, "middleware has not called neither next nor context's Abort()", c.AbortInfo().Reason())
						fs := c.AbortInfo().Fields()
						fkey2value := make(map[string]any)
						for _, f := range fs {
							fkey2value[f.Key] = f.Value
						}
						require.Equal(t, int64(1), fkey2value["mw_ind"].(StackInd))
						require.True(t, strings.HasSuffix(fkey2value["mw_func"].(string), "xchaincall/xwrapcall.TestNew.func8.4.newMwEarlyReturnNoAbort[...].3"))
						require.Equal(t, e.Error(), fkey2value["mw_err"].(string))
						return c, err
					},
					newMwEarlyReturnNoAbort[AbortableContext](e),
				).
				Do(fn)
			err := fn(actxNew())
			require.Equal(t, e, err)
		})
	})
	t.Run("middlewares/early-return-custom", func(t *testing.T) {
		mwErr := errors.New("early return")
		eaErr := errors.New("early return action called")
		fn := New[Context]().
			With(
				func(c Context, next Next[Context]) (Context, error) {
					return c, mwErr
				},
			).
			OnEarlyReturnDo(func(info EarlyReturnInfo[Context]) (Context, error) {
				c, err := info.ReturnedValues()
				require.Equal(t, mwErr, err)
				return c, eaErr
			}).
			Do(noopAction)
		err := fn(context.Background())
		require.Equal(t, eaErr, err)
	})
	t.Run("middlewares/call-next-twice", func(t *testing.T) {
		fn := New[Context]().
			With(
				func(c Context, next Next[Context]) (Context, error) {
					next(c)
					next(c)
					return c, nil
				},
			).
			Do(noopAction)
		require.Panics(t, func() {
			_ = fn(context.Background())
		})
	})
	t.Run("exporting-debug-info/non-nil", func(t *testing.T) {
		var d Debug = nil
		b := New[C]().
			ExportingDebugInfo(&d)
		require.Nil(t, d)
		_ = b.Do(noopAction)
		require.NotNil(t, d)
	})
	t.Run("abortable-ctx-smoke", func(t *testing.T) {
		fn := func(c AbortableContext) error { return nil }
		fn = New[AbortableContext]().Do(fn)
		err := fn(actxNew())
		require.NoError(t, err)
	})
	t.Run("abortable-ctx-action-abort", func(t *testing.T) {
		reason := "action aborted"
		fn := func(c AbortableContext) error {
			c.Abort(WithReason(reason))
			return nil
		}
		fn = New[AbortableContext]().
			With(
				func(c AbortableContext, next Next[AbortableContext]) (AbortableContext, error) {
					c, err := next(c)
					require.True(t, c.Aborted())
					require.Equal(t, reason, c.AbortInfo().Reason())
					return c, err
				},
			).
			Do(fn)
		err := fn(actxNew())
		require.NoError(t, err)
	})
	t.Run("abortable-ctx-middleware-abort-panics-on-next", func(t *testing.T) {
		reason := "middleware aborted"
		fn := func(c AbortableContext) error { return nil }
		fn = New[AbortableContext]().
			With(
				func(c AbortableContext, next Next[AbortableContext]) (retC AbortableContext, retErr error) {
					c.Abort(WithReason(reason))
					defer func() {
						r := recover()
						require.NotNil(t, r)
						require.Equal(t, ErrNextAfterAbort, r.(error))
						retC = c
						retErr = nil
					}()
					return next(c)
				},
			).
			Do(fn)
		ctx := actxNew()
		err := fn(ctx)
		require.NoError(t, err)
	})
}

func actxNew() *actx {
	return &actx{Context: context.Background(), abortInfo: &abortInfo{}}
}

// actx is primitive implementation of AbortableContext
type actx struct {
	context.Context

	aborted   bool
	abortInfo *abortInfo
}

func (c *actx) Abort(opts ...AbortOption) {
	c.aborted = true
	for _, opt := range opts {
		opt(c.abortInfo)
	}
}

func (c *actx) Aborted() bool {
	return c.aborted
}

func (c *actx) AbortInfo() AbortInfo {
	return c.abortInfo
}

type abortInfo struct {
	reason string
	fields []Field
}

var _ AbortInfoMutable = &abortInfo{}

func (i *abortInfo) Reason() string {
	return i.reason
}

func (i *abortInfo) SetReason(r string) {
	i.reason = r
}

func (i *abortInfo) Fields() []Field {
	return i.fields
}

func (i *abortInfo) SetFields(fields ...Field) {
	i.fields = fields
}

func noopAction(c Context) error {
	return nil
}

func newMwEarlyReturnNoAbort[CtxT AbortableContext](e error) Middleware[CtxT] {
	return func(c CtxT, next Next[CtxT]) (CtxT, error) {
		return c, e
	}
}
