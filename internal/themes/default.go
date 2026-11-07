package themes

import "biosky/internal/bio"

type static string

func (s static) Generate(int64) string { return string(s) }

func init() {
	bio.Register(bio.DefaultTheme, static("quietly existing, one day at a time"))
}
