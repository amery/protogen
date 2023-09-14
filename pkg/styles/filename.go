package styles

import (
	"path/filepath"
	"strings"

	"github.com/amery/protogen/pkg/protogen"
)

// FilenameWithSuffixes concatenates suffixes to the name of a File
// without extension
func (style Style) FilenameWithSuffixes(file *protogen.File, suffixes ...string) string {
	var singleDot bool
	var parts []string

	// split source file name without .proto
	dirName, fileName := filepath.Split(file.Base())

	switch {
	case fileName == "":
		// skip
	case style.FilenameSingleDot:
		// single-dot
		prefix := strings.Split(fileName, ".")
		parts = append(parts, prefix...)
		singleDot = true
	default:
		// multi-dot
		parts = append(parts, fileName)
	}

	parts = append(parts, suffixes...)

	if singleDot {
		// single dot
		if l := len(parts); l > 2 {
			// 3+ parts
			base := strings.Join(parts[:l-1], "_")
			ext := parts[l-1]

			parts = []string{base, ext}
		}
	}

	return dirName + strings.Join(parts, ".")
}
