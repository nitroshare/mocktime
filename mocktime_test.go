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
	select {
	case <-ch:
		t.Fatalf("unexpected read on channel")
	case <-time.After(10 * time.Millisecond):
	}
	Advance(2 * time.Second)
	select {
	case <-ch:
	case <-time.After(10 * time.Millisecond):
		t.Fatalf("unexpected block on channel")
	}
}
