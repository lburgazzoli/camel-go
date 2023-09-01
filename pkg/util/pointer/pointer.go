package pointer

func Of[T any](value T) *T {
	return &value
}

func Is[T comparable](ptr *T, value T) bool {
	return ptr != nil && *ptr == value
}
