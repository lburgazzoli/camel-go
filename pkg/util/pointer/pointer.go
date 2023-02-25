package pointer

// Of ---
// TODO: improve.
func Of[T any](value T) *T {
	return &value
}

// Is ---
// TODO: improve.
func Is[T comparable](ptr *T, value T) bool {
	return ptr != nil && *ptr == value
}
