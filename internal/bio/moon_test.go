package bio

import (
	"testing"
	"time"
)

var testReferenceNewMoon = time.Date(2000, 1, 6, 18, 14, 0, 0, time.UTC)

func TestMoonPhaseAtReference(t *testing.T) {
	if got := MoonPhaseAt(testReferenceNewMoon); got.Index != 0 || got.Name != "new moon" {
		t.Fatalf("at reference new moon got %+v", got)
	}
}

func TestMoonPhaseQuarters(t *testing.T) {
	cases := []struct {
		offset time.Duration
		index  int
		name   string
	}{
		{738226 * time.Second, 2, "first quarter"},
		{1275875 * time.Second, 4, "full moon"},
		{1913524 * time.Second, 6, "last quarter"},
	}
	for _, tc := range cases {
		got := MoonPhaseAt(testReferenceNewMoon.Add(tc.offset))
		if got.Index != tc.index || got.Name != tc.name {
			t.Errorf("offset %v: got %+v, want index %d %q", tc.offset, got, tc.index, tc.name)
		}
	}
}

func TestMoonPhaseAtSlotIsDeterministic(t *testing.T) {
	a := MoonPhaseAtSlot(20000)
	b := MoonPhaseAtSlot(20000)
	if a != b {
		t.Fatalf("same slot returned %+v and %+v", a, b)
	}
}
