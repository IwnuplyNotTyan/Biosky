//go:build !bio_kaomoji && !bio_moon

package bio

const theme = ThemeDefault

const defaultText = "quietly existing, one day at a time"

// Generate returns the fixed default bio.
func Generate(int64) string { return defaultText }
