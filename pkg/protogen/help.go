package protogen

import (
	"reflect"
	"sort"

	"github.com/amery/protogen/pkg/text"
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

// IsZero checks if a given value is zero, either using
// the IsZero() bool interface or reflection
func IsZero(vi any) bool {
	if p, ok := vi.(interface {
		IsZero() bool
	}); ok {
		return p.IsZero()
	}

	v := reflect.ValueOf(vi)

	switch v.Kind() {
	case reflect.Pointer:
		// (*T)(nil)
		return v.IsNil()
	case reflect.Invalid:
		// nil
		return true
	default:
		// T{}
		return v.IsZero()
	}
}

// IsNil checks if a given value is nil, regardless the type
func IsNil(vi any) bool {
	v := reflect.ValueOf(vi)
	switch v.Kind() {
	case reflect.Pointer:
		// *T(nl)
		return v.IsNil()
	case reflect.Invalid:
		// nil
		return true
	default:
		return false
	}
}

// SplitName the last element of a dot delimited name
func SplitName(name string) (before string, after string, found bool) {
	return text.CutLastFunc(name, func(r rune) bool { return r == '.' })
}

// Sort sorts a slice of pointers
func Sort[T any](s []*T, less func(a, b *T) bool) {
	sort.Slice(s, func(i, j int) bool {
		return less(s[i], s[j])
	})
}
