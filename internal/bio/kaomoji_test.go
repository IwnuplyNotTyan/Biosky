//go:build kaomoji

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

func TestGenerateReturnsListedKaomoji(t *testing.T) {
	seen := map[string]bool{}
	for _, k := range kaomojiList {
		seen[k] = true
	}
	for slot := int64(0); slot < 100; slot++ {
		if got := Generate(slot); !seen[got] {
			t.Fatalf("Generate returned unlisted kaomoji %q", got)
		}
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
