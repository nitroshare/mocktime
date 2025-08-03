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
