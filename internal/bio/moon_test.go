//go:build moon

package bio

import (
	"strings"
	"testing"
	"time"
)

var testReferenceNewMoon = time.Date(2000, 1, 6, 18, 14, 0, 0, time.UTC)

func TestMoonPhaseAtSlotIsDeterministic(t *testing.T) {
	a := MoonPhaseAtSlot(20000)
	b := MoonPhaseAtSlot(20000)
	if a != b {
		t.Fatalf("same slot returned %+v and %+v", a, b)
	}
}

func TestGenerateUsesPhase(t *testing.T) {
	const slot = int64(42)
	out := Generate(slot)
	phase := MoonPhaseAtSlot(slot)
	if !strings.Contains(out, phase.Emoji) {
		t.Fatalf("bio %q does not contain phase %+v", out, phase)
	}
}
