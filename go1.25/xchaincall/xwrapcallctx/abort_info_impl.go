package xwrapcallctx

import "github.com/Deimvis/go-ext/go1.25/xchaincall/xwrapcall"

// abortInfo is single-threaded implementation of xwrapcall.AbortInfoMutable.
type abortInfo struct {
	reason string
	fields []xwrapcall.Field
}

var _ xwrapcall.AbortInfoMutable = &abortInfo{}

func (i *abortInfo) Reason() string {
	return i.reason
}

func (i *abortInfo) SetReason(r string) {
	i.reason = r
}

func (i *abortInfo) Fields() []xwrapcall.Field {
	return i.fields
}

func (i *abortInfo) SetFields(fields ...xwrapcall.Field) {
	i.fields = fields
}
