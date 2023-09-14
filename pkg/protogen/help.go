package protogen

import (
	"sort"
	"strings"
)

func optional2[T any](p *T, fallback T) (T, bool) {
	if p != nil {
		return *p, true
	}
	return fallback, false
}

func optional[T any](p *T, fallback T) T {
	v, _ := optional2(p, fallback)
	return v
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

// SplitName the last element of a dot delimited name
func SplitName(fullname string) (before string, after string, found bool) {
	return CutLastFunc(fullname, func(r rune) bool { return r == '.' })
}

// Sort sorts a slice of pointers
func Sort[T any](s []*T, less func(a, b *T) bool) {
	sort.Slice(s, func(i, j int) bool {
		return less(s[i], s[j])
	})
}
