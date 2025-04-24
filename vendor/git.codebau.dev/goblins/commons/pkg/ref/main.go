package ref

// AsRef returns a reference to the element given as t.
func AsRef[T any](t T) *T {
	return &t
}
