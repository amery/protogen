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

// CutLastFunc slices s around the last match of a rune checker,
// returning the text before and after it.
// The found result reports whether there was a match.
// If there is no match in s, cut returns "", s, false.
func CutLastFunc(s string, f func(rune) bool) (before string, after string, found bool) {
	i := strings.LastIndexFunc(s, f)
	if i < 0 {
		return "", s, false
	}

	return s[:i], s[i+1:], true
}

// CutLastRune slices s around the last instance of the given rune,
// returning the text before and after it.
// The found result reports whether the rune appears in s.
// If the rune does not appear in s, cut returns "", s, false.
func CutLastRune(s string, r rune) (before string, after string, found bool) {
	return CutLastFunc(s, func(r2 rune) bool { return r == r2 })
}
