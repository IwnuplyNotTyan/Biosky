package bio

import (
	"math/rand/v2"
	"time"
)

// Slot converts a time to a whole UTC day since the Unix epoch.
func Slot(t time.Time) int64 {
	return t.UTC().Unix() / 86400
}

// Rand returns a deterministic random source for a slot and theme salt.
func Rand(slot int64, salt uint64) *rand.Rand {
	return rand.New(rand.NewPCG(uint64(slot), salt))
}
