package xio

import "io"

type ReadFn func([]byte) (int, error)
type SeekFn func(offset int64, whence int) (int64, error)

func NewReader(fn ReadFn) io.Reader {
	return &readerImpl{fn: fn}
}

func NewSeeker(fn SeekFn) io.Seeker {
	return &seekerImpl{fn: fn}
}

// NewReadSeeker constructs io.ReadSeeker from given functions.
// Obviously, these functions must work over the same data.
func NewReadSeeker(rfn ReadFn, sfn SeekFn) io.ReadSeeker {
	return &readSeekerImpl{rfn: rfn, sfn: sfn}
}

type readerImpl struct {
	fn ReadFn
}

func (r *readerImpl) Read(p []byte) (int, error) {
	return r.fn(p)
}

type seekerImpl struct {
	fn SeekFn
}

func (s *seekerImpl) Seek(offset int64, whence int) (int64, error) {
	return s.fn(offset, whence)
}

type readSeekerImpl struct {
	rfn ReadFn
	sfn SeekFn
}

func (rs *readSeekerImpl) Read(p []byte) (int, error) {
	return rs.rfn(p)
}

func (rs *readSeekerImpl) Seek(offset int64, whence int) (int64, error) {
	return rs.sfn(offset, whence)
}
