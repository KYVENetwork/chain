package util

func MinInt(a int, b int) int {
	if a < b {
		return a
	}
	return b
}

func MinUInt64(a uint64, b uint64) uint64 {
	if a < b {
		return a
	}
	return b
}

func MaxUInt64(a uint64, b uint64) uint64 {
	if a > b {
		return a
	}
	return b
}
