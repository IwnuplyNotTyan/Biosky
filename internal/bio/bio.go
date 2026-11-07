package bio

import (
	"math/rand/v2"
	"time"
)

// Theme names accepted by the compiled binary.
const (
	ThemeDefault = "default"
	ThemeKaomoji = "kaomoji"
	ThemeMoon    = "moon"
)

// Theme returns the name of the compiled-in bio.
func Theme() string { return theme }

// Slot converts a time to a whole UTC day since the Unix epoch.
func Slot(t time.Time) int64 {
	return t.UTC().Unix() / 86400
}

// Rand returns a deterministic random source for a slot and theme salt.
func Rand(slot int64, salt uint64) *rand.Rand {
	return rand.New(rand.NewPCG(uint64(slot), salt))
}
