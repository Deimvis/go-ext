package xstrconv

type SliceFormatter interface {
	Add(string)
	Format() (string, error)
}

func NewIncrementalSliceFormatter(addFn func(s string, incr string) string, formatFn func(s string) string) SliceFormatter {
	return &incrementalSliceFormatter{s: "", addFn: addFn, formatFn: formatFn}
}

type incrementalSliceFormatter struct {
	s string

	addFn    func(s string, incr string) string
	formatFn func(s string) string
}

func (isf *incrementalSliceFormatter) Add(incr string) {
	isf.s = isf.addFn(isf.s, incr)
}

func (isf *incrementalSliceFormatter) Format() (string, error) {
	return isf.formatFn(isf.s), nil
}
