package text

import "strings"

// CutPrefix removes, if present, a prefix from a string.
func CutPrefix(s, prefix string) (string, bool) {
	if !strings.HasPrefix(s, prefix) {
		return s, false
	}

	return s[len(prefix):], true
}

// CutSuffix removes, if present, a suffix from a string.
func CutSuffix(s, suffix string) (string, bool) {
	if !strings.HasSuffix(s, suffix) {
		return s, false
	}

	return s[:len(s)-len(suffix)], true
}
