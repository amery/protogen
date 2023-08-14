package protogen

import "strings"

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
