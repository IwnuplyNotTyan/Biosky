package bio

import (
	"math"
	"time"
)

const synodicMonthDays = 29.530588853

var referenceNewMoon = time.Date(2000, 1, 6, 18, 14, 0, 0, time.UTC)

var moonPhases = [8]MoonPhase{
	{0, "new moon", "🌑"},
	{1, "waxing crescent", "🌒"},
	{2, "first quarter", "🌓"},
	{3, "waxing gibbous", "🌔"},
	{4, "full moon", "🌕"},
	{5, "waning gibbous", "🌖"},
	{6, "last quarter", "🌗"},
	{7, "waning crescent", "🌘"},
}

// MoonPhase is one of eight phases with its display name and emoji.
type MoonPhase struct {
	Index int
	Name  string
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
