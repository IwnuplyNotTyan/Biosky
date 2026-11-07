//go:build kaomoji

package bio

import (
	_ "embed"
	"encoding/json"
	"math/rand/v2"
)

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

// Generate returns a random kaomoji.
func Generate(int64) string {
	return kaomojiList[rand.IntN(len(kaomojiList))]
}
