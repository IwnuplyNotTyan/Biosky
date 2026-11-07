//go:build bio_kaomoji

package bio

import (
	_ "embed"
	"encoding/json"
)

const theme = ThemeKaomoji

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

// Generate returns one kaomoji for the UTC day slot.
func Generate(slot int64) string {
	return kaomojiList[Rand(slot, 0x6b616f6d6f6a69).IntN(len(kaomojiList))]
}
