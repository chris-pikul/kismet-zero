package rnd

import "math/rand/v2"

// Die returns an unsigned integer in inclusive range of [1..N] where N is the
// `sides` argument. Simulation of a single die roll.
func Die(rng *rand.Rand, sides uint) uint {
	return rng.UintN(sides) + 1
}

// Dice performs a roll of N dice with given sides and returns the sum total.
func Dice(rng *rand.Rand, n, sides uint) uint {
	var result uint
	for i := uint(0); i < n; i++ {
		result += Die(rng, sides)
	}
	return result
}
