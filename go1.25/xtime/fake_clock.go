//go:build debug

package xtime

import "time"

func NewFakeClock() ConfigurableClock {
	return &fakeClock{now: nil, shift: 0}
}

type fakeClock struct {
	now   *time.Time
	shift time.Duration
}

func (fc *fakeClock) Now() time.Time {
	if fc.now != nil {
		return *fc.now
	}
	return time.Now().Add(fc.shift)
}

func (fc *fakeClock) Shift(d time.Duration) {
	fc.shift += d
}

func (fc *fakeClock) Stop() {
	fc.StopAt(fc.Now())
}

func (fc *fakeClock) StopAt(t time.Time) {
	fc.now = &t
	fc.shift = 0
}

func (fc *fakeClock) Reset() {
	fc.now = nil
	fc.shift = 0
}
