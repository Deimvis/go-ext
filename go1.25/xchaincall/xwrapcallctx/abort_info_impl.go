package xwrapcallctx

// abortInfo is single-threaded implementation of xwrapcall.AbortInfoMutable.
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
