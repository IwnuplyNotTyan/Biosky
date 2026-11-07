package bio

import (
	"strings"
	"testing"
	"time"
)

type fixedGen struct{ text string }

func (f fixedGen) Generate(int64) string { return f.text }

func resetRegistry() { themes = map[string]Generator{} }

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

func TestSelectDefault(t *testing.T) {
	resetRegistry()
	Register(DefaultTheme, fixedGen{"default"})
	sel, err := Select()
	if err != nil {
		t.Fatal(err)
	}
	if sel.Name != DefaultTheme || sel.Generator.Generate(0) != "default" {
		t.Fatalf("unexpected selection: %+v", sel)
	}
}

func TestSelectSingleTheme(t *testing.T) {
	resetRegistry()
	Register(DefaultTheme, fixedGen{"default"})
	Register("kaomoji", fixedGen{"kaomoji"})
	sel, err := Select()
	if err != nil {
		t.Fatal(err)
	}
	if sel.Name != "kaomoji" {
		t.Fatalf("selected %q, want kaomoji", sel.Name)
	}
}

func TestSelectMultipleThemesFails(t *testing.T) {
	resetRegistry()
	Register(DefaultTheme, fixedGen{"default"})
	Register("moon", fixedGen{"moon"})
	Register("mixed", fixedGen{"mixed"})
	_, err := Select()
	if err == nil {
		t.Fatal("expected an error")
	}
	if !strings.Contains(err.Error(), "mixed, moon") {
		t.Fatalf("error %q does not name both themes", err)
	}
}

func TestSelectEmptyFails(t *testing.T) {
	resetRegistry()
	if _, err := Select(); err == nil {
		t.Fatal("expected an error")
	}
}
