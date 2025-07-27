package mocktime

import (
	"time"
)

type afterChanData struct {
	expiry time.Time
	ch     chan<- time.Time
}

var (
	mockTime  time.Time
	afterChan []*afterChanData
)

// MockNow returns the current mocked time. Although this can be set by
// reassigning Now, this is typically handled automatically by Mock.
func MockNow() time.Time {
	return mockTime
}

func MockAfter(d time.Duration) <-chan time.Time {
	ch := make(chan time.Time)
	afterChan = append(afterChan, &afterChanData{
		expiry: mockTime.Add(d),
		ch:     ch,
	})
	return ch
}

// Set explicitly sets the mocked time.
func Set(t time.Time) {
	mockTime = t
	expInd := 0
	for i, v := range afterChan {
		if v.expiry.After(mockTime) {
			expInd = i
			break
		}
		go func() {
			v.ch <- v.expiry
		}()
	}
	afterChan = afterChan[expInd:]
}

// Advance advances the mocked time by the specified duration.
func Advance(d time.Duration) {
	Set(mockTime.Add(d))
}

var (

	// Now normally points to time.Now but can also be pointed to with MockNow.
	Now func() time.Time

	// After normally points to time.After but can also be pointed to with MockAfter.
	After func(d time.Duration) <-chan time.Time
)

func set() {
	Now = MockNow
	After = MockAfter
}

func reset() {
	Now = time.Now
	After = time.After
}

func init() {
	reset()
}

// Mock replaces the time functions in this package with their mocked equivalents.
func Mock() {
	set()
}

// Unmock replaces the mocked time functions with their original equivalents.
func Unmock() {
	reset()
}
