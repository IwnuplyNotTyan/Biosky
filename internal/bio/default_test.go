//go:build !kaomoji && !moon && !iwnuply

package bio

import "testing"

func TestGenerateDefault(t *testing.T) {
	if got := Generate(42); got != defaultText {
		t.Fatalf("Generate = %q, want %q", got, defaultText)
	}
}
