package xtime

import "time"

type Clock interface {
	Now() time.Time
}

type ConfigurableClock interface {
	Clock
	Shift(time.Duration)
	Stop()
	StopAt(time.Time)
	Reset()
}


