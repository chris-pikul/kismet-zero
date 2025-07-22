package lang

var PhonemePoolCommon = PhonemePool{
	Consonants: []Phoneme{
		// Plosives
		{Symbol: "p", Type: PhonemeTypeConsonant, Consonant: &ConsonantSpec{Manner: ConsonantMannerPlosive}, Weight: 1.0},
		{Symbol: "b", Type: PhonemeTypeConsonant, Consonant: &ConsonantSpec{Manner: ConsonantMannerPlosive}, Weight: 1.0},
		{Symbol: "t", Type: PhonemeTypeConsonant, Consonant: &ConsonantSpec{Manner: ConsonantMannerPlosive}, Weight: 1.0},
		{Symbol: "d", Type: PhonemeTypeConsonant, Consonant: &ConsonantSpec{Manner: ConsonantMannerPlosive}, Weight: 1.0},
		{Symbol: "k", Type: PhonemeTypeConsonant, Consonant: &ConsonantSpec{Manner: ConsonantMannerPlosive}, Weight: 1.0},
		{Symbol: "g", Type: PhonemeTypeConsonant, Consonant: &ConsonantSpec{Manner: ConsonantMannerPlosive}, Weight: 1.0},

		// Nasals
		{Symbol: "m", Type: PhonemeTypeConsonant, Consonant: &ConsonantSpec{Manner: ConsonantMannerNasal}, Weight: 1.0},
		{Symbol: "n", Type: PhonemeTypeConsonant, Consonant: &ConsonantSpec{Manner: ConsonantMannerNasal}, Weight: 1.0},
		{Symbol: "ŋ", Type: PhonemeTypeConsonant, Consonant: &ConsonantSpec{Manner: ConsonantMannerNasal}, Weight: 0.8},

		// Fricatives
		{Symbol: "f", Type: PhonemeTypeConsonant, Consonant: &ConsonantSpec{Manner: ConsonantMannerFricative}, Weight: 1.0},
		{Symbol: "v", Type: PhonemeTypeConsonant, Consonant: &ConsonantSpec{Manner: ConsonantMannerFricative}, Weight: 1.0},
		{Symbol: "s", Type: PhonemeTypeConsonant, Consonant: &ConsonantSpec{Manner: ConsonantMannerSibilants}, Weight: 1.0},
		{Symbol: "z", Type: PhonemeTypeConsonant, Consonant: &ConsonantSpec{Manner: ConsonantMannerSibilants}, Weight: 1.0},

		// Approximants
		{Symbol: "l", Type: PhonemeTypeConsonant, Consonant: &ConsonantSpec{Manner: ConsonantMannerLateralApproximant}, Weight: 1.0},
		{Symbol: "r", Type: PhonemeTypeConsonant, Consonant: &ConsonantSpec{Manner: ConsonantMannerTrill}, Weight: 0.8},
		{Symbol: "j", Type: PhonemeTypeConsonant, Consonant: &ConsonantSpec{Manner: ConsonantMannerSemivowel}, Weight: 1.0},
		{Symbol: "w", Type: PhonemeTypeConsonant, Consonant: &ConsonantSpec{Manner: ConsonantMannerSemivowel}, Weight: 1.0},
	},
	Vowels: []Phoneme{
		{Symbol: "a", Type: PhonemeTypeVowel, Vowel: &VowelSpec{Height: VowelHeightLow, Backness: VowelBacknessFront, Rounded: false}, Weight: 1.0},
		{Symbol: "e", Type: PhonemeTypeVowel, Vowel: &VowelSpec{Height: VowelHeightMid, Backness: VowelBacknessFront, Rounded: false}, Weight: 1.0},
		{Symbol: "i", Type: PhonemeTypeVowel, Vowel: &VowelSpec{Height: VowelHeightHigh, Backness: VowelBacknessFront, Rounded: false}, Weight: 1.0},
		{Symbol: "o", Type: PhonemeTypeVowel, Vowel: &VowelSpec{Height: VowelHeightMid, Backness: VowelBacknessBack, Rounded: true}, Weight: 1.0},
		{Symbol: "u", Type: PhonemeTypeVowel, Vowel: &VowelSpec{Height: VowelHeightHigh, Backness: VowelBacknessBack, Rounded: true}, Weight: 1.0},
		{Symbol: "ə", Type: PhonemeTypeVowel, Vowel: &VowelSpec{Height: VowelHeightMid, Backness: VowelBacknessCentral, Rounded: false}, Weight: 0.7},
	},
}
