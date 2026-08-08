package lib

import "math"

func MulOverflow(a, b int64) bool {
	if a == 0 || b == 0 {
		return false
	}

	if a == -1 {
		return b == math.MinInt64
	}
	if b == -1 {
		return a == math.MinInt64
	}

	if a > 0 {
		if b > 0 {
			return a > math.MaxInt64/b
		}
		return b < math.MinInt64/a
	}

	if b > 0 {
		return a < math.MinInt64/b
	}
	return a < math.MaxInt64/b
}
