package bio

import (
	"errors"
	"fmt"
	"math/rand/v2"
	"sort"
	"strings"
	"time"
)

// DefaultTheme is the name of the always-available fallback theme.
const DefaultTheme = "default"

// Generator produces a bio for a UTC day slot.
type Generator interface {
	Generate(slot int64) string
}

// Selection is the theme chosen for this build.
type Selection struct {
	Name      string
	Generator Generator
}

var themes = map[string]Generator{}

// Register adds a theme to the compile-time registry.
func Register(name string, g Generator) {
	if name == "" || g == nil {
		panic("bio: invalid theme registration")
	}
	if _, ok := themes[name]; ok {
		panic("bio: duplicate theme " + name)
	}
	themes[name] = g
}

// Select returns the single non-default theme, or the default theme.
func Select() (Selection, error) {
	if len(themes) == 0 {
		return Selection{}, errors.New("bio: no themes registered")
	}
	names := make([]string, 0, len(themes))
	for name := range themes {
		names = append(names, name)
	}
	sort.Strings(names)
	nonDefault := make([]string, 0, len(names))
	for _, name := range names {
		if name != DefaultTheme {
			nonDefault = append(nonDefault, name)
		}
	}
	switch len(nonDefault) {
	case 0:
		return Selection{Name: DefaultTheme, Generator: themes[DefaultTheme]}, nil
	case 1:
		return Selection{Name: nonDefault[0], Generator: themes[nonDefault[0]]}, nil
	default:
		return Selection{}, fmt.Errorf(
			"bio: multiple themes compiled in (%s); build with exactly one bio_* tag",
			strings.Join(nonDefault, ", "),
		)
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
