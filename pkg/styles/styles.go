// Package styles assists on composing names
package styles

// Style describes how we want to render names
type Style struct {
	// FilenameSingleDot tells only the last part of a filename should
	// be delimited by a dot
	FilenameSingleDot bool
}
