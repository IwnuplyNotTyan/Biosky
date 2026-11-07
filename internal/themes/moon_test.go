//go:build bio_moon

package themes

import (
	"strings"
	"testing"

	"biosky/internal/bio"
)

func TestMoonThemeUsesPhase(t *testing.T) {
	const slot = int64(42)
	out := moonTheme{}.Generate(slot)
	phase := bio.MoonPhaseAtSlot(slot)
	if !strings.Contains(out, phase.Emoji) || !strings.Contains(out, phase.Name) {
		t.Fatalf("bio %q does not contain phase %+v", out, phase)
	}
}
