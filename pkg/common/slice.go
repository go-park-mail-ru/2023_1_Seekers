package common

func RemoveFromSlice[T comparable](slice []T, el T) ([]T, bool) {
	for i, v := range slice {
		if v == el {
			return append(slice[:i], slice[i+1:]...), true
		}
	}
	return slice, false
}
