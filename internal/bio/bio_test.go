package bio

import (
	"strings"
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

func TestKaomojiList(t *testing.T) {
	if len(kaomojiList) < 100 {
		t.Fatalf("embedded list has only %d kaomoji", len(kaomojiList))
	}
	seen := map[string]bool{}
	for _, k := range kaomojiList {
		if strings.TrimSpace(k) == "" {
			t.Fatal("embedded list contains an empty kaomoji")
		}
		if seen[k] {
			t.Fatalf("embedded list contains duplicate %q", k)
		}
		seen[k] = true
	}
}

func TestGenerateIsDeterministic(t *testing.T) {
	if Generate(42) != Generate(42) {
		t.Fatal("same slot produced different bios")
	}
}

func TestGenerateVaries(t *testing.T) {
	seen := map[string]bool{}
	for slot := int64(0); slot < 100; slot++ {
		seen[Generate(slot)] = true
	}
	if len(seen) < 2 {
		t.Fatal("100 slots produced fewer than two distinct bios")
	}
}

func TestGenerateMatchesTheme(t *testing.T) {
	const slot = int64(42)
	switch theme {
	case ThemeKaomoji:
		if got, want := Generate(slot), kaomojiBio(slot); got != want {
			t.Fatalf("Generate = %q, want kaomoji %q", got, want)
		}
	case ThemeMoon:
		if got, want := Generate(slot), moonBio(slot); got != want {
			t.Fatalf("Generate = %q, want moon %q", got, want)
		}
	case ThemeDefault:
		if got, want := Generate(slot), defaultBio(); got != want {
			t.Fatalf("Generate = %q, want default %q", got, want)
		}
	default:
		t.Fatalf("unknown theme %q", theme)
	}
}

func TestMoonBioUsesPhase(t *testing.T) {
	const slot = int64(42)
	out := moonBio(slot)
	phase := MoonPhaseAtSlot(slot)
	if !strings.Contains(out, phase.Emoji) || !strings.Contains(out, phase.Name) {
		t.Fatalf("bio %q does not contain phase %+v", out, phase)
	}
}
