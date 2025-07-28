package mocktime

import (
	"testing"
	"time"
)

func compare[T comparable](t *testing.T, v1, v2 T, same bool) {
	t.Helper()
	if same {
		if v1 != v2 {
			t.Fatalf("%v != %v", v1, v2)
		}
	} else {
		if v1 == v2 {
			t.Fatalf("%v == %v", v1, v2)
		}
	}
}

var (
	fireMap = map[<-chan time.Time]bool{}
)

func init() {
	firedFn = func(ch <-chan time.Time) {
		fireMap[ch] = true
	}
}

func assertFired(t *testing.T, ch <-chan time.Time, shouldSucceed bool) {
	t.Helper()
	v, _ := fireMap[ch]
	if shouldSucceed && !v {
		t.Fatal("expected to be able to read from channel")
	}
	if !shouldSucceed && v {
		t.Fatal("unexpecedtly able to read from channel")
	}
}

func TestMockUnmock(t *testing.T) {
	compare(t, Now(), time.Time{}, false)
	Mock()
	compare(t, Now(), time.Time{}, true)
	Unmock()
	compare(t, Now(), time.Time{}, false)
}

func TestSetAdvance(t *testing.T) {
	Mock()
	defer Unmock()
	v := time.Date(2000, time.January, 1, 0, 0, 0, 0, time.UTC)
	Set(v)
	compare(t, Now(), v, true)
	d := 24 * time.Hour
	Advance(d)
	compare(t, Now(), v.Add(d), true)
}

func TestAfter(t *testing.T) {
	Mock()
	defer Unmock()
	ch := After(2 * time.Second)
	Advance(1 * time.Second)
	assertFired(t, ch, false)
	Advance(2 * time.Second)
	assertFired(t, ch, true)
}
