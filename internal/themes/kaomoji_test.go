//go:build bio_kaomoji

package themes

import (
	"testing"

	"biosky/internal/bio"
)

func TestKaomojiSelection(t *testing.T) {
	sel, err := bio.Select()
	if err != nil {
		t.Fatal(err)
	}
	if sel.Name != "kaomoji" {
		t.Fatalf("selected %q, want kaomoji", sel.Name)
	}
}

func TestKaomojiDeterministicAndVaried(t *testing.T) {
	g := kaomoji{}
	if g.Generate(42) != g.Generate(42) {
		t.Fatal("same slot produced different bios")
	}
	seen := map[string]bool{}
	for slot := int64(0); slot < 100; slot++ {
		seen[g.Generate(slot)] = true
	}
	if len(seen) < 2 {
		t.Fatal("100 slots produced fewer than two distinct bios")
	}
}
