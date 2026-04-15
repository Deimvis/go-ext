package xstrings

import (
	"golang.org/x/text/runes"
)

func BelongsTo(s string, charset runes.Set) bool {
	for _, r := range s {
		if !charset.Contains(r) {
			return false
		}
	}
	return true
}
