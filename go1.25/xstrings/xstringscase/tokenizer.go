package xstringscase

import (
	"slices"
	"unicode"

	"github.com/Deimvis/go-ext/go1.25/xoptional"
)

// Terms:
// - Token - sequence of chars that should be
// formatted according to case
// (e.g. delimited with '-' for snake case).
// - Statement - sequence of tokens,
// each sequence is formatted independently
// by each case
// (all characters between statements should be preserved).

// tokenizer expects single Statement
// on its input.
// If tokenizer reaches Statement end
// it stops.
type tokenizer struct {
	d []byte
	i int
}

// Next returns indices (i, j) indicating
// that token is in [i:j] slice.
func (t *tokenizer) Next() (int, int, bool) {
	tokenBegin := xoptional.New[int]()
	for ; t.i < len(t.d); t.i++ {
		if isStatementDelim(t.d[t.i]) {
			if tokenBegin.HasValue() {
				return tokenBegin.Value(), t.i, true
			}
			return 0, 0, false
		}
		if t.isTokenBegin(t.i-1, t.i) {
			if tokenBegin.HasValue() {
				return tokenBegin.Value(), t.i, true
			}
			tokenBegin.SetValue(t.i)
		} else if tokenBegin.HasValue() && t.isTokenEnd(t.i-1, t.i) {
			return tokenBegin.Value(), t.i, true
		}
	}
	if tokenBegin.HasValue() {
		return tokenBegin.Value(), len(t.d), true
	}
	return 0, 0, false
}

// isTokenBegin reports whether new token starts with char j
func (t *tokenizer) isTokenBegin(i, j int) bool {
	// algorithm description:
	// - text is splitted into tokens by
	//  - delimiters (continuous delimiters are counted as single delimiter)
	//  - capital characters after non-capitals
	// algorithm notes:
	// - not letter chars assumed to be a part of tokens
	//   and do not cause splitting into tokens
	// - by default
	if isTokenDelim(t.d[j]) {
		return false
	}
	if i == -1 {
		return true
	}
	if isTokenDelim(t.d[i]) {
		return true
	}
	return isLower(t.d[i]) && isUpper(t.d[j])
}

// isTokenEnd reports whether current token ends with char j.
// it expects that token is started and does not validate that.
func (t *tokenizer) isTokenEnd(i, j int) bool {
	if isTokenDelim(t.d[j]) {
		return true
	}
	return isLower(t.d[i]) && isUpper(t.d[j])
}

func isTokenDelim(b byte) bool {
	return slices.Contains(tokenDelimChars, b)
}

func isStatementDelim(b byte) bool {
	return !slices.Contains(tokenDelimChars, b) && unicode.IsSpace(rune(b))
}

func isUpper(b byte) bool {
	return 'A' <= b && b <= 'Z'
}

func isLower(b byte) bool {
	return 'a' <= b && b <= 'z'
}

var (
	tokenDelimChars = []byte{' ', '_', '-', '.'}
)
