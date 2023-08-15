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

// SplitName the last element of a dot delimited name
func SplitName(fullname string) (prefix string, name string, found bool) {
	i := strings.LastIndexByte(fullname, byte('.'))
	if i < 0 {
		// no prefix in name
		return "", fullname, false
	}

	return fullname[:i], fullname[i+1:], true
}

// Sort sorts a slice of pointers
func Sort[T any](s []*T, less func(a, b *T) bool) {
	sort.Slice(s, func(i, j int) bool {
		return less(s[i], s[j])
	})
}
