package mocktime

import (
	"time"

	"github.com/nitroshare/golist"
)

type afterChanData struct {
	expiry time.Time
	ch     chan time.Time
}

type mockLoop struct {
	mockTime           time.Time
	afterChan          golist.List[*afterChanData]
	chanNow            chan any
	chanAfter          chan time.Duration
	chanSet            chan time.Time
	chanAdvance        chan time.Duration
	chanAdvanceToAfter chan any
	chanAny            chan any
	chanTime           chan time.Time
	chanTimeChan       chan (<-chan time.Time)
	chanClose          chan any
}

func (m *mockLoop) send(expiry time.Time, ch chan<- time.Time) {
	go func() {
		ch <- expiry
	}()
}

func (m *mockLoop) sendElapsed() {
	for e := m.afterChan.Front; e != nil; e = e.Next {
		if !e.Value.expiry.After(m.mockTime) {
			m.send(e.Value.expiry, e.Value.ch)
			m.afterChan.Remove(e)
		}
	}
}

func (m *mockLoop) findEarliest() time.Time {
	var earliestTime time.Time
	for e := m.afterChan.Front; e != nil; e = e.Next {
		if earliestTime.IsZero() || e.Value.expiry.Before(earliestTime) {
			earliestTime = e.Value.expiry
		}
	}
	return earliestTime
}

func (m *mockLoop) run() {
	defer close(m.chanClose)
	waitingForAfter := false
	for {
		select {
		case <-m.chanNow:
			m.chanTime <- m.mockTime
		case d := <-m.chanAfter:
			var (
				expiry = m.mockTime.Add(d)
				ch     = make(chan time.Time)
			)
			if waitingForAfter {
				m.mockTime = expiry
				m.send(expiry, ch)
				waitingForAfter = false
				m.chanAny <- nil
			} else {
				m.afterChan.Add(&afterChanData{
					expiry: expiry,
					ch:     ch,
				})
			}
			m.chanTimeChan <- ch
		case t := <-m.chanSet:
			m.mockTime = t
			m.sendElapsed()
			m.chanAny <- nil
		case d := <-m.chanAdvance:
			m.mockTime = m.mockTime.Add(d)
			m.sendElapsed()
			m.chanAny <- nil
		case <-m.chanAdvanceToAfter:
			earliestTime := m.findEarliest()
			if earliestTime.IsZero() {
				waitingForAfter = true
				continue
			}
			m.mockTime = earliestTime
			m.sendElapsed()
			m.chanAny <- nil
		case <-m.chanClose:
			return
		}
	}
}

func newMockLoop() *mockLoop {
	m := &mockLoop{
		afterChan:          golist.List[*afterChanData]{},
		chanNow:            make(chan any),
		chanAfter:          make(chan time.Duration),
		chanSet:            make(chan time.Time),
		chanAdvance:        make(chan time.Duration),
		chanAdvanceToAfter: make(chan any),
		chanAny:            make(chan any),
		chanTime:           make(chan time.Time),
		chanTimeChan:       make(chan (<-chan time.Time)),
		chanClose:          make(chan any),
	}
	go m.run()
	return m
}

func (m *mockLoop) close() {
	m.chanClose <- nil
	<-m.chanClose
}
