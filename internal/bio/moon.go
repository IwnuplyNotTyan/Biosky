//go:build moon

package bio

import (
	"math"
	"time"
)

const theme = ThemeMoon

const synodicMonthDays = 29.530588853

var referenceNewMoon = time.Date(2000, 1, 6, 18, 14, 0, 0, time.UTC)

var moonPhases = [8]MoonPhase{
	{0, "🌑"},
	{1, "🌒"},
	{2, "🌓"},
	{3, "🌔"},
	{4, "🌕"},
	{5, "🌖"},
	{6, "🌗"},
	{7, "🌘"},
}

// MoonPhase is one of eight phases with its display name and emoji.
type MoonPhase struct {
	Index int
	Emoji string
}

// MoonPhaseAt returns the moon phase for a moment in time.
func MoonPhaseAt(t time.Time) MoonPhase {
	days := t.UTC().Sub(referenceNewMoon).Hours() / 24
	frac := math.Mod(days/synodicMonthDays, 1)
	if frac < 0 {
		frac++
	}
	index := int(math.Floor(frac*8+0.5)) % 8
	return moonPhases[index]
}

// MoonPhaseAtSlot returns the moon phase for a UTC day slot.
func MoonPhaseAtSlot(slot int64) MoonPhase {
	return MoonPhaseAt(time.Unix(slot*86400, 0).UTC())
}

// Generate returns the moon phase and its line for the UTC day slot.
func Generate(slot int64) string {
	phase := MoonPhaseAtSlot(slot)
	return phase.Emoji
}
