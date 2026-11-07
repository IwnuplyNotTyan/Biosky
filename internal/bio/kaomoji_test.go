//go:build bio_kaomoji

package bio

import (
	"strings"
	"testing"
)

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

func TestTheme(t *testing.T) {
	if theme != ThemeKaomoji {
		t.Fatalf("theme = %q, want %q", theme, ThemeKaomoji)
	}
	if Theme() != ThemeKaomoji {
		t.Fatalf("Theme() = %q", Theme())
	}
}
