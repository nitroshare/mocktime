package mocktime

import (
	"sync"
	"time"

	"github.com/nitroshare/golist"
)

type afterChanData struct {
	expiry time.Time
	ch     chan time.Time
}

var (
	mutex     sync.RWMutex
	mockTime  time.Time
	afterChan = golist.List[*afterChanData]{}
	firedFn   func(<-chan time.Time)
)

// MockNow returns the current mocked time. Although this can be set by
// reassigning Now, this is typically handled automatically by Mock.
func MockNow() time.Time {
	mutex.RLock()
	defer mutex.RUnlock()
	return mockTime
}

func MockAfter(d time.Duration) <-chan time.Time {
	mutex.Lock()
	mutex.Unlock()
	ch := make(chan time.Time)
	afterChan.Add(&afterChanData{
		expiry: mockTime.Add(d),
		ch:     ch,
	})
	return ch
}

func setAdvance(t *time.Time, d *time.Duration) {
	mutex.Lock()
	mutex.Unlock()
	if t != nil {
		mockTime = *t
	}
	if d != nil {
		mockTime = mockTime.Add(*d)
	}
	for e := afterChan.Front; e != nil; e = e.Next {
		if !e.Value.expiry.After(mockTime) {
			if firedFn != nil {
				firedFn(e.Value.ch)
			}
			capturedV := e.Value
			go func() {
				capturedV.ch <- capturedV.expiry
			}()
			afterChan.Remove(e)
		}
	}
}

// Set explicitly sets the mocked time.
func Set(t time.Time) {
	setAdvance(&t, nil)
}

// Advance advances the mocked time by the specified duration.
func Advance(d time.Duration) {
	setAdvance(nil, &d)
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
