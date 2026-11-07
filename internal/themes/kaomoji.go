//go:build bio_kaomoji

package themes

import "biosky/internal/bio"

type kaomoji struct{}

var kaomojiLines = []string{
	"(｡•̀ᴗ-)✧ quietly existing on the timeline",
	"(￣▽￣)ノ see you at the next ping",
	"(´-ω-`) holding it together, mostly",
	"( •_•)>⌐■-■ no thoughts, just vibes",
	"(⌐■_■) offline in spirit",
}

func (kaomoji) Generate(slot int64) string {
	r := bio.Rand(slot, 0x6b616f6d6f6a69)
	return kaomojiLines[r.IntN(len(kaomojiLines))]
}

func init() {
	bio.Register("kaomoji", kaomoji{})
}
