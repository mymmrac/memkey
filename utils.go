package memkey

// zero returns zero value for the specified type
func zero[T any]() T {
	var value T
	return value
}
