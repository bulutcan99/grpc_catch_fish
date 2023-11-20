package utility

func Contains[T comparable](elems []T, v T, fn func(value T, element T) bool) (bool, int) {
	if fn == nil {
		fn = func(value T, element T) bool {
			return value == element
		}
	}
	for index, s := range elems {
		if fn(v, s) {
			return true, index
		}
	}
	return false, -1
}
