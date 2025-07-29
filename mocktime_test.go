package mocktime

import (
	"testing"
	"time"

	"github.com/nitroshare/compare"
)

func TestMockUnmock(t *testing.T) {
	compare.Compare(t, Now(), time.Time{}, false)
	Mock()
	compare.Compare(t, Now(), time.Time{}, true)
	Unmock()
	compare.Compare(t, Now(), time.Time{}, false)
}

func TestSetAndAdvance(t *testing.T) {
	Mock()
	defer Unmock()
	v := time.Date(2000, time.January, 1, 0, 0, 0, 0, time.UTC)
	Set(v)
	compare.Compare(t, Now(), v, true)
	d := 24 * time.Hour
	Advance(d)
	compare.Compare(t, Now(), v.Add(d), true)
}

func TestAfter(t *testing.T) {
	Mock()
	defer Unmock()
	var (
		n          = Now()
		d          = 1 * time.Second
		chanClosed = make(chan any)
	)
	go func() {
		defer close(chanClosed)
		<-After(d)
	}()
	AdvanceToAfter()
	compare.Compare(t, Now(), n.Add(d), true)
	<-chanClosed
}

func TestAdvanceToAfter(t *testing.T) {
	Mock()
	defer Unmock()
	var (
		d          = 1 * time.Second
		afterChan  = After(d)
		chanClosed = make(chan any)
	)
	go func() {
		defer close(chanClosed)
		<-afterChan
	}()
	AdvanceToAfter()
	<-chanClosed
}
