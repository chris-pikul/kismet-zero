package lang

import "math/rand/v2"

type Phonology struct {
	Consonants PhonemeList `json:"consonants"`
	Vowels     PhonemeList `json:"vowels"`
}

func (p *Phonology) RandomConsonant(rng *rand.Rand) *Phoneme {
	return p.Consonants.WeightedChoice(rng)
}

func (p *Phonology) RandomVowel(rng *rand.Rand) *Phoneme {
	return p.Vowels.WeightedChoice(rng)
}

func NewPhonology(rng *rand.Rand, pool *PhonemePool) *Phonology {
	consonantCount := rng.IntN(7) + 6 // 6–12
	vowelCount := rng.IntN(4) + 3     // 3–6

	selected := make(map[string]struct{})
	consonants := make(PhonemeList, 0, consonantCount)
	vowels := make(PhonemeList, 0, vowelCount)

	// Shuffle pool copies to ensure variety
	globalConsonants := append([]Phoneme(nil), pool.Consonants...)
	globalVowels := append([]Phoneme(nil), pool.Vowels...)
	rand.Shuffle(len(globalConsonants), func(i, j int) {
		globalConsonants[i], globalConsonants[j] = globalConsonants[j], globalConsonants[i]
	})
	rand.Shuffle(len(globalVowels), func(i, j int) {
		globalVowels[i], globalVowels[j] = globalVowels[j], globalVowels[i]
	})

	// Select consonants
	for _, ph := range globalConsonants {
		if _, seen := selected[ph.Symbol]; seen {
			continue
		}
		selected[ph.Symbol] = struct{}{}

		ph.Weight = rng.Float32()*0.9 + 0.1 // Local weight override (0.1–1.0)
		consonants = append(consonants, ph)

		if len(consonants) >= consonantCount {
			break
		}
	}

	// Select vowels
	for _, ph := range globalVowels {
		if _, seen := selected[ph.Symbol]; seen {
			continue
		}
		selected[ph.Symbol] = struct{}{}

		ph.Weight = rng.Float32()*0.9 + 0.1
		vowels = append(vowels, ph)

		if len(vowels) >= vowelCount {
			break
		}
	}

	return &Phonology{
		Consonants: consonants,
		Vowels:     vowels,
	}
}
