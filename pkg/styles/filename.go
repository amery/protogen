package styles

import (
	"strings"

	"github.com/amery/protogen/pkg/protogen"
)

// FilenameWithSuffixes concatenates suffixes to the name of a File
// without extension
func (style Style) FilenameWithSuffixes(file *protogen.File, suffixes ...string) string {
	var singleDot bool
	var parts []string

	base := file.Base()

	switch {
	case base == "":
		// skip
	case style.FilenameSingleDot:
		// single-dot
		prefix := strings.Split(base, ".")
		parts = append(parts, prefix...)
		singleDot = true
	default:
		// multi-dot
		parts = append(parts, base)
	}

	parts = append(parts, suffixes...)

	if singleDot {
		// single dot
		if l := len(parts); l > 2 {
			// 3+ parts
			base := strings.Join(parts[:l-1], "_")
			ext := parts[l-1]
			return base + "." + ext
		}
	}

	return strings.Join(parts, ".")
}
