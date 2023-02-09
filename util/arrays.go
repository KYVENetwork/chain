package util

// RemoveFromUint64ArrayStable removes the first occurrence of `match` in the given `array`.
// It returns True if one element was removed. False otherwise.
// The order of the remaining elements is not changed.
func RemoveFromUint64ArrayStable(array []uint64, match uint64) ([]uint64, bool) {
	for i, other := range array {
		if other == match {
			return append(array[0:i], array[i+1:]...), true
		}
	}
	return array, false
}

func RemoveFromStringArrayStable(array []string, match string) ([]string, bool) {
	for i, other := range array {
		if other == match {
			return append(array[0:i], array[i+1:]...), true
		}
	}
	return array, false
}

func ContainsUint64(array []uint64, match uint64) bool {
	for _, other := range array {
		if other == match {
			return true
		}
	}
	return false
}

func ContainsString(array []string, match string) bool {
	for _, other := range array {
		if other == match {
			return true
		}
	}
	return false
}
