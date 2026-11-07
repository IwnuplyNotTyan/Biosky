package bio

import (
	_ "embed"
	"encoding/json"
	"math"
	"math/rand/v2"
	"time"
)

// Theme names accepted by the compiled binary.
const (
	ThemeDefault = "default"
	ThemeKaomoji = "kaomoji"
	ThemeMoon    = "moon"
)

// theme is the bio this binary publishes. Change it and rebuild.
const theme = ThemeKaomoji

func init() {
	switch theme {
	case ThemeDefault, ThemeKaomoji, ThemeMoon:
	default:
		panic("bio: unknown theme " + theme)
	}
}

// Theme returns the name of the compiled-in bio.
func Theme() string { return theme }

// kaomoji.json is a flat list derived from
// https://github.com/asakura89/Minimanimoji, which is released under the
// Unlicense (public domain).
//
//go:embed kaomoji.json
var kaomojiJSON []byte

var kaomojiList = loadKaomoji()

func loadKaomoji() []string {
	var list []string
	if err := json.Unmarshal(kaomojiJSON, &list); err != nil {
		panic("bio: invalid kaomoji.json: " + err.Error())
	}
	if len(list) == 0 {
		panic("bio: kaomoji.json is empty")
	}
	return list
}

const defaultText = "quietly existing, one day at a time"

func defaultBio() string { return defaultText }

func kaomojiBio(slot int64) string {
	return kaomojiList[Rand(slot, 0x6b616f6d6f6a69).IntN(len(kaomojiList))]
}

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

func moonBio(slot int64) string {
	phase := MoonPhaseAtSlot(slot)
	return phase.Emoji + " " + phase.Name + " · " + moonLines[phase.Index]
}

// Generate returns the bio for a UTC day slot.
func Generate(slot int64) string {
	switch theme {
	case ThemeKaomoji:
		return kaomojiBio(slot)
	case ThemeMoon:
		return moonBio(slot)
	default:
		return defaultBio()
	}
}

// Slot converts a time to a whole UTC day since the Unix epoch.
func Slot(t time.Time) int64 {
	return t.UTC().Unix() / 86400
}

// Rand returns a deterministic random source for a slot and theme salt.
func Rand(slot int64, salt uint64) *rand.Rand {
	return rand.New(rand.NewPCG(uint64(slot), salt))
}

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
