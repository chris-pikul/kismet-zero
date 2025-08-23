package phoneme

import "math/rand/v2"

// Phoneme represents a single speech sound unit with its articulatory features
// and metadata for language generation.
type Phoneme struct {
	Symbol string      `json:"symbol"`
	Type   PhonemeType `json:"type"`
	Weight float32     `json:"weight"`

	Consonant *ConsonantSpec `json:"consonant,omitempty"`
	Vowel     *VowelSpec     `json:"vowel,omitempty"`
}

// PhonemeList is a collection of phonemes that supports both weighted and
// unweighted random selection.
type PhonemeList []Phoneme

// UnweightedChoice selects a random phoneme from the list with equal probability
// for each phoneme, ignoring their weight values.
func (pl PhonemeList) UnweightedChoice(rng *rand.Rand) *Phoneme {
	if len(pl) == 0 {
		return nil
	}
	return &pl[rng.IntN(len(pl))]
}

// WeightedChoice selects a random phoneme from the list using their weight
// values to determine selection probability. Higher weights increase the chance
// of selection.
func (pl PhonemeList) WeightedChoice(rng *rand.Rand) *Phoneme {
	if len(pl) == 0 {
		return nil
	}

	var total float32
	for _, p := range pl {
		total += p.Weight
	}

	if total == 0 {
		panic("PhonemeList with zero total weight cannot be sampled from")
	}

	target := rng.Float32() * total
	var cumulative float32
	for i := range pl {
		cumulative += pl[i].Weight
		if cumulative >= target {
			return &pl[i]
		}
	}

	// Edge case fallback — shouldn't normally happen
	return &pl[len(pl)-1]
}
