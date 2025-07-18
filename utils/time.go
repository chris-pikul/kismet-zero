package utils

import (
	"math"
	"time"
)

// NanoSecondDuration converts a [uint] number of nanoseconds to a [time.Duration]
// while protecting from overflows.
func NanoSecondDuration(ns uint64) time.Duration {
	if ns > math.MaxInt64 {
		return 0
	}
	return time.Duration(ns)
}
