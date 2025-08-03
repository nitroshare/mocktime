package mocktime

import (
	"time"
)

var (
	loop *mockLoop
)

// MockNow returns the current mocked time. Although this can be set by
// reassigning Now, this is typically handled automatically by Mock.
func MockNow() time.Time {
	loop.chanNow <- nil
	return <-loop.chanTime
}

// MockAfter returns a channel that will send when the specified duration
// elapses via calls to Set or Advance.
func MockAfter(d time.Duration) <-chan time.Time {
	if loop.chanTest != nil {
		<-loop.chanTest
	}
	loop.chanAfter <- d
	return <-loop.chanTimeChan
}

// Set explicitly sets the mocked time.
func Set(t time.Time) {
	loop.chanSet <- t
	<-loop.chanAny
}

// Advance advances the mocked time by the specified duration.
func Advance(d time.Duration) {
	loop.chanAdvance <- d
	<-loop.chanAny
}

// AdvanceToAfter is a synchronization function that will advance to the time
// of the next After() call, waiting for a call to After() if necessary. Note
// that this will also include timers.
func AdvanceToAfter() {
	loop.chanAdvanceToAfter <- nil
	<-loop.chanAny
}

var (
	// Now normally points to time.Now but can also be pointed to with
	// MockNow.
	Now func() time.Time

	// After normally points to time.After but can also be pointed to with
	// MockAfter.
	After func(d time.Duration) <-chan time.Time
)

func reset() {
	Now = time.Now
	After = time.After
}

func init() {
	reset()
}

// Mock replaces the time functions in this package with their mocked
// equivalents. Note that this is intended to be called at the beginning of a
// test and used with an associated defer Unmock() call. It is not safe for
// Mock or Unmock to be called from other goroutines, although the other
// mocked functions and types may be used anywhere while mocking is active.
func Mock() {
	loop = newMockLoop()
	Now = MockNow
	After = MockAfter
}

// Unmock replaces the mocked time functions with their original equivalents.
func Unmock() {
	loop.close()
	loop = nil
	reset()
}

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
// specified duration expires. When mocked,
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
