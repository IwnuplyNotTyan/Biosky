//go:build bio_mixed

package themes

import (
	"strings"
	"testing"

	"biosky/internal/bio"
)

func TestMixedThemeUsesKaomojiAndPhase(t *testing.T) {
	const slot = int64(42)
	out := mixedTheme{}.Generate(slot)
	phase := bio.MoonPhaseAtSlot(slot)
	if !strings.Contains(out, phase.Emoji) {
		t.Fatalf("bio %q does not contain the moon emoji %q", out, phase.Emoji)
	}
	found := false
	for _, k := range mixedKaomoji {
		if strings.Contains(out, k) {
			found = true
			break
		}
	}
	if !found {
		t.Fatalf("bio %q contains no kaomoji", out)
	}
}
