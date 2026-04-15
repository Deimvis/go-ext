package xtime

import "time"

var StdClock Clock = NewStdClock()

func NewStdClock() Clock {
	return stdClock{}
}

type stdClock struct{}

var _ Clock = stdClock{}

func (std stdClock) Now() time.Time {
	return time.Now()
}
