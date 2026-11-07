//go:build bio_mixed

package themes

import "biosky/internal/bio"

type mixedTheme struct{}

var mixedKaomoji = []string{
	"(´-ω-`)",
	"(｡•̀ᴗ-)✧",
	"(￣▽￣)ノ",
	"( •_•)>⌐■-■",
}

var mixedLines = [8]string{
	"new moon, new excuse to disappear",
	"waxing crescent, quietly becoming",
	"first quarter, half committed",
	"waxing gibbous, almost legible",
	"full moon, nowhere to hide",
	"waning gibbous, shedding signal",
	"last quarter, holding on by a thread",
	"waning crescent, almost offline",
}

func (mixedTheme) Generate(slot int64) string {
	r := bio.Rand(slot, 0x6d69786564)
	k := mixedKaomoji[r.IntN(len(mixedKaomoji))]
	phase := bio.MoonPhaseAtSlot(slot)
	return k + " " + phase.Emoji + " " + mixedLines[phase.Index]
}

func init() {
	bio.Register("mixed", mixedTheme{})
}
