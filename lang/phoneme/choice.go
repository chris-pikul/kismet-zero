package phoneme

import "math/rand/v2"

// ChoiceOptions provides configuration for phoneme selection strategies.
type ChoiceOptions struct {
	// UseWeighted determines whether to use weighted random selection
	// or uniform random selection.
	UseWeighted bool

	// FilterByType restricts selection to a specific phoneme type.
	// If nil, all phonemes are considered.
	FilterByType *PhonemeType

	// FilterByManner restricts consonant selection to a specific manner.
	// Only applies when FilterByType is PhonemeTypeConsonant.
	FilterByManner *ConsonantManner

	// FilterByHeight restricts vowel selection to a specific height.
	// Only applies when FilterByType is PhonemeTypeVowel.
	FilterByHeight *VowelHeight
}

// SelectFromPool selects phonemes from a pool based on the given options.
// Returns a slice of selected phonemes and any error encountered.
func SelectFromPool(pool *PhonemePool, rng *rand.Rand, count int, options ChoiceOptions) (PhonemeList, error) {
	if count <= 0 {
		return PhonemeList{}, nil
	}

	// Handle nil pool
	if pool == nil {
		return PhonemeList{}, nil
	}

	// Apply filters to get candidate phonemes
	candidates := filterPool(pool, options)
	if len(candidates) == 0 {
		return PhonemeList{}, nil
	}

	// Select phonemes
	selected := make(PhonemeList, 0, count)
	for i := 0; i < count && len(candidates) > 0; i++ {
		var chosen *Phoneme
		if options.UseWeighted {
			chosen = candidates.WeightedChoice(rng)
		} else {
			chosen = candidates.UnweightedChoice(rng)
		}

		if chosen != nil {
			selected = append(selected, *chosen)
			// Remove chosen phoneme to avoid duplicates
			candidates = removePhoneme(candidates, chosen.Symbol)
		}
	}

	return selected, nil
}

// filterPool applies the given filters to a phoneme pool and returns
// a filtered list of candidates.
func filterPool(pool *PhonemePool, options ChoiceOptions) PhonemeList {
	var candidates PhonemeList

	// Add consonants if not filtered out
	if options.FilterByType == nil || *options.FilterByType == PhonemeTypeConsonant {
		for _, consonant := range pool.Consonants {
			if options.FilterByManner == nil || consonant.Consonant.Manner == *options.FilterByManner {
				candidates = append(candidates, consonant)
			}
		}
	}

	// Add vowels if not filtered out
	if options.FilterByType == nil || *options.FilterByType == PhonemeTypeVowel {
		for _, vowel := range pool.Vowels {
			if options.FilterByHeight == nil || vowel.Vowel.Height == *options.FilterByHeight {
				candidates = append(candidates, vowel)
			}
		}
	}

	return candidates
}

// removePhoneme removes a phoneme with the given symbol from the list.
func removePhoneme(list PhonemeList, symbol string) PhonemeList {
	for i, phoneme := range list {
		if phoneme.Symbol == symbol {
			return append(list[:i], list[i+1:]...)
		}
	}
	return list
}

// SelectRandomConsonants selects a specified number of random consonants
// from the pool, optionally filtered by manner.
func SelectRandomConsonants(pool *PhonemePool, rng *rand.Rand, count int, manner *ConsonantManner) (PhonemeList, error) {
	consonantType := PhonemeTypeConsonant
	options := ChoiceOptions{
		UseWeighted:    false,
		FilterByType:   &consonantType,
		FilterByManner: manner,
	}
	return SelectFromPool(pool, rng, count, options)
}

// SelectRandomVowels selects a specified number of random vowels
// from the pool, optionally filtered by height.
func SelectRandomVowels(pool *PhonemePool, rng *rand.Rand, count int, height *VowelHeight) (PhonemeList, error) {
	vowelType := PhonemeTypeVowel
	options := ChoiceOptions{
		UseWeighted:    false,
		FilterByType:   &vowelType,
		FilterByHeight: height,
	}
	return SelectFromPool(pool, rng, count, options)
}

// SelectWeightedConsonants selects a specified number of consonants
// using weighted random selection, optionally filtered by manner.
func SelectWeightedConsonants(pool *PhonemePool, rng *rand.Rand, count int, manner *ConsonantManner) (PhonemeList, error) {
	consonantType := PhonemeTypeConsonant
	options := ChoiceOptions{
		UseWeighted:    true,
		FilterByType:   &consonantType,
		FilterByManner: manner,
	}
	return SelectFromPool(pool, rng, count, options)
}

// SelectWeightedVowels selects a specified number of vowels
// using weighted random selection, optionally filtered by height.
func SelectWeightedVowels(pool *PhonemePool, rng *rand.Rand, count int, height *VowelHeight) (PhonemeList, error) {
	vowelType := PhonemeTypeVowel
	options := ChoiceOptions{
		UseWeighted:    true,
		FilterByType:   &vowelType,
		FilterByHeight: height,
	}
	return SelectFromPool(pool, rng, count, options)
}
