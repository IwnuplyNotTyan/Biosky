package bio

import (
	"testing"
	"time"
)

func TestSlot(t *testing.T) {
	ts := time.Date(2026, 11, 7, 23, 59, 59, 0, time.UTC)
	got := Slot(ts)
	want := ts.Unix() / 86400
	if got != want {
		t.Fatalf("Slot(%v) = %d, want %d", ts, got, want)
	}
}

func TestRandIsDeterministicPerSlot(t *testing.T) {
	a := Rand(123, 0xabc)
	b := Rand(123, 0xabc)
	for i := 0; i < 20; i++ {
		if x, y := a.IntN(1000), b.IntN(1000); x != y {
			t.Fatalf("draw %d differs: %d != %d", i, x, y)
		}
	}
}
