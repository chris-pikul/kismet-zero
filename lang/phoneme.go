package lang

import "math/rand/v2"

type Phoneme struct {
	Symbol string      `json:"symbol"`
	Type   PhonemeType `json:"type"`
	Weight float32     `json:"weight"`

	Consonant *ConsonantSpec `json:"consonant,omitempty"`
	Vowel     *VowelSpec     `json:"vowel,omitempty"`
}

type PhonemeList []Phoneme

func (pl PhonemeList) UnweightedChoice(rng *rand.Rand) *Phoneme {
	if len(pl) == 0 {
		return nil
	}
	return &pl[rng.IntN(len(pl))]
}

func (pl PhonemeList) WeightedChoice(rng *rand.Rand) *Phoneme {
	var total float32
	for _, p := range pl {
		total += p.Weight
	}

	if total == 0 {
		panic("empty PhonemeList cannot be sampled from")
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
