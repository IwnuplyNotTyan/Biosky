//go:build bio_moon

package themes

import "biosky/internal/bio"

type moonTheme struct{}

var moonLines = [8]string{
	"starting over quietly",
	"growing into something",
	"half here, half elsewhere",
	"almost ready to be seen",
	"everything visible, nothing said",
	"letting go slowly",
	"holding on by a thread",
	"between sleep and signal",
}

func (moonTheme) Generate(slot int64) string {
	phase := bio.MoonPhaseAtSlot(slot)
	return phase.Emoji + " " + phase.Name + " · " + moonLines[phase.Index]
}

func init() {
	bio.Register("moon", moonTheme{})
}
