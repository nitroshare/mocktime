package mocktime

import (
	"time"

	"github.com/nitroshare/golist"
)

type afterChanData struct {
	expiry time.Time
	ch     chan time.Time
	timer  *Timer
}

type newTimerParams struct {
	duration time.Duration
	timer    *Timer
}

type mockLoop struct {
	mockTime           time.Time
	waitingForAfter    bool
	afterChan          golist.List[*afterChanData]
	chanNow            chan any
	chanAfter          chan time.Duration
	chanSet            chan time.Time
	chanAdvance        chan time.Duration
	chanAdvanceToAfter chan any
	chanNewTimer       chan *newTimerParams
	chanStopTimer      chan *Timer
	chanAny            chan any
	chanBool           chan bool
	chanTime           chan time.Time
	chanTimeChan       chan (<-chan time.Time)
	chanTest           chan any
	chanClose          chan any
}

func (m *mockLoop) send(expiry time.Time, ch chan<- time.Time) {
	go func() {
		ch <- expiry
	}()
}

func (m *mockLoop) after(d time.Duration, t *Timer) <-chan time.Time {
	var (
		expiry = m.mockTime.Add(d)
		ch     = make(chan time.Time)
	)
	if m.waitingForAfter {
		m.mockTime = expiry
		m.send(expiry, ch)
		m.waitingForAfter = false
		m.chanAny <- nil
	} else {
		m.afterChan.Add(&afterChanData{
			expiry: expiry,
			ch:     ch,
			timer:  t,
		})
	}
	return ch
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

func (m *mockLoop) stopTimer(t *Timer) bool {
	for e := m.afterChan.Front; e != nil; e = e.Next {
		if t == e.Value.timer {
			m.afterChan.Remove(e)
			return true
		}
	}
	return false
}

func (m *mockLoop) run() {
	defer close(m.chanClose)
	for {
		select {
		case <-m.chanNow:
			m.chanTime <- m.mockTime
		case d := <-m.chanAfter:
			m.chanTimeChan <- m.after(d, nil)
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
				if m.chanTest != nil {
					m.chanTest <- nil
				}
				m.waitingForAfter = true
				continue
			}
			m.mockTime = earliestTime
			m.sendElapsed()
			m.chanAny <- nil
		case p := <-m.chanNewTimer:
			m.chanTimeChan <- m.after(p.duration, p.timer)
		case t := <-m.chanStopTimer:
			m.chanBool <- m.stopTimer(t)
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
		chanNewTimer:       make(chan *newTimerParams),
		chanStopTimer:      make(chan *Timer),
		chanAny:            make(chan any),
		chanBool:           make(chan bool),
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
