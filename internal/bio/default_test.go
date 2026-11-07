//go:build !kaomoji && !moon

package bio

import "testing"

func TestGenerateDefault(t *testing.T) {
	if got := Generate(42); got != defaultText {
		t.Fatalf("Generate = %q, want %q", got, defaultText)
	}
}

func TestTheme(t *testing.T) {
	if theme != ThemeDefault {
		t.Fatalf("theme = %q, want %q", theme, ThemeDefault)
	}
	if Theme() != ThemeDefault {
		t.Fatalf("Theme() = %q", Theme())
	}
}
