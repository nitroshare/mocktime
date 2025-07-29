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
// of the next After() call, waiting for a call to After() if necessary.
func AdvanceToAfter() {
	loop.chanAdvanceToAfter <- nil
	<-loop.chanAny
}

var (
	// Now normally points to time.Now but can also be pointed to with MockNow.
	Now func() time.Time

	// After normally points to time.After but can also be pointed to with MockAfter.
	After func(d time.Duration) <-chan time.Time
)

func reset() {
	Now = time.Now
	After = time.After
}

func init() {
	reset()
}

// Mock replaces the time functions in this package with their mocked equivalents.
func Mock() {
	loop = newMockLoop()
	Now = MockNow
	After = MockAfter
}

// Unmock replaces the mocked time functions with their original equivalents.
func Unmock() {
	loop.close()
	reset()
}
