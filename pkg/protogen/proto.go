package protogen

// Pointer returns a pointer to the given value
func Pointer[T any](v T) *T {
	return &v
}

// PointerOrNil returns a pointer to the given value unless
// it's zero, in which case it will return nil
func PointerOrNil[T any](v T) *T {
	if IsZero(v) {
		return nil
	}
	return &v
}
