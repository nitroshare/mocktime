package mocktime

import (
	"time"
)

// Timer provides a drop-in replacement for time.Timer.
type Timer struct {
	C     <-chan time.Time
	timer *time.Timer
}

// Stop ends the timer.
func (t *Timer) Stop() bool {
	if t.timer != nil {
		return t.timer.Stop()
	}
	loop.chanStopTimer <- t
	return <-loop.chanBool
}

// NewTimer creates a new Timer that will send on the channel C when the
// specified duration expires.
func NewTimer(d time.Duration) *Timer {
	if loop == nil {
		timer := time.NewTimer(d)
		return &Timer{
			C:     timer.C,
			timer: timer,
		}
	}
	t := &Timer{}
	loop.chanNewTimer <- &newTimerParams{
		duration: d,
		timer:    t,
	}
	t.C = <-loop.chanTimeChan
	return t
}
