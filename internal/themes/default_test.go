package themes

import "testing"

func TestStaticGenerate(t *testing.T) {
	g := static("hello")
	for _, slot := range []int64{0, 1, 99999} {
		if got := g.Generate(slot); got != "hello" {
			t.Fatalf("slot %d: got %q", slot, got)
		}
	}
}
